package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
)
const version = byte(0x00)
type Wallet	struct {

	PrivateKey ecdsa.PrivateKey
	Publickey []byte
}

func Newwallet() *Wallet{
	private,public := newKeyPair()
	wallet:= Wallet{private,public}
	return &wallet
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	//生成椭圆曲线	specp256r1曲线	比特币当中的曲线是secp256k1
	curve:= elliptic.P256()
	private,err:= ecdsa.GenerateKey(curve,rand.Reader)
	if err != nil {
		log.Println("error")
	}
	pubkey := append(private.PublicKey.X.Bytes(),private.PublicKey.Y.Bytes()...)
	return *private,pubkey
}

func (w Wallet) GetAddress() []byte {
	pubkeyHash := HashPubkey(w.Publickey)
	versionPayload := append([]byte{version},pubkeyHash...)
	check :=checksum(versionPayload)
	fullPayload := append(versionPayload,check...)
	address:= Base58Encode(fullPayload)
	return address
}

func HashPubkey(pubkey []byte) []byte{
	pubkeyHash256:= sha256.Sum256(pubkey)
	PIPEMD160Hasher := ripemd160.New()

	_,err:= PIPEMD160Hasher.Write(pubkeyHash256[:])
	if err!=nil{
		fmt.Println("error")
	}
	publicRIPEMD160 :=PIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160
}

func checksum(paload []byte) []byte {
	firstSHA := sha256.Sum256(paload)
	secondSHA := sha256.Sum256(firstSHA[:])
	checksum := secondSHA[:4]
	return checksum
}

func ValidateAddress(address []byte) bool {
	pubkeyHash := Base58Decode(address)
	actualCheckSum := pubkeyHash[len(pubkeyHash)-4:]
	publickeyHash := pubkeyHash[1:len(pubkeyHash)-4]
	targetChecksum := checksum(append([]byte{0x00},publickeyHash...))
	return bytes.Compare(actualCheckSum, targetChecksum) ==0
}