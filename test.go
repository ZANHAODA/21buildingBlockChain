package main

import "fmt"

func TestCreateMerkelTreeRoot() {
	block:= &Block{
		Version:2,
		PrevBlockHash:[]byte{},
		Merkleroot:[]byte{},
		Hash:[]byte{},
		Time:1418755788,
		Bits:404454261,
		Nonce:0,
		Transations:[]*Transation{},
		Height:0,
	}

	txin := TXInput{[]byte{}, -1, nil,nil}
	txout := NewTXOutput(subsidy, "first")
	tx := Transation{nil, []TXInput{txin}, []TXOutput{*txout}}

	txin2 := TXInput{[]byte{}, -1, nil,nil}
	txout2 := NewTXOutput(subsidy, "second")
	tx2 := Transation{nil, []TXInput{txin2}, []TXOutput{*txout2}}

	var Transations []*Transation
	Transations = append(Transations, &tx, &tx2)
	//return ransations
	block.createMerkelTreeNode(Transations)
	fmt.Printf("%x\n", block.Merkleroot)
}

func TestNewSerialize() {
	block:= &Block{
		Version:2,
		PrevBlockHash:[]byte{},
		Merkleroot:[]byte{},
		Hash:[]byte{},
		Time:1418755788,
		Bits:404454261,
		Nonce:0,
		Transations:[]*Transation{},
		Height:0,
	}
	deBlock:=DeserializeBlock(block.Serialize())
	deBlock.String()
}
func TestPow() {
	block:= &Block{
		Version:2,
		PrevBlockHash:[]byte{},
		Merkleroot:[]byte{},
		Hash:[]byte{},
		Time:1418755788,
		Bits:404454261,
		Nonce:0,
		Transations:[]*Transation{},
		Height:0,
	}
	pow :=NewProofofWork(block)
	nonce,_ := pow.Run()
	block.Nonce = nonce
	fmt.Println(pow.Validate())
}

func TestBoltDB() {
	blockchain := NewBlockchain("1NeBzmfLDxingHwNdzoA5y8cfY")
	blockchain.MineBlock([]*Transation{})
	blockchain.MineBlock([]*Transation{})
	blockchain.MineBlock([]*Transation{})
	blockchain.printBlockchain()
}
func main() {
	TestBoltDB()
}
