package commandline

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/mohammad-hakemi22/blockchain/blockchain"
	"github.com/mohammad-hakemi22/blockchain/utility"
	"github.com/mohammad-hakemi22/blockchain/wallet"
)

type CommandLine struct{}

func (cli *CommandLine) PrintHelp() {
	fmt.Println("Help:")
	fmt.Println("getbalance -address ADDRESS => get the balance for address")
	fmt.Println("createblockchain -address ADDRESS => create a blockchain")
	fmt.Println("send -from FROM -to TO -amount AMOUNT => send amount of token to other address")
	fmt.Println("print => prints all block in the chain")
	fmt.Println("createwallet => Creates a new Wallet")
	fmt.Println("listaddresses => Lists the addresses in our wallet file")
}

func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintHelp()
		runtime.Goexit()
	}
}

func (cli *CommandLine) PrintChain() {
	chain := blockchain.ContinueBlockchain("")
	defer chain.Database.Close()
	iter := chain.Iterator()
	for {
		block := iter.Next()
		fmt.Printf("Pervious Hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockchain(address)
	defer chain.Database.Close()
	balance := 0
	uTxs := chain.FindUTx(address)
	for _, out := range uTxs {
		balance += out.Value
	}
	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockchain(from)
	defer chain.Database.Close()
	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Printf("Successfully transfer %d token, from %s to %s.", amount, from, to)
}

func (cli *CommandLine) createBlockchain(address string) {
	chain := blockchain.InitBlockChain(address)
	chain.Database.Close()
	fmt.Println("finished!")
}

func (cli *CommandLine) listAddresses() {
	wallets, _ := wallet.CreateWallets()
	addresses := wallets.GetAllAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}

func (cli *CommandLine) createWallet() {
	wallets, _ := wallet.CreateWallets()
	address := wallets.AddWallet()
	wallets.SaveFile()

	fmt.Printf("New address is: %s\n", address)
}

func (cli *CommandLine) Run() {
	cli.ValidateArgs()
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)
	listAddressesCmd := flag.NewFlagSet("listaddresses", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "Balance address")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "Create blockchain address")
	sendFrom := sendCmd.String("from", "", "Source address")
	sendTo := sendCmd.String("to", "", "Destination address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		utility.ErrorHandler("can't parse args", err)
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		utility.ErrorHandler("can't parse args", err)
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		utility.ErrorHandler("can't parse args", err)
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		utility.ErrorHandler("can't parse args", err)
	case "listaddresses":
		err := listAddressesCmd.Parse(os.Args[2:])
		utility.ErrorHandler("can't parse args", err)
	case "createwallet":
		err := createWalletCmd.Parse(os.Args[2:])
		utility.ErrorHandler("can't parse args", err)
	default:
		cli.PrintHelp()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}
	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockchain(*createBlockchainAddress)
	}
	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
	if printChainCmd.Parsed() {
		cli.PrintChain()
	}
	if createWalletCmd.Parsed() {
		cli.createWallet()
	}
	if listAddressesCmd.Parsed() {
		cli.listAddresses()
	}
}
