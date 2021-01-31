package cargo

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"fmt"
)

func Resolve(conn string) (string, error) {
	if filepath.IsAbs(conn) {
		return conn, nil
	}

	here, err := os.Executable()
	if err != nil {
		return conn, err
	}

	path := filepath.Join(filepath.Dir(here), conn)
	conn, err = filepath.EvalSymlinks(path)
	return conn, err
}

func Wrap(conn string) string {
	params := []string{
		"_foreign_keys=1",
		"cached=shared",
	}
	return fmt.Sprintf("file:%s?%s", conn, strings.Join(params, "&"))
}

func Open(conn string, tables string) (*sql.DB, error) {
	bootstrap := false

	if conn == ":memory:" {
		bootstrap = true
	} else {
		conn, err := Resolve(conn)
		if err != nil {
			return nil, err
		}

		_, err = os.Stat(conn)
		if os.IsNotExist(err) {
			bootstrap = true
		}
	}

	uri := Wrap(conn)
	db, err := sql.Open("sqlite3", uri)
	if bootstrap && err == nil {
		_, err = db.Exec(tables)
	}

	return db, err
}

var Tables string = `
	CREATE TABLE internal (
		id INTEGER PRIMARY KEY,
		uuid BLOB(32) UNIQUE NOT NULL,
		added DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
		flag INTEGER DEFAULT 0 NOT NULL,
		type VARCHAR(64) NOT NULL,
		origin VARCHAR(64) NOT NULL,
		data BLOB NOT NULL,
		CHECK (flag >= 0 AND flag <= 255) -- force unsigned int8
	);
	CREATE TABLE external (
		id INTEGER PRIMARY KEY,
		uuid BLOB(32) UNIQUE NOT NULL,
		added DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
		flag INTEGER DEFAULT 0 NOT NULL,
		type VARCHAR(64) NOT NULL,
		name VARCHAR(64) NOT NULL,
		body TEXT NOT NULL,
		data INTEGER, -- allowed to be null
		FOREIGN KEY (data) REFERENCES internal(id) ON DELETE SET NULL,
		CHECK (flag >= 0 AND flag <= 255) -- force unsigned int8
	);
	CREATE TABLE tag (
		id INTEGER PRIMARY KEY,
		uuid BLOB(32) UNIQUE NOT NULL,
		added DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
		updated DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
		flag INTEGER DEFAULT 0 NOT NULL,
		label VARCHAR(128) UNIQUE NOT NULL,
		CHECK (flag >= 0 AND flag <= 255) -- force unsigned int8
	);
	CREATE TABLE mapping (
		id INTEGER PRIMARY KEY,
		internal_id INTEGER, -- allowed to be null
		external_id INTEGER, -- allowed to be null
		tag_id INTEGER, -- allowed to be null
		UNIQUE (internal_id, tag_id) ON CONFLICT IGNORE,
		UNIQUE (external_id, tag_id) ON CONFLICT IGNORE, -- prevent duplicates
		UNIQUE (internal_id, external_id) ON CONFLICT IGNORE,
		FOREIGN KEY (internal_id) REFERENCES internal(id) ON DELETE CASCADE,
		FOREIGN KEY (external_id) REFERENCES external(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tag(id) ON DELETE CASCADE,
		CHECK (
			(
				internal_id IS NOT NULL
			AND external_id IS NOT NULL
			AND tag_id IS NULL
			)
		OR  (
				internal_id IS NULL
			AND external_id IS NOT NULL
			AND tag_id IS NOT NULL
			)
		OR  (
				internal_id IS NOT NULL
			AND external_id IS NULL
			AND tag_id IS NOT NULL
			)
		) -- force one pair to be used only
	);
`
