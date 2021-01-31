package repo

import (
	"fmt"
	"time"
	"database/sql"

	"github.com/aewens/armada/cargo/model"
)

type External struct {
	Store  *sql.DB
	Crates map[int64]*model.External
}

func NewExternal(store *sql.DB) *External {
	return &External{
		Store:  store,
		Crates: make(map[int64]*model.External),
	}
}

func (self *External) Create() (model.Entity, error) {
	return model.NewExternal(self.Store)
}

func (self *External) Load(stream Stream) {
	for entity := range stream {
		external, ok := entity.(*model.External)
		if !ok {
			continue
		}

		self.Crates[external.ID] = external
	}
}

func (self *External) Import(
	id      int64,
	uuid    []byte,
	added   time.Time,
	updated time.Time,
	flag    uint8,
	etype   string,
	name    string,
	body    string,
	link    sql.NullInt64,
) (model.Entity, error) {
	entity, err := self.Create()
	if err != nil {
		return entity, err
	}

	external, ok := entity.(*model.External)
	if !ok {
		return entity, fmt.Errorf("Cannot cast to External: %#v", entity)
	}

	external.ID = id
	external.UUID = uuid
	external.Added = added
	external.Updated = updated
	external.Flag = flag
	external.Type = etype
	external.Name = name
	external.Body = body

	if link.Valid {
		internals := NewInternal(self.Store)
		ientity, err := internals.Get(link.Int64)
		if err != nil {
			return entity, err
		}

		err = external.Link(ientity)
		if err != nil {
			return entity, err
		}
	}

	return entity, nil
}

func (self *External) Get(id int64) (model.Entity, error) {
	statement, err := self.Store.Prepare(`
		SELECT uuid, added, updated, flag, type, name, body, data
		FROM external WHERE id = ?;
	`)

	if err != nil {
		return nil, err
	}

	var (
		uuid    []byte
		added   time.Time
		updated time.Time
		flag    uint8
		etype   string
		name    string
		body    string
		link    sql.NullInt64
	)

	defer statement.Close()
	err = statement.QueryRow(id).Scan(
		&uuid,
		&added,
		&updated,
		&flag,
		&etype,
		&name,
		&body,
		&link,
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
		etype,
		name,
		body,
		link,
	)
}

func (self *External) All() Stream {
	stream := make(Stream)

	go func() {
		rows, err := self.Store.Query(`
			SELECT id, uuid, added, updated, flag, type, name, body, data
			FROM external;
		`)

		if err != nil {
			return
		}

		defer rows.Close()
		for rows.Next() {
			var (
				id      int64
				uuid    []byte
				added   time.Time
				updated time.Time
				flag    uint8
				etype   string
				name    string
				body    string
				link    sql.NullInt64
			)

			err = rows.Scan(
				&id,
				&uuid,
				&added,
				&updated,
				&flag,
				&etype,
				&name,
				&body,
				&link,
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
				etype,
				name,
				body,
				link,
			)

			if err != nil {
				continue
			}

			stream <- entity
		}

		close(stream)
	}()

	return stream
}

func (self *External) Lookup(ids ...int64) Stream {
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
