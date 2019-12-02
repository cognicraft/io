package io

import (
	"encoding/json"
	"testing"
)

type Foo struct {
	Bar string
}

type Bar struct {
	Nil string
}

func TestTypeRegistry(t *testing.T) {
	r := NewTypeRegistry()

	err := r.Register("", Foo{})
	if err == nil {
		t.Errorf("expected an error")
	}

	err = r.Register("foo", Foo{})
	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}
	err = r.Register("o-foo", Foo{})
	if err == nil {
		t.Errorf("expected an error")
	}
	err = r.Register("foo", Bar{})
	if err == nil {
		t.Errorf("expected an error")
	}

	name, data, err := r.Marshal(json.Marshal, Bar{})
	if err == nil {
		t.Errorf("expected an error")
	}

	a := Foo{Bar: "Baz"}
	name, data, err = r.Marshal(json.Marshal, a)
	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}
	if name != "foo" {
		t.Errorf("want: %s, got: %s", "foo", name)
	}

	v, err := r.Unmarshal(json.Unmarshal, "n-a", data)
	if err == nil {
		t.Errorf("expected an error")
	}
	v, err = r.Unmarshal(json.Unmarshal, name, data)
	if err != nil {
		t.Errorf("expected no error, but got: %v", err)
	}
	if a != v {
		t.Errorf("want: %#v, got: %#v", a, v)
	}
}
