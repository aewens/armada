package repo

import (
	//"time"

	"github.com/aewens/armada/cargo/model"
)

type Stream chan model.Entity

type Reader interface {
	All() Stream
	Get(int64) (model.Entity, error)
	Lookup(...int64) Stream
	Contains(string, string) Stream
	Equals(string, string) Stream
	//Before(time.Time) Stream
	//After(time.Time) Stream
	//Between(time.Time, time.Time) Stream
}

type Entity interface {
	Reader
	Create() (model.Entity, error)
	Load(Stream)
}
