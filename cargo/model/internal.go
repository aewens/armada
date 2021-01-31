package model

import (
	"io"
	"fmt"
	"log"
	"encoding/json"
	"database/sql"
)

type Internal struct {
	Common
	Type    string          `json:"type"`
	Origin  string          `json:"origin"`
	Data    []byte          `json:"data"`
	Tags    []Entity        `json:"tags"`
	Mapping map[int64]int64 `json:"-"`
}

func NewInternal(store *sql.DB) (*Internal, error) {
	var self *Internal

	common, err := NewCommon(store, "internal")
	if err != nil {
		return self, err
	}

	self = &Internal{
		Common:  common,
		Tags:    []Entity{},
		Mapping: make(map[int64]int64),
	}

	return self, nil
}

func (self *Internal) Display() {
	//log.Printf("%#v\n", self)
	log.Printf(
		"Internal<gid:%d uid:%x add:%s upd:%s flg:%d>\n",
		self.ID,
		self.UUID[:8],
		self.Added.String(),
		self.Updated.String(),
		self.Flag,
	)
}

func (self *Internal) Encode(w io.Writer) error {
	err := json.NewEncoder(w).Encode(&self)
	return err
}

func (self *Internal) Set(key string, value []byte) error {
	switch key {
	case "flag":
		self.Flag = value[0]
	case "type":
		val := string(value)
		if len(val) > 64 {
			return fmt.Errorf("Type is over 64 characters: %s", val)
		}
		self.Type = val
	case "origin":
		val := string(value)
		if len(val) > 64 {
			return fmt.Errorf("Origin is over 64 characters: %s", val)
		}
		self.Origin = val
	case "data":
		self.Data = value
	default:
		return fmt.Errorf("Invalid key: %s", key)
	}

	return nil
}

func (self *Internal) Save() error {
	if len(self.UUID) != 32 {
		return fmt.Errorf("UUID is not 32 bytes: %x", self.UUID)
	}

	if len(self.Type) > 64 {
		return fmt.Errorf("Type is over 64 characters: %s", self.Type)
	}

	if len(self.Origin) > 64 {
		return fmt.Errorf("Origin is over 64 characters: %s", self.Origin)
	}

	if len(self.Data) == 0 {
		return fmt.Errorf("Data is missing: %x", self.Data)
	}

	statement, err := self.Store.Prepare(`
		INSERT INTO internal (uuid, flag, type, origin, data)
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
		self.Origin,
		self.Data,
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

func (self *Internal) Update() error {
	self.Updated = Now()
	statement, err := self.Store.Prepare(`
		UPDATE internal
		SET updated = ?, flag = ?, type = ?, origin = ?, data = ?
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
		self.Origin,
		self.Data,
		self.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (self *Internal) Delete() error {
	statement, err := self.Store.Prepare(`
		DELETE FROM internal WHERE id = ?;
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

func (self *Internal) ExportMetadata() (int64, string) {
	return self.ID, self.Mapper
}

func (self *Internal) Map(entity Entity) error {
	id, mapper := entity.ExportMetadata()
	if self.Mapper == mapper {
		return fmt.Errorf("Cannot create mapping with: %s", mapper)
	}

	statement, err := self.Store.Prepare(fmt.Sprintf(`
		INSERT INTO mapping (internal_id, %s_id) VALUES (?, ?);
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

func (self *Internal) Unmap(entity Entity) error {
	id, mapper := entity.ExportMetadata()
	if self.Mapper == mapper {
		return fmt.Errorf("Cannot delete mapping with: %s", mapper)
	}

	mappingID, ok := self.Mapping[id]
	if !ok {
		statement, err := self.Store.Prepare(fmt.Sprintf(`
			DELETE FROM mapping
			WHERE internal_id = ? AND %s_id = ?;
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
