package model

import (
	"io"
	"time"
)

type Displayer interface {
	Display()
}

type Encoder interface {
	Encode(io.Writer) error
}

type Setter interface {
	Set(string, []byte) error
}

type Entity interface {
	Displayer
	Encoder
	Writer
	Mapper
}

type Crate interface{
	Entity
	Setter
	Updater
}

type CrateStream chan Crate

type Mapper interface {
	ExportMetadata() (int64, string)
	Map(Entity) error
	Unmap(Entity) error
}

type Saver interface {
	Save() error
}

type Updater interface {
	Update(map[string][]byte) error
}

type Deleter interface {
	Delete() error
}

type Writer interface {
	Saver
	Deleter
}

type Finder interface {
	All() CrateStream
	Lookup(...int) CrateStream
	Contains(string, []byte) CrateStream
	Equals(string, []byte) CrateStream
	Before(time.Time) CrateStream
	After(time.Time) CrateStream
	Between(time.Time, time.Time) CrateStream
}
