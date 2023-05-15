package main

import (
	"errors"
)

var (
	errNotFound     = errors.New("key not found")
	errAlreadyExist = errors.New("key already exists")
)

type Trie interface {
	Root() Hash
	Get(key string) (UserData, error)
	Put(key string, value UserData) error
}

var _ Trie = (*MapIsNotATrie)(nil)

// MapIsNotATrie is a dummy implementation of Trie interface.
type MapIsNotATrie struct {
	data map[string]UserData
}

func (m *MapIsNotATrie) Root() Hash {
	return Hash{}
}

func (m *MapIsNotATrie) Get(key string) (UserData, error) {
	data, ok := m.data[key]
	if !ok {
		return UserData{}, errNotFound
	}
	return data, nil
}

func (m *MapIsNotATrie) Put(key string, value UserData) error {
	if m.data == nil {
		m.data = map[string]UserData{}
	}

	_, ok := m.data[key]
	if ok {
		return errAlreadyExist
	}
	m.data[key] = value
	return nil
}
