package main

import (
	"os"

	cmd "github.com/mohammad-hakemi22/blockchain/commandline"
)

func main() {
	defer os.Exit(0)
	cli := cmd.CommandLine{}
	cli.Run()
}
