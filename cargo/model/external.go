package model

import (
	"io"
	"fmt"
	"log"
	"encoding/json"
	"database/sql"
)

type External struct {
	Common
	Type    string          `json:"type"`
	Name    string          `json:"name"`
	Body    string          `json:"body"`
	Data    []byte          `json:"data"`
	Tags    []Entity        `json:"tags"`
	Mapping map[int64]int64 `json:"-"`
	Meta    *Internal       `json:"-"`
}

func NewExternal(store *sql.DB) (*External, error) {
	var self *External
	var data []byte

	common, err := NewCommon(store, "external")
	if err != nil {
		return self, err
	}

	self = &External{
		Common:  common,
		Data:    data,
		Tags:    []Entity{},
		Mapping: make(map[int64]int64),
	}

	return self, nil
}

func (self *External) Display() {
	log.Printf(
		"External<gid:%d uid:%x add:%s upd:%s flg:%d>\n",
		self.ID,
		self.UUID[:8],
		self.Added.String(),
		self.Updated.String(),
		self.Flag,
	)
}

func (self *External) Encode(w io.Writer) error {
	err := json.NewEncoder(w).Encode(&self)
	return err
}

func (self *External) Set(key string, value []byte) error {
	switch key {
	case "flag":
		self.Flag = value[0]
	case "type":
		val := string(value)
		if len(val) > 64 {
			return fmt.Errorf("Type is over 64 characters: %s", val)
		}
		self.Type = val
	case "name":
		val := string(value)
		if len(val) > 64 {
			return fmt.Errorf("Name is over 64 characters: %s", val)
		}
		self.Name = val
	case "body":
		self.Body = string(value)
	default:
		return fmt.Errorf("Invalid key: %s", key)
	}

	return nil
}

func (self *External) Save() error {
	if len(self.UUID) != 32 {
		return fmt.Errorf("UUID is not 32 bytes: %x", self.UUID)
	}

	if len(self.Type) > 64 {
		return fmt.Errorf("Type is over 64 characters: %s", self.Type)
	}

	if len(self.Name) > 64 {
		return fmt.Errorf("Name is over 64 characters: %s", self.Name)
	}

	if len(self.Body) == 0 {
		return fmt.Errorf("Body is missing: %s", self.Body)
	}

	if len(self.Data) > 0 && len(self.Data) != 32 {
		return fmt.Errorf("Data is invalid: %x", self.Data)
	}

	statement, err := self.Store.Prepare(`
		INSERT INTO external (uuid, flag, type, name, body)
		VALUES (?, ?, ?, ?, ?);
	`)

	if err != nil {
		return err
	}

	defer statement.Close()
	result, err := statement.Exec(
		self.UUID,
		self.Flag,
		self.Type,
		self.Name,
		self.Body,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	self.ID = id
	return nil
}

func (self *External) Update() error {
	self.Updated = Now()
	statement, err := self.Store.Prepare(`
		UPDATE external
		SET updated = ?, flag = ?, type = ?, name = ?, body = ?
		WHERE id = ?
	`)

	if err != nil {
		return err
	}

	defer statement.Close()
	_, err = statement.Exec(
		self.Updated,
		self.Flag,
		self.Type,
		self.Name,
		self.Body,
		self.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (self *External) Delete() error {
	statement, err := self.Store.Prepare(`
		DELETE FROM external WHERE id = ?;
	`)

	if err != nil {
		return err
	}

	defer statement.Close()
	_, err = statement.Exec(
		self.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (self *External) ExportMetadata() (int64, string) {
	return self.ID, self.Mapper
}

func (self *External) Map(entity Entity) error {
	id, mapper := entity.ExportMetadata()
	if self.Mapper == mapper {
		return fmt.Errorf("Cannot create mapping with: %s", mapper)
	}

	statement, err := self.Store.Prepare(fmt.Sprintf(`
		INSERT INTO mapping (external_id, %s_id) VALUES (?, ?);
	`, mapper))

	if err != nil {
		return err
	}

	defer statement.Close()
	result, err := statement.Exec(
		self.ID,
		id,
	)

	if err != nil {
		return err
	}

	mappingID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	self.Mapping[id] = mappingID
	if mapper == "tag" {
		self.Tags = append(self.Tags, entity)
	}

	return nil
}

func (self *External) Unmap(entity Entity) error {
	id, mapper := entity.ExportMetadata()
	if self.Mapper == mapper {
		return fmt.Errorf("Cannot delete mapping with: %s", mapper)
	}

	mappingID, ok := self.Mapping[id]
	if !ok {
		statement, err := self.Store.Prepare(fmt.Sprintf(`
			DELETE FROM mapping
			WHERE external_id = ? AND %s_id = ?;
		`, mapper))

		if err != nil {
			return err
		}

		defer statement.Close()
		_, err = statement.Exec(
			self.ID,
			id,
		)

		if err != nil {
			return err
		}
	} else {
		delete(self.Mapping, id)

		statement, err := self.Store.Prepare(`
			DELETE FROM mapping WHERE id = ?;
		`)

		if err != nil {
			return err
		}

		defer statement.Close()
		_, err = statement.Exec(
			mappingID,
		)

		if err != nil {
			return err
		}
	}

	if mapper == "tag" {
		tags := []Entity{}
		for _, tag := range self.Tags {
			if tag == entity {
				continue
			}

			tags = append(tags, tag)
		}
		self.Tags = tags
	}

	return nil
}

func (self *External) Link(entity Entity) error {
	meta, ok := entity.(*Internal)
	if !ok {
		return fmt.Errorf("Cannot cast to Internal: %#v", entity)
	}

	self.Meta = meta
	self.Data = self.Meta.UUID

	self.Updated = Now()
	statement, err := self.Store.Prepare(`
		UPDATE external SET updated = ?, data = ? WHERE id = ?;
	`)

	if err != nil {
		return err
	}

	defer statement.Close()
	_, err = statement.Exec(
		self.Updated,
		self.Meta.ID,
		self.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (self *External) Unlink() error {
	var meta *Internal = nil

	self.Meta = meta
	self.Data = []byte{}

	self.Updated = Now()
	statement, err := self.Store.Prepare(`
		UPDATE external SET updated = ?, data = NULL WHERE id = ?;
	`)

	if err != nil {
		return err
	}

	defer statement.Close()
	_, err = statement.Exec(
		self.Updated,
		self.ID,
	)

	if err != nil {
		return err
	}

	return nil
}
