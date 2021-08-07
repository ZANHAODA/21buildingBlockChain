package main

import (
	"bytes"
	"crypto/sha256"
	"math/big"
)

type ProofOfWork struct {
	block *Block
	target *big.Int
}
const targetBits = 16 //挖矿难度

func NewProofofWork(b*Block) *ProofOfWork{
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - targetBits))
	pow := &ProofOfWork{b, target}
	return pow
}

func (pow * ProofOfWork) prepareData(nonce int32) []byte {
	data := bytes.Join(
		[][]byte{
			IntToHex(pow.block.Version),
			pow.block.PrevBlockHash,
			pow.block.Merkleroot,
			IntToHex(pow.block.Time),
			IntToHex(pow.block.Bits),
			IntToHex(pow.block.Nonce)},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Run() (int32, []byte){
	var nonce int32
	var secondhash [32]byte
	nonce = 0
	var currenthash big.Int
	for nonce < maxnonce{
		data:= pow.prepareData(nonce)  //xuliehua

		firstHash:= sha256.Sum256(data)
		secondhash = sha256.Sum256(firstHash[:])
//		fmt.Printf("%x\n", secondhash)
		currenthash.SetBytes(secondhash[:])

		if currenthash.Cmp(pow.target) == -1 {
			break
		} else{
			nonce++
		}
	}
	return nonce,secondhash[:]
}


func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.block.Nonce)
	firstHash := sha256.Sum256(data)
	secondhash := sha256.Sum256(firstHash[:])
	hashInt.SetBytes(secondhash[:])
	isValid:= hashInt.Cmp(pow.target) == -1
	return isValid
}
