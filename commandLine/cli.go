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

type CommandLine struct {
	Blockchain *blockchain.BlockChain
}

func (cli *CommandLine) PrintHelp() {
	fmt.Println("Help:")
	fmt.Println("add -block BLOCK_DATA => add a block to the chain")
	fmt.Println("print => prints all block in the chain")
}

func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintHelp()
		runtime.Goexit()
	}
}

func (cli *CommandLine) AddBlock(data string) {
	cli.Blockchain.AddBlock(data)
	fmt.Println("Block added.")
}

func (cli *CommandLine) PrintChain() {
	iter := cli.Blockchain.Iterator()
	for {
		block := iter.Next()
		fmt.Printf("Pervious Hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("POW: %s\n", strconv.FormatBool(pow.Validate()))
		if len(block.PrevHash) == 0 {
			break
		}
	}
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
