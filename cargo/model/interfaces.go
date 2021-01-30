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
	//Updater
	//Deleter
}

type CrateStream chan Crate

type Saver interface {
	Save() error
}

type Updater interface {
	Update() error
}

type Deleter interface {
	Delete() error
}

type Writer interface {
	Saver
	Updater
	Deleter
}

type WhereQuery interface {
	Contains(string) CrateStream
	Equals(int) CrateStream
	Is(string) CrateStream
	Before(time.Time) CrateStream
	After(time.Time) CrateStream
}

type FindQuery interface {
	//Where(string) WhereQuery
	All() CrateStream
}

type Finder interface {
	Find(Crate) FindQuery
}
