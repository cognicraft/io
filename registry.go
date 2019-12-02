package io

import (
	"fmt"
	"reflect"
	"sync"
)

func NewTypeRegistry() *TypeRegistry {
	return &TypeRegistry{
		nameToConcreteType: map[string]reflect.Type{},
		concreteTypeToName: map[reflect.Type]string{},
	}
}

type TypeRegistry struct {
	registerLock       sync.RWMutex
	nameToConcreteType map[string]reflect.Type
	concreteTypeToName map[reflect.Type]string
}

func (r *TypeRegistry) Register(name string, value interface{}) error {
	if name == "" {
		return fmt.Errorf("attempt to register empty name")
	}
	r.registerLock.Lock()
	defer r.registerLock.Unlock()
	ut := reflect.TypeOf(value)
	if t, ok := r.nameToConcreteType[name]; ok && t != ut {
		return fmt.Errorf("registering duplicate types for %q: %s != %s", name, t, ut)
	}
	if n, ok := r.concreteTypeToName[ut]; ok && n != name {
		return fmt.Errorf("registering duplicate names for %s: %q != %q", ut, n, name)
	}
	r.nameToConcreteType[name] = ut
	r.concreteTypeToName[ut] = name
	return nil
}

func (r *TypeRegistry) Type(name string) (reflect.Type, error) {
	r.registerLock.RLock()
	defer r.registerLock.RUnlock()
	if t, ok := r.nameToConcreteType[name]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("name not registered: %s", name)
}

func (r *TypeRegistry) Name(value interface{}) (string, error) {
	r.registerLock.RLock()
	defer r.registerLock.RUnlock()
	ut := reflect.TypeOf(value)
	if n, ok := r.concreteTypeToName[ut]; ok {
		return n, nil
	}
	return "", fmt.Errorf("type not registered: %T", value)
}

func (r *TypeRegistry) Marshal(marshal MarshalStrategy, value interface{}) (string, []byte, error) {
	typeName, err := r.Name(value)
	if err != nil {
		return "", nil, err
	}
	data, err := marshal(value)
	if err != nil {
		return "", nil, err
	}
	return typeName, data, nil
}

func (r *TypeRegistry) Unmarshal(unmarshal UnmarshalStrategy, name string, data []byte) (interface{}, error) {
	typ, err := r.Type(name)
	if err != nil {
		return nil, err
	}
	pointer := reflect.New(typ)
	err = unmarshal(data, pointer.Interface())
	return pointer.Elem().Interface(), err
}

type MarshalStrategy func(interface{}) ([]byte, error)

type UnmarshalStrategy func([]byte, interface{}) error
