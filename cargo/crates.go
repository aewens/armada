package cargo

import (
	"fmt"
	"time"
	"database/sql"

	"github.com/aewens/nautical/cargo/model"
	"github.com/aewens/nautical/cargo/repo"
)

type Hold struct {
	Store  *sql.DB
}

func Now() time.Time {
	return model.Now()
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

func (self *Hold) NewCrate(crateType string) (model.Entity, error) {
	var crate model.Entity = nil

	switch crateType {
	case "internal":
		return model.NewInternal(self.Store)
	case "external":
		return model.NewExternal(self.Store)
	}
	return crate, fmt.Errorf("Invalid crate type: %s", crateType)
}

func (self *Hold) NewTag() (*model.Tag, error) {
	return model.NewTag(self.Store)
}

func (self *Hold) NewRepo(repoType string) (repo.Entity, error) {
	var repository repo.Entity = nil

	switch repoType {
	case "internal":
		return repo.NewInternal(self.Store), nil
	case "external":
		return repo.NewExternal(self.Store), nil
	case "tag":
		return repo.NewTag(self.Store), nil
	}

	return repository, fmt.Errorf("Invalid repo type: %s", repoType)
}
