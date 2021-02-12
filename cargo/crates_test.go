package cargo

import (
	"testing"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/aewens/nautical/cargo/model"
	"github.com/aewens/nautical/cargo/repo"
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

	tag, err := hold.NewTag()
	catch(t, err)
	if tag == nil {
		t.Fatal("New tag is nil")
	}

	err = tag.Set("label", []byte("test"))
	catch(t, err)

	err = tag.Save()
	catch(t, err)

	err = tag.Set("flag", []byte{2})
	catch(t, err)

	err = tag.Update()
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

	err = icrate.Unmap(tag)
	catch(t, err)

	if len(internal.Tags) != 0 {
		t.Fatal("Did not remove tag entry")
	}

	if len(internal.Mapping) != 0 {
		t.Fatal("Did not remove tag mapping")
	}

	err = icrate.Set("data", []byte{1})
	catch(t, err)

	err = icrate.Update()
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

	err = ecrate.Unmap(tag)
	catch(t, err)

	if len(external.Tags) != 0 {
		t.Fatal("Did not remove tag entry")
	}

	if len(external.Mapping) != 0 {
		t.Fatal("Did not remove tag mapping")
	}

	err = ecrate.Set("name", []byte("changed"))
	catch(t, err)

	err = ecrate.Set("body", []byte("changed"))
	catch(t, err)

	err = ecrate.Update()
	catch(t, err)

	err = external.Link(icrate)
	catch(t, err)

	if external.Meta.ID != internal.ID {
		t.Fatal("Failed to link internal to meta")
	}

	if len(external.Data) == 0 {
		t.Fatal("Failed to link UUID to data")
	}

	err = external.Unlink()
	catch(t, err)

	if external.Meta != nil {
		t.Fatal("Failed to unlink meta")
	}

	if len(external.Data) != 0 {
		t.Fatal("Failed to unlink UUID")
	}

	err = ecrate.Delete()
	catch(t, err)

	err = icrate.Delete()
	catch(t, err)

	err = tag.Delete()
	catch(t, err)
}

func StreamSize(stream repo.Stream) int {
	count := 0
	for range stream {
		count = count + 1
	}
	return count
}

func TestRepos(t *testing.T) {
	now := Now()
	hold, err := New(":memory:")
	catch(t, err)

	if hold.Store == nil {
		t.Fatal("Store is nil")
	}

	defer hold.Store.Close()

	irepo, err := hold.NewRepo("internal")
	catch(t, err)

	ic := 3
	for i := 0; i < ic; i++ {
		entity, err := irepo.Create()
		catch(t, err)
		err = entity.Set("type", []byte(fmt.Sprintf("test%d", i)))
		catch(t, err)
		err = entity.Set("origin", []byte(fmt.Sprintf("test%d", i)))
		catch(t, err)
		err = entity.Set("data", []byte{byte(i)})
		catch(t, err)
		err = entity.Save()
		catch(t, err)
	}

	istream := irepo.All()
	irepo.Load(istream)

	iirepo, ok := irepo.(*repo.Internal)
	if !ok {
		t.Fatal("Could not cast to Internal")
	}

	if len(iirepo.Crates) != ic {
		t.Fatal("Did not load all entities")
	}

	count := StreamSize(irepo.Lookup(2))
	if count != 1 {
		t.Fatal("Could not lookup entity")
	}

	count = StreamSize(irepo.Contains("origin", "1"))
	if count != 1 {
		t.Fatal("Could not lookup entity")
	}

	count = StreamSize(irepo.Equals("origin", "test1"))
	if count != 1 {
		t.Fatal("Could not lookup entity")
	}

	count = StreamSize(irepo.Equals("origin", "test1"))
	if count != 1 {
		t.Fatal("Could not lookup entity")
	}

	count = StreamSize(irepo.Before("added", now.Add(1 * time.Minute)))
	if count != ic {
		t.Fatal("Could not lookup entities")
	}

	count = StreamSize(irepo.After("added", now.Add(-1 * time.Minute)))
	if count != ic {
		t.Fatal("Could not lookup entities")
	}

	count = StreamSize(irepo.Between(
		"added",
		now.Add(-1 * time.Minute),
		now.Add(1 * time.Minute),
	))
	if count != ic {
		t.Fatal("Could not lookup entities")
	}

	erepo, err := hold.NewRepo("external")
	catch(t, err)

	ec := 3
	for i := 0; i < ec; i++ {
		entity, err := erepo.Create()
		catch(t, err)
		err = entity.Set("type", []byte(fmt.Sprintf("test%d", i)))
		catch(t, err)
		err = entity.Set("name", []byte(fmt.Sprintf("test%d", i)))
		catch(t, err)
		err = entity.Set("body", []byte(fmt.Sprintf("body%d", i)))
		catch(t, err)
		err = entity.Save()
		catch(t, err)
	}

	estream := erepo.All()
	erepo.Load(estream)

	eerepo, ok := erepo.(*repo.External)
	if !ok {
		t.Fatal("Could not cast to External")
	}

	if len(eerepo.Crates) != ec {
		t.Fatalf("Did not load all entities: %d", len(eerepo.Crates))
	}

	count = StreamSize(erepo.Lookup(2))
	if count != 1 {
		t.Fatal("Could not lookup entity")
	}

	count = StreamSize(erepo.Contains("name", "1"))
	if count != 1 {
		t.Fatal("Could not lookup entity")
	}

	count = StreamSize(erepo.Equals("name", "test1"))
	if count != 1 {
		t.Fatal("Could not lookup entity")
	}

	count = StreamSize(erepo.Before("added", now.Add(1 * time.Minute)))
	if count != ec {
		t.Fatal("Could not lookup entities")
	}

	count = StreamSize(erepo.After("added", now.Add(-1 * time.Minute)))
	if count != ec {
		t.Fatal("Could not lookup entities")
	}

	count = StreamSize(erepo.Between(
		"added",
		now.Add(-1 * time.Minute),
		now.Add(1 * time.Minute),
	))
	if count != ec {
		t.Fatal("Could not lookup entities")
	}

	trepo, err := hold.NewRepo("tag")
	catch(t, err)

	tc := 3
	for i := 0; i < tc; i++ {
		entity, err := trepo.Create()
		catch(t, err)
		err = entity.Set("label", []byte(fmt.Sprintf("test%d", i)))
		catch(t, err)
		err = entity.Save()
		catch(t, err)
	}

	tstream := trepo.All()
	trepo.Load(tstream)

	ttrepo, ok := trepo.(*repo.Tag)
	if !ok {
		t.Fatal("Could not cast to Tag")
	}

	if len(ttrepo.Crates) != tc {
		t.Fatal("Did not load all entities")
	}

	count = StreamSize(trepo.Lookup(2))
	if count != 1 {
		t.Fatal("Could not lookup entity")
	}

	count = StreamSize(trepo.Contains("label", "1"))
	if count != 1 {
		t.Fatal("Could not lookup entity")
	}

	count = StreamSize(trepo.Equals("label", "test1"))
	if count != 1 {
		t.Fatal("Could not lookup entity")
	}

	count = StreamSize(trepo.Before("added", now.Add(1 * time.Minute)))
	if count != tc {
		t.Fatal("Could not lookup entities")
	}

	count = StreamSize(trepo.After("added", now.Add(-1 * time.Minute)))
	if count != tc {
		t.Fatal("Could not lookup entities")
	}

	count = StreamSize(trepo.Between(
		"added",
		now.Add(-1 * time.Minute),
		now.Add(1 * time.Minute),
	))
	if count != tc {
		t.Fatal("Could not lookup entities")
	}
}
