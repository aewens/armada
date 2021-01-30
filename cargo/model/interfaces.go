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
	Set(key string, value []byte) error
}

type Crate interface{
	Displayer
	Encoder
	Setter
	Saver
	Updater
	Deleter
}

type CrateStream chan Crate

type Saver interface {
	Save() error
}

type Updater interface {
	Update(changes map[string][]byte) error
}

type Deleter interface {
	Delete() error
}

type Writer interface {
	Saver
	Updater
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
