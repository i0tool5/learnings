package reflectnames

import (
	"reflect"
	"strings"
	"testing"
)

type callable struct{}

func (c callable) funcOne() {}
func (c callable) funcTwo() {}
func (c callable) FuncOne() {}
func (c callable) FuncTwo() {}

func TestGetMethodNames(t *testing.T) {
	result := GetMethodNames(&callable{})
	want := []string{"FuncOne", "FuncTwo"}
	if len(result) != 2 {
		t.Fatalf("got %v names, but want 2", len(result))
	}
	if !reflect.DeepEqual(result, want) {
		t.Fatalf("want %v got %v", want, result)
	}
}

func TestWrapping(t *testing.T) {
	var r string
	s := Some{}
	ds := NewDecorated(s)
	r = ds.CallFnOne()
	t.Log(r)
	if !strings.HasPrefix(r, "Calling decorated") {
		t.Fatal("decoration prefix wasn't found")
	}

	r = ds.CallFnTwo()
	if !strings.HasPrefix(r, "Calling decorated") {
		t.Fatal("decoration prefix wasn't found")
	}
	r = ds.CallFnThree()
	if !strings.HasPrefix(r, "Calling decorated") {
		t.Fatal("decoration prefix wasn't found")
	}
}
