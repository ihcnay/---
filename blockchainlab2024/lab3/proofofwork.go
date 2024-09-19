package main

import (
	"bytes"
	"crypto/sha256"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 8

// ProofOfWork represents a proof-of-work
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork builds and returns a ProofOfWork
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

// Run performs a proof-of-work
// implement
func (pow *ProofOfWork) Run() (int64, []byte) {
	nonce := int64(0)
	pow.block.SetNonce(nonce)
	var hash [32]byte
	for nonce < int64(maxNonce) {
		buffer := bytes.Buffer{}
		buffer.Write(IntToHex(pow.block.Header.Version))
		buffer.Write(pow.block.Header.PrevBlockHash[:])
		buffer.Write(pow.block.Header.MerkleRoot[:])
		buffer.Write(IntToHex(pow.block.Header.Timestamp))
		buffer.Write(IntToHex(int64(targetBits)))
		buffer.Write(IntToHex(pow.block.Header.Nonce))
		header := buffer.Bytes()
		hash = sha256.Sum256(header)
		var num big.Int
		num.SetBytes(hash[:])
		if num.Cmp(pow.target) < 0 {
			break
		} else {
			nonce++
			pow.block.SetNonce(nonce)
		}
	}
	return nonce, nil
}

// Validate validates block's PoW
// implement
func (pow *ProofOfWork) Validate() bool {
	buffer := bytes.Buffer{}
	buffer.Write(IntToHex(pow.block.Header.Version))
	buffer.Write(pow.block.Header.PrevBlockHash[:])
	buffer.Write(pow.block.Header.MerkleRoot[:])
	buffer.Write(IntToHex(pow.block.Header.Timestamp))
	buffer.Write(IntToHex(int64(targetBits)))
	buffer.Write(IntToHex(pow.block.Header.Nonce))
	header := buffer.Bytes()
	hash := sha256.Sum256(header)
	var num big.Int
	num.SetBytes(hash[:])
	return num.Cmp(pow.target) < 0
}
