package internal

import (
	"crypto"
	"strings"
)

var _ Trie = (*MerkleTrie)(nil)

type MerkleTrie struct {
	RootNode *BranchNode
}

func (m MerkleTrie) Root() Hash {
	return m.RootNode.Root()
}

func (m MerkleTrie) Get(key string) (UserData, error) {
	return m.RootNode.Get(key)
}

func (m MerkleTrie) Put(key string, value UserData) error {
	//TODO implement me
	panic("implement me")
}

type LeafNode struct {
	Key   string
	Value UserData
}

func (ln *LeafNode) Root() Hash {
	if ln == nil {
		return Hash{}
	}
	return ln.Value.Hash()
}

func (ln *LeafNode) Get(key string) (UserData, error) {
	if ln.Key == key {
		return ln.Value, nil
	}
	return UserData{}, ErrNotFound
}

func (ln *LeafNode) Put(key string, value UserData) error {
	//TODO implement me
	panic("implement me")
}

type BranchNode struct {
	PathPrefix string
	Children   [16]Trie
}

func (bn *BranchNode) Root() Hash {
	if bn == nil {
		return Hash{}
	}

	emptyValueHash := Hash{}.Bytes()
	hasher := crypto.SHA256.New()
	for _, child := range bn.Children {
		if child == nil {
			hasher.Write(emptyValueHash)
		} else {
			hasher.Write(child.Root().Bytes())
		}
	}

	result := Hash{}
	copy(result[:], hasher.Sum(nil))
	return result
}

func (bn *BranchNode) Get(key string) (UserData, error) {
	if key == bn.PathPrefix {
		panic("sanity check: this supposed to be a leaf node")
	}

	if !strings.HasPrefix(key, bn.PathPrefix) {
		return UserData{}, ErrNotFound
	}

	childIdx := nibbleIndex(key[len(bn.PathPrefix)])
	child := bn.Children[childIdx]
	if child == nil {
		return UserData{}, ErrNotFound
	}

	return child.Get(key)
}

func (bn *BranchNode) Put(key string, value UserData) error {
	//TODO implement me
	panic("implement me")
}

func nibbleIndex(nibble byte) int {
	if nibble >= byte('0') && nibble <= byte('9') {
		return int(nibble - byte('0'))
	}
	if nibble >= byte('a') && nibble <= byte('f') {
		return 10 + int(nibble-byte('a'))
	}
	panic("invalid nibble")
}
