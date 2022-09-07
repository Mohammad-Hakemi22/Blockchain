package blockchain

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/mohammad-hakemi22/blockchain/utility"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

func InitBlockChain() *BlockChain {
	var lastHash []byte
	opt := badger.DefaultOptions(dbPath)
	opt.ValueDir = dbPath
	opt.Dir = dbPath

	db, err := badger.Open(opt)
	utility.ErrorHandler("something wrong in opening database", err)
	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("lh")); err == badger.ErrKeyNotFound { // check for existing blockchain
			fmt.Println("Not existing Blockchain!")
			genesis := GenesisBlock()
			fmt.Println("Genesis created!")
			err := txn.Set(genesis.Hash, genesis.Serialize())
			utility.ErrorHandler("can't set block in database", err)
			err = txn.Set([]byte("lh"), genesis.Hash)
			lastHash = genesis.Hash
			return err
		} else {
			item, err := txn.Get([]byte("lh"))
			utility.ErrorHandler("can't get last hash", err)
			err = item.Value(func(val []byte) error {
				lastHash = val
				return nil
			})
			return err
		}
	})
	utility.ErrorHandler("can't do R/W operation on database", err)
	blockchain := BlockChain{LastHash: lastHash, Database: db}
	return &blockchain
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utility.ErrorHandler("can't get last hash", err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})
	utility.ErrorHandler("can't get block from database", err)
	newBlock := CreateBlock(data, lastHash)
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		utility.ErrorHandler("can't set block in database", err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		chain.LastHash = newBlock.Hash
		return err
	})
	utility.ErrorHandler("can't do R/W operation on database", err)
}
