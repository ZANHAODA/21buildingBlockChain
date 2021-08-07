package main

func main() {
	bc := NewBlockchain("1NeBzmfLDxingHwNdzoA5y8cfY")
	cli := CLI{bc}
	cli.Run()
	//wallet := Newwallet()
	//fmt.Printf("私钥： %x",wallet.PrivateKey.D.Bytes())
	//fmt.Printf("公钥： %x",wallet.Publickey)
	//fmt.Printf("地址： %x",wallet.GetAddress())
	//
	//a,_:=hex.DecodeString("")
	//fmt.Printf("%d", ValidateAddress(a))
}
