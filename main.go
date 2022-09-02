package main

import (
	"fmt"
	"github.com/mohammad-hakemi22/blockchain/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()

	chain.AddBlock("First Block!")
	chain.AddBlock("Second Block!")
	chain.AddBlock("Third Block!")

	for _, block := range chain.Blocks {
		fmt.Printf("Pervious Hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
	}
}