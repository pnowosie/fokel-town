package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootOnEmptyTrie(t *testing.T) {
	var (
		emptyLeaf   Trie = (*LeafNode)(nil)
		emptyBranch Trie = (*BranchNode)(nil)
	)

	t.Run("Leaf", func(t *testing.T) {
		assert.Equal(t, Hash{}, emptyLeaf.Root())
	})
	t.Run("Branch", func(t *testing.T) {
		assert.Equal(t, Hash{}, emptyBranch.Root())
	})
}

func makeTestTrie(t *testing.T) Trie {
	trie := &MerkleTrie{}
	assert.NoError(t, trie.Put("001a00", UserData{Id: "001a00"}))
	assert.NoError(t, trie.Put("001a01", UserData{Id: "001a01"}))
	assert.NoError(t, trie.Put("001aba", UserData{Id: "001aba"}))
	assert.NoError(t, trie.Put("001abb", UserData{Id: "001abb"}))

	return trie
}

func TestSearchingOnMerkle(t *testing.T) {
	tests := map[string]struct {
		keyToFind     string
		expectedErr   error
		expectedValue UserData
	}{
		"Successful 001a00": {
			keyToFind:     "001a00",
			expectedValue: UserData{Id: "001a00"},
		},
		"Successful 001a01": {
			keyToFind:     "001a01",
			expectedValue: UserData{Id: "001a01"},
		},
		"Successful 001aba": {
			keyToFind:     "001aba",
			expectedValue: UserData{Id: "001aba"},
		},
		"Successful 001abb": {
			keyToFind:     "001abb",
			expectedValue: UserData{Id: "001abb"},
		},
		"Not found at start": {
			keyToFind:   "f0000f",
			expectedErr: ErrNotFound,
		},
		"Almost found": {
			keyToFind:   "001a0e",
			expectedErr: ErrNotFound,
		},
	}

	trie := makeTestTrie(t)
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			value, err := trie.Get(test.keyToFind)
			if err != nil {
				assert.Equal(t, test.expectedErr, err)
				return
			}
			assert.Equal(t, test.expectedValue, value)
		})
	}
}

func TestRootCalcOnMerkle(t *testing.T) {
	t.Run("Empty trie", func(t *testing.T) {
		expectedHash := "0000000000000000000000000000000000000000000000000000000000000000"
		emptyMerkle := &MerkleTrie{RootNode: nil}
		assert.Equal(t, expectedHash, emptyMerkle.Root().String())
	})

	t.Run("Test trie", func(t *testing.T) {
		expectedHash := "6370101c0992860d2b0b6b9c604ad7fece1728decd34e94817aee12e4c264531"
		trie := makeTestTrie(t)
		assert.Equal(t, expectedHash, trie.Root().String())
	})
}

func TestPutMerkle(t *testing.T) {
	trie := &MerkleTrie{}

	key := "001a00"
	err := trie.Put(key, UserData{Id: key})
	assert.NoError(t, err)

	userData, err := trie.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, UserData{Id: key}, userData)
}
