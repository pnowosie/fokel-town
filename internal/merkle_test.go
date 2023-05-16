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

	assert.Equal(t, Hash{}, emptyLeaf.Root())
	assert.Equal(t, Hash{}, emptyBranch.Root())
}

func MakeTrieByHand() Trie {
	_ = `
	// keys := [
	// 	"001a00",
	// 	"001a01",
	// 	"001aba",
	// 	"001abb" ]

	// Prefix "001ab"
	var _001ab = &BranchNode{
		PathPrefix: "001ab",
		Children:   [16]merkleTreeNode{},
	}
	_001ab.Children[10] = &LeafNode{
		Key: "001aba", Value: UserData{Id: "001aba"},
	}
	_001ab.Children[11] = &LeafNode{
		Key: "001abb", Value: UserData{Id: "001abb"},
	}

	// Prefix "001ab"
	var _001a0 = &BranchNode{
		PathPrefix: "001a0",
		Children:   [16]merkleTreeNode{},
	}
	_001a0.Children[0] = &LeafNode{
		Key: "001a00", Value: UserData{Id: "001a00"},
	}
	_001a0.Children[1] = &LeafNode{
		Key: "001a01", Value: UserData{Id: "001a01"},
	}

	// Prefix "001a"
	var _001a = &BranchNode{
		PathPrefix: "001a",
		Children:   [16]merkleTreeNode{},
	}
	_001a.Children[0] = _001a0
	_001a.Children[11] = _001ab

	var root = &BranchNode{
		PathPrefix: "",
		Children:   [16]merkleTreeNode{_001a},
	}
`
	trie := &MerkleTrie{}
	trie.Put("001a00", UserData{Id: "001a00"})
	trie.Put("001a01", UserData{Id: "001a01"})
	trie.Put("001aba", UserData{Id: "001aba"})
	trie.Put("001abb", UserData{Id: "001abb"})

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

	trie := MakeTrieByHand()
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
	expectedHash := "0000000000000000000000000000000000000000000000000000000000000000"
	emptyMerkle := &MerkleTrie{RootNode: nil}
	assert.Equal(t, expectedHash, emptyMerkle.Root().String())

	expectedHash = "6370101c0992860d2b0b6b9c604ad7fece1728decd34e94817aee12e4c264531"
	trie := MakeTrieByHand()
	assert.Equal(t, expectedHash, trie.Root().String())
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
