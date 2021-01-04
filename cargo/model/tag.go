package model

import (
	"io"
	"fmt"
	"log"
	"encoding/json"
)

type Tag struct {
	Common
	Label  string `json:"label"`
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
