package internal

import (
	"crypto"
	"math"
	"strings"
)

var _ Trie = (*MerkleTrie)(nil)

type MerkleTrie struct {
	RootNode  *BranchNode
	elemCount int
}

func (m *MerkleTrie) Root() Hash {
	return m.RootNode.Root()
}

func (m *MerkleTrie) Get(key string) (UserData, error) {
	return m.RootNode.Get(key)
}

func (m *MerkleTrie) Put(key string, value UserData) error {
	if m.RootNode == nil {
		m.RootNode = &BranchNode{Children: [16]merkleTreeNode{}}
	}

	if err := m.RootNode.Put(key, value); err != nil {
		return err
	}

	m.elemCount++
	return nil
}

func (m *MerkleTrie) Count() int {
	return m.elemCount
}

type merkleTreeNode interface {
	Trie
	putNode(key string, value UserData) (merkleTreeNode, error)
}

var _ merkleTreeNode = (*LeafNode)(nil)

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

func (ln *LeafNode) putNode(key string, value UserData) (merkleTreeNode, error) {
	// No node here, create a new leaf node
	if ln == nil {
		return &LeafNode{Key: key, Value: value}, nil
	}

	if ln.Key == key {
		return nil, ErrAlreadyExists
	}

	prefix, i := longestCommonPrefix(key, ln.Key)
	newBranch := &BranchNode{PathPrefix: prefix, Children: [16]merkleTreeNode{}}
	newBranch.Children[nibbleIndex(ln.Key[i])] = ln
	newBranch.Children[nibbleIndex(key[i])] = &LeafNode{Key: key, Value: value}

	return newBranch, nil
}

func longestCommonPrefix(first, second string) (string, int) {
	shortestLen := int(math.Min(float64(len(first)), float64(len(second))))

	var i = 0
	for ; i < shortestLen; i++ {
		if first[i] != second[i] {
			break
		}
	}
	return first[:i], i
}

var _ merkleTreeNode = (*BranchNode)(nil)

type BranchNode struct {
	PathPrefix string
	Children   [16]merkleTreeNode
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
	if bn == nil {
		return UserData{}, ErrNotFound
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
	if bn == nil {
		panic("sanity check: BranchNode.Put never called on nil")
	}

	childIdx := nibbleIndex(key[len(bn.PathPrefix)])
	child := bn.Children[childIdx]
	if child == nil {
		bn.Children[childIdx] = &LeafNode{Key: key, Value: value}
		return nil
	}

	newChild, err := child.putNode(key, value)
	if err != nil {
		return err
	}
	bn.Children[childIdx] = newChild
	return nil
}

func (bn *BranchNode) putNode(key string, value UserData) (merkleTreeNode, error) {
	if bn == nil {
		return (*LeafNode)(nil).putNode(key, value)
	}

	if !strings.HasPrefix(key, bn.PathPrefix) {
		prefix, i := longestCommonPrefix(key, bn.PathPrefix)
		newBranch := &BranchNode{PathPrefix: prefix, Children: [16]merkleTreeNode{}}
		newBranch.Children[nibbleIndex(bn.PathPrefix[i])] = bn

		return newBranch.putNode(key, value)
	}

	if err := bn.Put(key, value); err != nil {
		return nil, err
	}
	return bn, nil
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
