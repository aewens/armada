package repo

import (
	//"fmt"
	//"time"
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
