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
	Type   string `json:"type"`
	Name   string `json:"name"`
	Body   string `json:"body"`
	Data   []byte `json:"data"`
	Tags   []*Tag `json:"tags"`
}

func NewExternal(store *sql.DB) (*External, error) {
	var self *External
	var data []byte

	common, err := NewCommon(store)
	if err != nil {
		return self, err
	}

	self = &External{
		Common: common,
		Data:   data,
		Tags:   []*Tag{},
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

	statement, err := self.Store.Prepare(`
		UPDATE external SET type = ?, name = ?, body = ?, data = ?
		WHERE id = ?
	`)

	if err != nil {
		return err
	}

	defer statement.Close()
	_, err = statement.Exec(
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
