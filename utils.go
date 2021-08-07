package main

import (
	"bytes"
	_ "crypto/sha256"
	"encoding/binary"
	_ "encoding/hex"
	_ "fmt"
	"log"
	_ "math/big"
)

func min(a int, b int ) int{
	if (a>b) {
		return b
	}
	return a
}

func IntToHex(num int32) []byte{
	buff := new(bytes.Buffer)

	err := binary.Write(buff, binary.LittleEndian, num)
	if err != nil{
		log.Panic(err)
	}
	return buff.Bytes()
}
//将类型转换为了字节数组，大端
func IntToHex2(num int32) []byte{
	buff := new(bytes.Buffer)

	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil{
		log.Panic(err)
	}
	return buff.Bytes()
}
//字节反转
func ReverseBytes(data []byte) {
	for i,j := 0,len(data)-1; i < j; i,j = i+1,j-1{
		data[i], data[j] = data[j],data[i]
	}
}
