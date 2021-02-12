package repo

import (
	"time"

	"github.com/aewens/nautical/cargo/model"
)

type Stream chan model.Entity

type Reader interface {
	All() Stream
	Get(int64) (model.Entity, error)
	Lookup(...int64) Stream
	Contains(string, string) Stream
	Equals(string, string) Stream
	Before(string, time.Time) Stream
	After(string, time.Time) Stream
	Between(string, time.Time, time.Time) Stream
}

type Entity interface {
	Reader
	Create() (model.Entity, error)
	Load(Stream)
}
