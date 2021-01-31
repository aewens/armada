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
	case "data":
		if len(value) != 32 {
			return fmt.Errorf("Data is not 32 bytes: %x", value)
		}
		self.Data = value
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
		INSERT INTO external (uuid, flag, type, name, body, data)
		VALUES (?, ?, ?, ?, ?, ?);
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

func (self *External) Update(changes map[string][]byte) error {
	for key, value := range changes {
		switch key {
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
			val := string(value)
			if len(val) == 0 {
				return fmt.Errorf("Body is missing: %s", val)
			}
			self.Body = val
		case "data":
			if len(value) == 0 {
				return fmt.Errorf("Data is missing: %x", value)
			}
			self.Data = value
		default:
			return fmt.Errorf("Invalid key: %s", key)
		}
	}

	self.Updated = Now()
	statement, err := self.Store.Prepare(`
		UPDATE external
		SET updated = ?, type = ?, name = ?, body = ?, data = ?
		WHERE id = ?
	`)

	if err != nil {
		return err
	}

	defer statement.Close()
	_, err = statement.Exec(
		self.Updated,
		self.Type,
		self.Name,
		self.Body,
		self.Data,
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
