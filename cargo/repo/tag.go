package repo

import (
	"fmt"
	"time"
	"database/sql"

	"github.com/aewens/armada/cargo/model"
)

type Tag struct {
	Store  *sql.DB
	Crates []*model.Tag
}

func NewTag(store *sql.DB) *Tag {
	return &Tag{
		Store:  store,
		Crates: []*model.Tag{},
	}
}

func (self *Tag) Create() (model.Entity, error) {
	return model.NewTag(self.Store)
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

			err = rows.Scan(
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
	}()

	return stream
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

func (self *Tag) Load(stream Stream) {
	for entity := range stream {
		tag, ok := entity.(*model.Tag)
		if !ok {
			continue
		}

		self.Crates[tag.ID] = tag
	}
}
