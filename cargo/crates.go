package cargo

import (
	"fmt"
	"database/sql"

	"github.com/aewens/armada/cargo/model"
)

type Manifest []*model.Crate

type Hold struct {
	Store   *sql.DB
	Crates  map[string]Manifest
	Mapping map[string]Manifest
}

func New(conn string) (*Hold, error) {
	var hold *Hold

	store, err := Open(conn, Tables)
	if err != nil {
		return hold, err
	}

	hold = &Hold{
		Store:   store,
		Crates:  make(map[string]Manifest),
		Mapping: make(map[string]Manifest),
	}

	return hold, nil
}

func (self *Hold) NewCrate(crateType string) (model.Crate, error) {
	var crate model.Crate = nil

	switch crateType {
	case "internal":
		return model.NewInternal()
	case "external":
		return model.NewExternal()
	}

	return crate, fmt.Errorf("Invalid crate type: %s", crateType)
}
