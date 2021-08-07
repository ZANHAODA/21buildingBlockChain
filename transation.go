package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"math/big"
	"math/rand"
)
const subsidy = 100
type Transation struct {
	ID []byte	//交易的hash
	Vin []TXInput
	Vout []TXOutput
}
//把from 理解为  某人recv
type TXInput struct {
	TXid []byte   //前一笔交易的hash
	Voutindex int
	Signature []byte
	Pubkey []byte
}
//计算哈希值
type TXOutput struct{
	Value int
	PubkeyHash []byte  //要输出给谁  公钥的hash  再加数字签名
}

type TXOutputs struct {
	Outputs []TXOutput
}

func (outs TXOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc:= gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

func DeserializeOutputs(data []byte) TXOutputs {
	var outputs TXOutputs
	dec := gob.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&outputs)

	if err != nil {
		log.Panic(err)
	}
	return outputs
}

func (out *TXOutput) Lock(address []byte) {
	decodeAddress := Base58Decode(address)

	pubkeyhash := decodeAddress[1:len(decodeAddress) -4]
	out.PubkeyHash = pubkeyhash
}

func (tx Transation) Serialize() []byte {
	var encoded bytes.Buffer
	enc:= gob.NewEncoder(&encoded)

	err:= enc.Encode(tx)
	if err != nil{
		log.Panic(err)
	}
	return encoded.Bytes()
}

func (tx *Transation) Hash() []byte {
	txcopy:= *tx
	txcopy.ID = []byte{}
	hash:= sha256.Sum256(txcopy.Serialize())
	return hash[:]
}
//根据金额与地址新建一个输出
func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))
	return txo
}
//第一笔coinbase交易
func NewCoinbaseTx(to,data string) *Transation{
	txin := TXInput{[]byte{}, -1, nil,[]byte(data)}
	txout := NewTXOutput(subsidy, to)
	tx := Transation{nil, []TXInput{txin}, []TXOutput{*txout}}
	tx.ID = tx.Hash()
	return &tx
}

func (out *TXOutput) CanBeUnlockedWith(pubkeyhash []byte) bool {

	return bytes.Compare(out.PubkeyHash, pubkeyhash)==0
}

func (in *TXInput) canUnlockOutputWith(unlockdata []byte) bool {
	lockinghash := HashPubkey(in.Pubkey)
	return bytes.Compare(lockinghash, unlockdata) == 0
}

func (tx Transation) IsCoinBase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].TXid) == 0 && tx.Vin[0].Voutindex == -1
}

func (tx *Transation) Sign(privkey ecdsa.PrivateKey, prevTXs map[string]Transation) {
	if tx.IsCoinBase() {
		return
	}
	for _,vin := range tx.Vin{
		if prevTXs[hex.EncodeToString(vin.TXid)].ID == nil {
			log.Panic("Error")
		}
	}
	txcopy:= tx.TrimmedCopy()
	for inID,vin :=range txcopy.Vin {
		prevTX := prevTXs[hex.EncodeToString(vin.TXid)] //前一笔交易的结构体
		txcopy.Vin[inID].Signature = nil
		txcopy.Vin[inID].Pubkey = prevTX.Vout[vin.Voutindex].PubkeyHash //当前输入引用的前一笔交易的输出  的公钥哈希
		txcopy.ID = txcopy.Hash()
		r,s, err := ecdsa.Sign(rand.Reader, &privkey, txcopy.ID)  //交易的hash结果做了签名
		if err != nil {
			log.Panic(err)
		}
		signature:= append(r.Bytes(), s.Bytes()...)
		tx.Vin[inID].Signature = signature
	}
}

func (tx *Transation) TrimmedCopy() Transation {
	var inputs []TXInput
	var outputs []TXOutput

	for _,vin := range tx.Vin{
		inputs= append(inputs, TXInput{vin.TXid, vin.Voutindex,nil,nil})
	}
	for _,vout := range tx.Vout{
		outputs = append(outputs,TXOutput{vout.Value,vout.PubkeyHash})
	}
	txCopy:= Transation{tx.ID,inputs,outputs}
	return txCopy
}

func (tx *Transation) Verify(prevTxs map[string]Transation) bool {
	if tx.IsCoinBase() {
		return true
	}
	for _,vin := range tx.Vin {
		if prevTxs[hex.EncodeToString(vin.TXid)].ID == nil {
			log.Panic("ERRor")
		}
	}
	txcopy :=tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID,vin := range tx.Vin {
		prevTx := prevTxs[hex.EncodeToString(vin.TXid)]
		txcopy.Vin[inID].Signature = nil
		txcopy.Vin[inID].Pubkey = prevTx.Vout[vin.Voutindex].PubkeyHash
		txcopy.ID = txcopy.Hash()

		//r:=big.Int{}
		//s := big.Int{}
		//siglen := len(vin.Signature)
		//r.SetBytes(vin.Signature[:siglen/2])
		//s.SetBytes(vin.Signature[siglen/2:])

		x:=big.Int{}
		y := big.Int{}

		kenlen := len(vin.Pubkey)
		txcopy.Vin[inID].Pubkey = nil
		x.SetBytes(vin.Pubkey[:kenlen/2])   //公钥前一半为x  后一半为Y
		y.SetBytes(vin.Pubkey[kenlen/2:])

		rawPubkey := ecdsa.PublicKey{curve, &x, &y}

		if ecdsa.Verify(&rawPubkey,txcopy.ID,&x, &y) == false{
			return false
		}

		txcopy.Vin[inID].Pubkey = nil

	}
	return true
}


func NewUTXOTransation(from,to string,amount int, bc *Blockchain) *Transation {
	var inputs []TXInput
	var outputs []TXOutput

	wallets,err := NewWallets()
	if err!=nil{
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)
	acc,validoutputs := bc.FindSpendableOutputs(HashPubkey(wallet.Publickey), amount)
	if acc < amount{
		log.Panic("Error: Not enough funds")
	}
	for txid,outs:= range validoutputs{  //把找到的from的所有可消费的输出都 构建对应的输入
		txID, err:= hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		for _,out := range outs {
			input:= TXInput{txID, out, nil,wallet.Publickey}
			inputs = append(inputs,input)
		}
	}
	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc > amount{
		outputs = append(outputs, *NewTXOutput(acc - amount, from))   //把剩余的还给自己
	}
	tx:= Transation{nil,inputs,outputs}
	tx.ID = tx.Hash()
	bc.SignTransation(&tx, wallet.PrivateKey)
	return &tx
}

