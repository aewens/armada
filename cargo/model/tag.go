package model

import (
	"io"
	"fmt"
	"log"
	"encoding/json"
	"database/sql"
)

type Tag struct {
	Common
	Label  string `json:"label"`
}

func NewTag(store *sql.DB, label string) (*Tag, error) {
	var self *Tag

	common, err := NewCommon(store, "tag")
	if err != nil {
		return self, err
	}

	if len(label) == 0 {
		return self, fmt.Errorf("Missing label: %s", label)
	}

	self = &Tag{
		Common: common,
		Label:  label,
	}

	return self, nil
}

func (self *Tag) Display() {
	log.Printf(
		"Tag<gid:%d uid:%x add:%s upd:%s flg:%d>\n",
		self.ID,
		self.UUID[:8],
		self.Added.String(),
		self.Updated.String(),
		self.Flag,
	)
}

func (self *Tag) Encode(w io.Writer) error {
	err := json.NewEncoder(w).Encode(&self)
	return err
}

func (self *Tag) Save() error {
	if len(self.UUID) != 32 {
		return fmt.Errorf("UUID is not 32 bytes: %x", self.UUID)
	}

	if len(self.Label) > 128 {
		return fmt.Errorf("Label is over 128 characters: %s", self.Label)
	}

	statement, err := self.Store.Prepare(`
		INSERT INTO tag (uuid, flag, label)
		VALUES (?, ?, ?);
	`)

	if err != nil {
		return err
	}

	defer statement.Close()
	result, err := statement.Exec(
		self.UUID,
		self.Flag,
		self.Label,
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

func (self *Tag) Delete() error {
	statement, err := self.Store.Prepare(`
		DELETE FROM tag WHERE id = ?;
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

func (self *Tag) ExportMetadata() (int64, string) {
	return self.ID, self.Mapper
}

func (self *Tag) Map(entity Entity) error {
	return fmt.Errorf("Cannot create mapping from %s", self.Mapper)
}

func (self *Tag) Unmap(entity Entity) error {
	return fmt.Errorf("Cannot delete mapping from %s", self.Mapper)
}
