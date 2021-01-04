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
	Type   string `json:"type"`
	Origin string `json:"origin"`
	Data   []byte `json:"data"`
	Tags   []*Tag `json:"tags"`
}

func NewInternal() (*Internal, error) {
	var self *Internal

	common, err := NewCommon()
	if err != nil {
		return self, err
	}

	self = &Internal{
		Common: common,
		Tags:   []*Tag{},
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

func (self *Internal) Set(key string, value interface{}) error {
	switch key {
	case "type":
		val, ok := value.(string)
		if !ok {
			return fmt.Errorf("Invalid value type: %#v", value)
		}
		if len(val) > 64 {
			return fmt.Errorf("Type is over 64 characters: %s", val)
		}
		self.Type = val
	case "origin":
		val, ok := value.(string)
		if !ok {
			return fmt.Errorf("Invalid value type: %#v", value)
		}
		if len(val) > 64 {
			return fmt.Errorf("Origin is over 64 characters: %s", val)
		}
		self.Origin = val
	case "data":
		val, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("Invalid value type: %#v", value)
		}
		self.Data = val
	default:
		return fmt.Errorf("Invalid key: %s", key)
	}

	return nil
}

func (self *Internal) Save(store *sql.DB) error {
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

	statement, err := store.Prepare(`
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

func (self *Internal) Search(crateType) SearchQuery
