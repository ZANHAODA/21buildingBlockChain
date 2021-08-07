package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
)

const nodeversion = 0x00

var nodeAddress string
var blockInTransit [][]byte
const commonLength = 12
type Version struct {
	Version int
	BestHeight int32
	AddrFrom string
}

func (ver *Version) String(){
	fmt.Printf("Version: %d \n", ver.Version)
	fmt.Printf("Version: %d \n", ver.BestHeight)
	fmt.Printf("Version: %d \n", ver.AddrFrom)
}
var knownNodes = []string{"localhost:3000"}
//区块少的a向b 发送版本号  b发现自己区块多，也向a发送版本号; a区块少，执行sendGetBlock ，发送getblocks命令+信息
//b调用处理GetBlock函数, 处理所有区块信息，调用sendInv;  a解析inv, 调用sendGetData
func StartServer(nodeID,minerAddress string, bc *Blockchain) {

	nodeAddress = fmt.Sprintf("localhost:%s",nodeID)
	ln,err := net.Listen("tcp",nodeAddress)
	if err!=nil{
		log.Panic(err)
	}
	defer ln.Close()
//	bc := NewBlockchain("1NeBzmfLd")
	if nodeAddress != knownNodes[0] {
		sendVersion(knownNodes[0],bc)
	}

	for {
		conn,err:=ln.Accept()
		if err!=nil{
			log.Panic(err)
		}
		go handlerConnection(conn, bc)  //协程
	}

}

func handlerConnection(conn net.Conn, bc *Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err!=nil{
		log.Panic(err)
	}
	//获取命令
	command := bytesToCommand(request[:12])
	switch command {
	case "version":
		handleVersion(request, bc)
	case "getblocks":
		handleGetBlock(request, bc)
	case "inv":
		handleInv(request,bc)
	case "getdata":
		handleGetData(request, bc)
	case "block":
		handleBlock(request, bc)
	}

}

func handleBlock(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload blocksend

	buff.Write(request[commonLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err!= nil {
		log.Panic(err)
	}
	blockdata := payload.Block
	block := DeserializeBlock(blockdata)
	bc.AddBlock(block)
	fmt.Printf("Receive a new block")
	if len(blockInTransit) > 0 {
		blockHash := blockInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)
		blockInTransit = blockInTransit[1:]
	}else {
		set := UTXOSet{}
		set.Reindex()
	}
}

func handleGetData(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload getdata
	buff.Write(request[commonLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err!= nil {
		log.Panic(err)
	}
	if payload.Type=="block"{
		block,err := bc.GetBlock([]byte(payload.ID))
		if err!= nil {
			log.Panic(err)
		}
		sendBlock(payload.AddrFrom, &block)
	}
}

type blocksend struct {
	AddrFrom string
	Block []byte
}

func sendBlock(addr string, block *Block) {
	data := blocksend{nodeAddress, block.Serialize()}
	payload := gobEncode(data)
	request:= append(commandToBytes("block"), payload...)

	sendData(addr, request)
}

func handleInv(request []byte, blockchain *Blockchain) {  //处理Inv节点是缺少的节点
	var buff bytes.Buffer

	var payload inv
	buff.Write(request[commonLength:])
	dec:= gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err!= nil {
		log.Panic(err)
	}

	fmt.Printf("Receive inventory %d %s", len(payload.Items), payload.Type)

	if payload.Type=="block" {
		blockInTransit = payload.Items
		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _,b:= range blockInTransit{
			if bytes.Compare(b, blockHash)!=0 {
				newInTransit = append(newInTransit,b)
			}
		}
		blockInTransit = newInTransit
	}
}

type getdata struct {
	AddrFrom string
	Type string
	ID []byte
}

func sendGetData(addr string, kind string, id []byte) {
	payload := gobEncode(getdata{addr, kind, id})
	request:= append(commandToBytes("getdata"), payload...)
	sendData(addr, request)
}

func handleGetBlock(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload getblocks

	buff.Write(request[commonLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err !=nil{
		log.Panic(err)
	}
	block := bc.Getblockhash()
	sendInv(payload.Addrfrom, "block",block)
}

type inv struct {
	AddrFrom string
	Type string
	Items [][]byte
}

func sendInv(addr string, kind string, items [][]byte) {
	inventory := inv{nodeAddress, kind, items}
	payload := gobEncode(inventory)
	request := append(commandToBytes("inv"), payload...)
	sendData(addr,request)
}

func handleVersion(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload Version
	buff.Write(request[commonLength:])

	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err!=nil{
		log.Panic(err)
	}
	payload.String()
	myBestHeight := bc.GetBestHeight()

	foreignerBestHeight := payload.BestHeight
	if myBestHeight < foreignerBestHeight {
		sendGetBlock(payload.AddrFrom)
	} else {
		sendVersion(payload.AddrFrom, bc)
	}

	if !nodeIsKnow(payload.AddrFrom){
		knownNodes = append(knownNodes, payload.AddrFrom)
	}
}

type getblocks struct {
	Addrfrom string
}

func sendGetBlock(address string) {
	payload := gobEncode(getblocks{nodeAddress})

	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func nodeIsKnow(addr string) bool {
	for _,node := range knownNodes{
		if node == addr {
			return true
		}
	}
	return false
}

func sendVersion(addr string, bc *Blockchain) {
	bestHeight := bc.GetBestHeight()

	payload := gobEncode(Version{nodeversion,bestHeight,nodeAddress})
	request := append(commandToBytes("version"), payload...)
	sendData(addr,request)
}

func sendData(addr string, data []byte) {
	con,err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("%s is not available ", addr)

		var updateNodes []string
		for _,node := range knownNodes {
			if node != addr{
				updateNodes = append(updateNodes, node)
			}
		}
		knownNodes = updateNodes
	}
	defer con.Close()
	_,err = io.Copy(con, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func commandToBytes(command string) []byte{
	var bates [commonLength]byte

	for i,c := range command{
		bates[i] = byte(c)
	}
	return bates[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte
	for _,b:= range bytes{
		if b!=0x00{
			command = append(command, b)
		}
	}
	return fmt.Sprintf("%s", command)
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)

	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
