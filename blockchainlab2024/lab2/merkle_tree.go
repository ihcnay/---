package main

import (
	"bytes"
	"crypto/sha256"
)

// MerkleTree represent a Merkle tree
type MerkleTree struct {
	RootNode *MerkleNode
	Leaf     [][]byte
}

// MerkleNode represent a Merkle tree node
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

// NewMerkleTree creates a new Merkle tree from a sequence of data
// implement
func NewMerkleTree(data [][]byte) *MerkleTree {
	length := len(data)
	if length%2 != 0 {
		data = append(data, data[length-1])
	}
	var merklenodes []*MerkleNode
	i := 0
	for i < len(data) {
		merklenodes = append(merklenodes, NewMerkleNode(nil, nil, data[i]))
		i++
	}
	for len(merklenodes) > 1 {
		var tnodes []*MerkleNode
		length = len(merklenodes)
		if length%2 != 0 {
			merklenodes = append(merklenodes, merklenodes[length-1])
			length = length + 1
		}
		t := 0
		for t < length {
			tnodes = append(tnodes, NewMerkleNode(merklenodes[t], merklenodes[t+1], nil))
			t = t + 2
		}
		merklenodes = tnodes
	}

	return &MerkleTree{
		merklenodes[0],
		data,
	}
}

// NewMerkleNode creates a new Merkle tree node
// implement
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	if data != nil {
		hash := sha256.Sum256(data)
		return &MerkleNode{
			left,
			right,
			hash[:],
		}
	} else {
		hash := sha256.Sum256(append(left.Data, right.Data...))
		return &MerkleNode{
			left,
			right,
			hash[:],
		}
	}
}

func (t *MerkleTree) SPVproof(index int) ([][]byte, error) {
	n := len(t.Leaf)
	if index >= n {
		return nil, nil
	}
	var cpath []rune
	ptr := t.RootNode.Left
	for ptr != nil {
		if index%2 == 1 {
			cpath = append(cpath, 'l')
			index = (index - 1) / 2
		} else {
			cpath = append(cpath, 'r')
			index = index / 2
		}
		ptr = ptr.Left
	}
	ptr = t.RootNode
	var path [][]byte
	length := len(cpath)
	for length > 0 {
		if cpath[length-1] == 'r' {
			path = append(path, ptr.Right.Data)
			ptr = ptr.Left
		} else {
			path = append(path, ptr.Left.Data)
			ptr = ptr.Right
		}
		length--
	}
	return path, nil
}

func (t *MerkleTree) VerifyProof(index int, path [][]byte) (bool, error) {
	data := t.Leaf[index]
	hashdata := sha256.Sum256(data)
	n := len(path)
	for n > 0 {
		if index%2 == 1 {
			hashdata = sha256.Sum256(append(path[n-1], hashdata[:]...))
			index = (index - 1) / 2
		} else {
			hashdata = sha256.Sum256(append(hashdata[:], path[n-1]...))
			index = index / 2
		}
		n--
	}
	if bytes.Equal(hashdata[:], t.RootNode.Data) {
		return true, nil
	}
	return false, nil
}
