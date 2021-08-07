package main

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"strconv"
	"time"
)

var (
	maxnonce int32 = math.MaxInt32
)

type Block struct {
	Version int32
	PrevBlockHash []byte
	Merkleroot []byte
	Hash []byte
	Time int32
	Bits int32 //难度
	Nonce int32  //随机数
	Transations []*Transation
	Height int32
}

func (block *Block) serialize() []byte{

	result := bytes.Join(
		[][]byte{
			IntToHex(block.Version),
			block.PrevBlockHash,
			block.Merkleroot,
			IntToHex(block.Time),
			IntToHex(block.Bits),
			IntToHex(block.Nonce)},
		[]byte{},
		)
	return result
}

func(b* Block) Serialize() []byte{
	var encoded bytes.Buffer
	enc:= gob.NewEncoder(&encoded)

	err:= enc.Encode(b)
	if err != nil{
		log.Panic(err)
	}
	return encoded.Bytes()
}

func DeserializeBlock(d []byte) *Block{
	var block Block
	decode := gob.NewDecoder(bytes.NewReader(d))
	err := decode.Decode(&block)
	if err != nil{
		log.Panic(err)
	}
	return &block
}
//计算困难度
func CalculateTargetFast(bits []byte) []byte {
	var result []byte
	exponent := bits[:1]
	fmt.Printf("%x\n", exponent)

	coeffient := bits[1:]
	fmt.Printf("%x\n", coeffient)

	str:= hex.EncodeToString(exponent)
	fmt.Printf("str=%s\n", str)
	//
	exp,_:= strconv.ParseInt(str, 16, 8)
	fmt.Printf("exp = %d\n", exp)

	result = append(bytes.Repeat([]byte{0x00}, 32-int(exp)), coeffient...)
	result = append(result, bytes.Repeat([]byte{0x00},32-len(result))...)
	return result
}

//交易拿到区块的默克尔根
func (b *Block) createMerkelTreeNode(transations []*Transation) {
	var tranHash [][]byte
	for _,tx:= range transations{
		tranHash = append(tranHash, tx.Hash())
	}
	mTree := NewMerkleTree(tranHash)
	b.Merkleroot = mTree.RootNode.Data
}

func (b*Block)String() {
	fmt.Printf("version :%s \n ", strconv.FormatInt(int64(b.Version), 10))
	fmt.Printf("prev.BlockHash :%x \n ", b.PrevBlockHash)
	fmt.Printf("Merkleroot :%x \n ", b.Merkleroot)
	fmt.Printf("Hash :%x \n ", b.Hash)
	fmt.Printf("time :%s \n ", strconv.FormatInt(int64(b.Time), 10))
	fmt.Printf("bits :%s \n ", strconv.FormatInt(int64(b.Bits), 10))
	fmt.Printf("Nonce :%s \n ", strconv.FormatInt(int64(b.Nonce), 10))
}

func NewBlock(transations []*Transation, preBlockHash []byte, height int32) *Block {
	block:= &Block{
		Version:2,
		PrevBlockHash:preBlockHash,
		Merkleroot:[]byte{},
		Hash:[]byte{},
		Time:int32(time.Now().Unix()),
		Bits:404454261,
		Nonce:0,
		Transations:transations,
		Height:height,
	}
	pow := NewProofofWork(block)
	nonce,hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash
	return block
}

func NewGensisBlock(transations []*Transation) *Block{
	block:= &Block{
		Version:2,
		PrevBlockHash:[]byte{},
		Merkleroot:[]byte{},
		Hash:[]byte{},
		Time:int32(time.Now().Unix()),
		Bits:404454261,
		Nonce:0,
		Transations:transations,
		Height:0,
	}
	pow := NewProofofWork(block)
	nonce,hash:=pow.Run()
	block.Nonce=nonce
	block.Hash=hash
	block.String()
	return block
}




