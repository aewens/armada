package model

import (
	"io"
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

type Mapper interface {
	ExportMetadata() (int64, string)
	Map(Entity) error
	Unmap(Entity) error
}

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

type Entity interface {
	Displayer
	Encoder
	Setter
	Writer
	Mapper
}
