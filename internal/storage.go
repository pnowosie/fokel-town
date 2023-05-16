package internal

import (
	"errors"
	"sync"
)

var (
	ErrNotFound     = errors.New("key not found")
	ErrAlreadyExist = errors.New("key already exists")
)

type Trie interface {
	Root() Hash
	Get(key string) (UserData, error)
	Put(key string, value UserData) error
}

var _ Trie = (*ThreadSafeTrie)(nil)

// ThreadSafeTrie implements thread synchronization around in-memory trie.
type ThreadSafeTrie struct {
	Trie
	lock sync.Mutex
}

func (t *ThreadSafeTrie) Root() Hash {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.Trie.Root()
}

func (t *ThreadSafeTrie) Get(key string) (UserData, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.Trie.Get(key)
}

func (t *ThreadSafeTrie) Put(key string, value UserData) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.Trie.Put(key, value)
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
		return UserData{}, ErrNotFound
	}
	return data, nil
}

func (m *MapIsNotATrie) Put(key string, value UserData) error {
	if m.data == nil {
		m.data = map[string]UserData{}
	}

	_, ok := m.data[key]
	if ok {
		return ErrAlreadyExist
	}
	m.data[key] = value
	return nil
}
