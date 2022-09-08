package main

import (
	"os"

	"github.com/mohammad-hakemi22/blockchain/blockchain"
	cmd "github.com/mohammad-hakemi22/blockchain/commandline"
)

func main() {
	defer os.Exit(0)
	chain := blockchain.InitBlockChain()
	defer chain.Database.Close()

	cli := cmd.CommandLine{Blockchain: chain}
	cli.Run()
}
