package model

import (
	"io"
	"log"
	"encoding/json"
)

type Mapping struct {
	ID       uint   `json:"-"`
	Internal []byte `json:"internal"`
	External []byte `json:"external"`
	Tag      []byte `json:"tag"`
}

func (self *Mapping) Display() {
	InternalUUID := self.Internal
	if len(InternalUUID) > 0 {
		InternalUUID = InternalUUID[:8]
	}

	ExternalUUID := self.External
	if len(ExternalUUID) > 0 {
		ExternalUUID = ExternalUUID[:8]
	}

	TagUUID := self.Tag
	if len(TagUUID) > 0 {
		TagUUID = TagUUID[:8]
	}


	log.Printf(
		"Mapping<gid:%d int:%x ext:%x tag:%x>\n",
		self.ID,
		InternalUUID,
		ExternalUUID,
		TagUUID,
	)
}

func (self *Mapping) Encode(w io.Writer) error {
	err := json.NewEncoder(w).Encode(&self)
	return err
}
