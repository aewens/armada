package main

import (
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/aewens/armada/cargo"
)

func main() {
	path := "test.db"
	//path := "file::memory?cached=shared"
	_, err := cargo.New(path)
	if err != nil {
		log.Fatalln(err)
	}

	//scribe.DB.Create()

	//var entities []*storage.Raw
	//scribe.DB.Find(&entities)
	//for _, entity := range entities {
	//	entity.Display()
	//}
}
