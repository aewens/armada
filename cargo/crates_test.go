package cargo

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/aewens/armada/cargo/model"
)

func catch(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestHold(t *testing.T) {
	hold, err := New(":memory:")
	catch(t, err)

	if hold.Store == nil {
		t.Fatal("Store is nil")
	}

	defer hold.Store.Close()

	_, err = hold.NewCrate("tag")
	if err == nil {
		t.Fatal("New tag crate is not invalid")
	}

	tag, err := hold.NewTag("test")
	catch(t, err)
	if tag == nil {
		t.Fatal("New tag is nil")
	}

	err = tag.Save()
	catch(t, err)

	icrate, err := hold.NewCrate("internal")
	catch(t, err)
	if icrate == nil {
		t.Fatal("New internal crate is nil")
	}

	err = icrate.Set("type", []byte("test"))
	catch(t, err)
	err = icrate.Set("origin", []byte("test"))
	catch(t, err)
	err = icrate.Set("data", []byte{0})
	catch(t, err)

	err = icrate.Save()
	catch(t, err)

	err = icrate.Map(tag)
	catch(t, err)

	internal, ok := icrate.(*model.Internal)
	if !ok {
		t.Fatal("Crate is not Internal")
	}

	if len(internal.Tags) != 1 {
		t.Fatal("Missing tag entry")
	}

	if len(internal.Mapping) != 1 {
		t.Fatal("Missing tag mapping")
	}

	changes := make(map[string][]byte)
	changes["data"] = []byte{1}
	err = icrate.Update(changes)
	catch(t, err)

	err = icrate.Delete()
	catch(t, err)

	ecrate, err := hold.NewCrate("external")
	catch(t, err)
	if ecrate == nil {
		t.Fatal("New external crate is nil")
	}

	err = ecrate.Set("type", []byte("test"))
	catch(t, err)
	err = ecrate.Set("name", []byte("test"))
	catch(t, err)
	err = ecrate.Set("body", []byte("test"))
	catch(t, err)

	err = ecrate.Save()
	catch(t, err)

	err = ecrate.Map(tag)
	catch(t, err)

	external, ok := ecrate.(*model.External)
	if !ok {
		t.Fatal("Crate is not External")
	}

	if len(external.Tags) != 1 {
		t.Fatal("Missing tag entry")
	}

	if len(external.Mapping) != 1 {
		t.Fatal("Missing tag mapping")
	}

	changes = make(map[string][]byte)
	changes["name"] = []byte("changed")
	changes["body"] = []byte("changed")
	err = ecrate.Update(changes)
	catch(t, err)

	err = ecrate.Delete()
	catch(t, err)
}
