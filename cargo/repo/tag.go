package repo

import (
	"fmt"
	"time"
	"database/sql"

	"github.com/aewens/nautical/cargo/model"
)

type Tag struct {
	Store  *sql.DB
	Crates map[int64]*model.Tag
}

func NewTag(store *sql.DB) *Tag {
	return &Tag{
		Store:  store,
		Crates: make(map[int64]*model.Tag),
	}
}

func (self *Tag) Create() (model.Entity, error) {
	return model.NewTag(self.Store)
}

func (self *Tag) Load(stream Stream) {
	for entity := range stream {
		tag, ok := entity.(*model.Tag)
		if !ok {
			continue
		}

		self.Crates[tag.ID] = tag
	}
}

func (self *Tag) Import(
	id      int64,
	uuid    []byte,
	added   time.Time,
	updated time.Time,
	flag    uint8,
	label   string,
) (model.Entity, error) {
	entity, err := self.Create()
	if err != nil {
		return entity, err
	}

	tag, ok := entity.(*model.Tag)
	if !ok {
		return entity, fmt.Errorf("Cannot cast to Tag: %#v", entity)
	}

	tag.ID = id
	tag.UUID = uuid
	tag.Added = added
	tag.Updated = updated
	tag.Flag = flag
	tag.Label = label

	return entity, nil
}

func (self *Tag) Get(id int64) (model.Entity, error) {
	statement, err := self.Store.Prepare(`
		SELECT uuid, added, updated, flag, label
		FROM tag WHERE id = ?;
	`)

	if err != nil {
		return nil, err
	}

	var (
		uuid    []byte
		added   time.Time
		updated time.Time
		flag    uint8
		label   string
	)

	defer statement.Close()
	err = statement.QueryRow(id).Scan(
		&uuid,
		&added,
		&updated,
		&flag,
		&label,
	)

	if err != nil {
		return nil, err
	}

	return self.Import(
		id,
		uuid,
		added,
		updated,
		flag,
		label,
	)
}

func (self *Tag) Process(stream Stream, rows *sql.Rows) {
	defer rows.Close()
	for rows.Next() {
		var (
			id      int64
			uuid    []byte
			added   time.Time
			updated time.Time
			flag    uint8
			label   string
		)

		err := rows.Scan(
			&id,
			&uuid,
			&added,
			&updated,
			&flag,
			&label,
		)

		if err != nil {
			continue
		}

		entity, err := self.Import(
			id,
			uuid,
			added,
			updated,
			flag,
			label,
		)

		if err != nil {
			continue
		}

		stream <- entity
	}

	close(stream)
}

func (self *Tag) All() Stream {
	stream := make(Stream)

	go func() {
		rows, err := self.Store.Query(`
			SELECT id, uuid, added, updated, flag, label
			FROM tag;
		`)

		if err != nil {
			return
		}

		self.Process(stream, rows)
	}()

	return stream
}

func (self *Tag) Lookup(ids ...int64) Stream {
	stream := make(Stream)

	go func() {
		for _, id := range ids {
			entity, err := self.Get(id)
			if err == nil {
				stream <- entity
			}
		}

		close(stream)
	}()

	return stream
}

func (self *Tag) Contains(field string, search string) Stream {
	stream := make(Stream)

	go func() {
		statement, err := self.Store.Prepare(fmt.Sprintf(`
			SELECT id, uuid, added, updated, flag, label
			FROM tag WHERE %s LIKE ?;
		`, field))

		if err != nil {
			return
		}

		defer statement.Close()
		rows, err := statement.Query("%" + search + "%")

		if err != nil {
			return
		}

		self.Process(stream, rows)
	}()

	return stream
}

func (self *Tag) Equals(field string, search string) Stream {
	stream := make(Stream)

	go func() {
		statement, err := self.Store.Prepare(fmt.Sprintf(`
			SELECT id, uuid, added, updated, flag, label
			FROM tag WHERE %s = ?;
		`, field))

		if err != nil {
			return
		}

		defer statement.Close()
		rows, err := statement.Query(search)

		if err != nil {
			return
		}

		self.Process(stream, rows)
	}()

	return stream
}

func (self *Tag) Before(field string, search time.Time) Stream {
	stream := make(Stream)

	go func() {
		statement, err := self.Store.Prepare(fmt.Sprintf(`
			SELECT id, uuid, added, updated, flag, label
			FROM tag WHERE %s < ?;
		`, field))

		if err != nil {
			return
		}

		defer statement.Close()
		rows, err := statement.Query(search)

		if err != nil {
			return
		}

		self.Process(stream, rows)
	}()

	return stream
}

func (self *Tag) After(field string, search time.Time) Stream {
	stream := make(Stream)

	go func() {
		statement, err := self.Store.Prepare(fmt.Sprintf(`
			SELECT id, uuid, added, updated, flag, label
			FROM tag WHERE %s > ?;
		`, field))

		if err != nil {
			return
		}

		defer statement.Close()
		rows, err := statement.Query(search)

		if err != nil {
			return
		}

		self.Process(stream, rows)
	}()

	return stream
}

func (self *Tag) Between(
	field string,
	before time.Time,
	after time.Time,
) Stream {
	stream := make(Stream)

	go func() {
		statement, err := self.Store.Prepare(fmt.Sprintf(`
			SELECT id, uuid, added, updated, flag, label
			FROM tag WHERE %s > ? AND %s < ?;
		`, field, field))

		if err != nil {
			return
		}

		defer statement.Close()
		rows, err := statement.Query(before, after)

		if err != nil {
			return
		}

		self.Process(stream, rows)
	}()

	return stream
}
