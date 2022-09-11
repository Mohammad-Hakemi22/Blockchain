package commandline

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/mohammad-hakemi22/blockchain/blockchain"
	"github.com/mohammad-hakemi22/blockchain/utility"
)

type CommandLine struct {}

func (cli *CommandLine) PrintHelp() {
	fmt.Println("Help:")
	fmt.Println("getbalance -address ADDRESS => get the balance for address")
	fmt.Println("createblockchain -address ADDRESS => create a blockchain")
	fmt.Println("send -from FROM -to TO -amount AMOUNT => send amount of token to other address")
	fmt.Println("print => prints all block in the chain")
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

func (cli *CommandLine) Run() {
	cli.ValidateArgs()
	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		utility.ErrorHandler("can't parse args", err)
	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		utility.ErrorHandler("can't parse args", err)
	default:
		cli.PrintHelp()
		runtime.Goexit()
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.AddBlock(*addBlockData)
	}
	if printChainCmd.Parsed() {
		cli.PrintChain()
	}
}
