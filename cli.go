package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
	bc *Blockchain
}
func (cli *CLI) addBlock() {
	cli.bc.MineBlock([]*Transation{})
}
func (cli *CLI) validateArgs() {
	if len(os.Args) < 1{
		fmt.Println("参数小于1")
		os.Exit(1)
	}

}
func (cli *CLI) printChain() {
	cli.bc.printBlockchain()
}

func (cli *CLI) send(from,to string, amount int) {
	tx:=NewUTXOTransation(from, to , amount, cli.bc)
	newBlock := cli.bc.MineBlock([]*Transation{tx})
	set := UTXOSet{cli.bc}
	set.update(newBlock)
	cli.getBalance("")
	fmt.Print("success! \n")
}
//终端传入地址，程序解码后  包含pubkeyhash
func (cli *CLI) getBalance(address string) {
	balance := 0
	decodeAddress := Base58Decode([]byte(address))
	pubkeyHash := decodeAddress[1:len(decodeAddress)-4]

	set := UTXOSet{cli.bc}
	UTXOs:= set.FindUTXObyPubkeyHash(pubkeyHash)
	for _,out:= range UTXOs{
		balance+=out.Value
	}
	fmt.Printf("\nbalance of '%s':%d \n",address, balance)
}

func (cli *CLI) createWallet() {
	wallets, _ := NewWallets()
	address:=wallets.CreateWallet()
	wallets.SaveToFile()
	fmt.Printf("your address：%s \n", address)
}
func (cli *CLI) printUsage() {
	fmt.Printf("Usage")
	fmt.Printf("addBlock-增加区块")

	fmt.Printf("printChain: 打印区块链")
}

func (cli *CLI) listAddress() {
	wallets,err:= NewWallets()
	if err!=nil{
		log.Panic(err)
	}
	addresses:=wallets.getAddress()
	for _,address :=range addresses{
		fmt.Printf(address)
	}
}
func (cli *CLI) Run() {
	cli.validateArgs()
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		fmt.Printf("NODE_ID is not set")
		os.Exit(1)
	}

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	
	startNodeCmd := flag.NewFlagSet("startnode", flag.ExitOnError)
	startNodeminner := startNodeCmd.String("minner", "", "minner address")

	getBalanceAddress := getBalanceCmd.String("address", "","the address to get balance of")  //如果有address这个参数， 就赋值给变量
	sendCmd:= flag.NewFlagSet("send",flag.ExitOnError)
	sendFrom:=sendCmd.String("from","", "source wallet address")
	sendTo := sendCmd.String("to","","Destination wallet address")

	sendAmount := sendCmd.Int("amount",0,"Amount to send")

	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressCmd := flag.NewFlagSet("listaddress", flag.ExitOnError)

	getBestHeightCmd := flag.NewFlagSet("getBestHeight", flag.ExitOnError)
	switch os.Args[1] {
	case "startnode":
		err := startNodeCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "getBestHeight":
		err := getBestHeightCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "listaddress":
		err := listAddressCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "addblock":
		err:= addBlockCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	case "printChain":
		err:= printChainCmd.Parse(os.Args[2:])
		if err!=nil{
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}
	if addBlockCmd.Parsed() {
		cli.addBlock()
	}
	if printChainCmd.Parsed() {
		cli.printChain()
	}
	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			os.Exit(1)
		}
		cli.getBalance(*getBalanceAddress)
	}
	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0{
			os.Exit(1)
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
	if createWalletCmd.Parsed() {
		cli.createWallet()
	}
	if listAddressCmd.Parsed() {
		cli.listAddress()
	}
	if getBestHeightCmd.Parsed() {
		cli.getBestHeight()
	}
	if startNodeCmd.Parsed() {
		nodeID := os.Getenv("NODE_ID")
		if nodeID == "" {
			startNodeCmd.Usage()
			os.Exit(1)
		}
		cli.startNode(nodeID, *startNodeminner)
	}
}

func (cli *CLI) getBestHeight() {
	cli.bc.GetBestHeight()

	fmt.Println(cli.bc.GetBestHeight())
}

func (cli *CLI) startNode(nodeID string, minnerAddress string) {
	fmt.Printf("Starting node%s \n", nodeID)
	if len(minnerAddress) > 0 {
		if ValidateAddress([]byte(minnerAddress)) {
			fmt.Printf("minner is on ", minnerAddress)
		} else {
			log.Panic("error minner Address")
		}
	}

	StartServer(nodeID, minnerAddress, cli.bc)
}


