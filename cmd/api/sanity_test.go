package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDummyStorage(t *testing.T) {
	m1 := MapIsNotATrie{}

	_, err := m1.Get("00aa00")
	assert.Equal(t, errNotFound, err)

	user := UserData{"00aa00", "John", "Doe"}
	err = m1.Put(user.Id, user)
	assert.NoError(t, err)
	userFromStorage, err := m1.Get(user.Id)
	assert.Equal(t, user, userFromStorage)
}

func TestUserDataValidation(t *testing.T) {
	validUser := UserData{"00aa00", "John", "Doe"}
	assert.True(t, validUser.IsValid())

	invalidShorter := UserData{"aa00", "John", "Doe"}
	assert.False(t, invalidShorter.IsValid())

	invalidLonger := UserData{"1234567", "John", "Doe"}
	assert.False(t, invalidLonger.IsValid())
}

func TestUserHashCalculation(t *testing.T) {
	user1 := UserData{"00aa00", "John", "Doe"}
	user2 := UserData{"00aa01", "John", "Deer"}

	assert.Equal(t, "7156f64e8fea558be429b43bde862f837e69bbe046abb440a07a7936522db3ff", user1.Hash().String())
	assert.Equal(t, "c69bf8bac2352f4e7eeda1de54d6bb143c80539d8abcd80deae81cad434e18ef", user2.Hash().String())
}

func TestThreadSafeTrieAddMany(t *testing.T) {
	trie := ThreadSafeTrie{Trie: &MapIsNotATrie{}}
	var waitGroup sync.WaitGroup

	for i := 0; i < 3; i++ {
		waitGroup.Add(1)
		go func(idx int) {
			defer waitGroup.Done()
			for j := 0; j < 1000; j++ {
				id := fmt.Sprintf("%d%05d", idx, j)
				user := UserData{Id: id}
				err := trie.Put(user.Id, user)
				assert.NoError(t, err)
			}
		}(i)
	}

	waitGroup.Wait()

	assert.Equal(t, 3*1000, len(trie.Trie.(*MapIsNotATrie).data))

	for i := 0; i < 3; i++ {
		for _, idx := range []int{0, 501, 999} {
			id := fmt.Sprintf("%d%05d", i, idx)
			user, err := trie.Get(id)
			assert.NoError(t, err)
			assert.Equal(t, id, user.Id)
		}
	}
}
