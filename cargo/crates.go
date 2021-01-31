package cargo

import (
	"fmt"
	"database/sql"

	"github.com/aewens/armada/cargo/model"
)

type Hold struct {
	Store  *sql.DB
}

func New(conn string) (*Hold, error) {
	var hold *Hold

	store, err := Open(conn, Tables)
	if err != nil {
		return hold, err
	}

	hold = &Hold{
		Store:  store,
	}

	return hold, nil
}

func (self *Hold) NewCrate(crateType string) (model.Crate, error) {
	var crate model.Crate = nil

	switch crateType {
	case "internal":
		return model.NewInternal(self.Store)
	case "external":
		return model.NewExternal(self.Store)
	}
	return crate, fmt.Errorf("Invalid crate type: %s", crateType)
}

func (self *Hold) NewTag(label string) (*model.Tag, error) {
	return model.NewTag(self.Store, label)
}
