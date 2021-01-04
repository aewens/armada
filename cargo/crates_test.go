package cargo

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
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
	icrate, err := hold.NewCrate("internal")
	catch(t, err)
	if icrate == nil {
		t.Fatal("New internal crate is nil")
	}

	err = icrate.Set("type", "test")
	catch(t, err)
	err = icrate.Set("origin", "test")
	catch(t, err)
	err = icrate.Set("data", []byte{0})
	catch(t, err)

	err = Save(hold.Store, icrate)
	catch(t, err)

	ecrate, err := hold.NewCrate("external")
	catch(t, err)
	if ecrate == nil {
		t.Fatal("New external crate is nil")
	}

	err = ecrate.Set("type", "test")
	catch(t, err)
	err = ecrate.Set("name", "test")
	catch(t, err)
	err = ecrate.Set("body", "test")
	catch(t, err)

	err = Save(hold.Store, ecrate)
	catch(t, err)

	_, err = hold.NewCrate("tag")
	if err == nil {
		t.Fatal("New tag crate is not invalid")
	}
}
