package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/mohammad-hakemi22/blockchain/utility"
)

const (
	dbPath      = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
	genesisData = "First Transaction from genesis block"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockchainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func InitBlockChain(address string) *BlockChain {
	if DBexists() {
		fmt.Println("Blockchain already exists.")
		runtime.Goexit()
	}
	var lastHash []byte
	opt := badger.DefaultOptions(dbPath)
	opt.ValueDir = dbPath
	opt.Dir = dbPath
	db, err := badger.Open(opt)
	utility.ErrorHandler("something wrong in opening database", err)

	err = db.Update(func(txn *badger.Txn) error {
		cbtx := CoinbaseTx(address, genesisData)
		genesis := GenesisBlock(cbtx)
		fmt.Println("Genesis created!")
		err := txn.Set(genesis.Hash, genesis.Serialize())
		utility.ErrorHandler("can't set block in database", err)
		err = txn.Set([]byte("lh"), genesis.Hash)
		lastHash = genesis.Hash
		return err
	})
	utility.ErrorHandler("can't do R/W operation on database", err)
	blockchain := BlockChain{LastHash: lastHash, Database: db}
	return &blockchain
}

func ContinueBlockchain(address string) *BlockChain {
	if !DBexists() {
		fmt.Println("No existing Blockchain, created one")
		runtime.Goexit()
	}
	var lastHash []byte
	opt := badger.DefaultOptions(dbPath)
	opt.ValueDir = dbPath
	opt.Dir = dbPath
	db, err := badger.Open(opt)
	utility.ErrorHandler("something wrong in opening database", err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		utility.ErrorHandler("can't get block from database", err)
		item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		return err
	})
	utility.ErrorHandler("can't do R/W operation on database", err)
	blockchain := BlockChain{LastHash: lastHash, Database: db}
	return &blockchain
}

func (chain *BlockChain) AddBlock(transaction []*Transaction) {
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
	newBlock := CreateBlock(transaction, lastHash)
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		utility.ErrorHandler("can't set block in database", err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		chain.LastHash = newBlock.Hash
		return err
	})
	utility.ErrorHandler("can't do R/W operation on database", err)
}

func (chain *BlockChain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{chain.LastHash, chain.Database}
}

func (iter *BlockchainIterator) Next() *Block {
	var block *Block
	var encodedBlock []byte
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		utility.ErrorHandler("can't get block from database", err)
		err = item.Value(func(val []byte) error {
			encodedBlock = val
			return nil
		})
		block = Deserialize(encodedBlock)
		return err
	})
	utility.ErrorHandler("can't get blocks from database", err)
	iter.CurrentHash = block.PrevHash
	return block
}

func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxs []Transaction
	spentTxs := make(map[string][]int)
	iter := chain.Iterator()

	for {
		block := iter.Next()
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)
		Output:
			for outIdx, out := range tx.Output {
				if spentTxs[txID] != nil {
					for _, spentOut := range spentTxs[txID] {
						if spentOut == outIdx {
							continue Output
						}
					}
				}
				if out.CanBeUnlock(address) {
					unspentTxs = append(unspentTxs, *tx)
				}
			}
			if !tx.IsCoinbase() {
				for _, in := range tx.Input {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spentTxs[inTxID] = append(spentTxs[inTxID], in.Out)
					}
				}
			}
		}
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return unspentTxs
}

func (chain *BlockChain) FindUTx(address string) []TxOutput {
	var UTxs []TxOutput
	unspentTransactions := chain.FindUnspentTransactions(address)
	for _, tx := range unspentTransactions {
		for _, out := range tx.Output {
			if out.CanBeUnlock(address) {
				UTxs = append(UTxs, out)
			}
		}
	}
	return UTxs
}

func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)
		for outIdx, out := range tx.Output {
			if out.CanBeUnlock(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)
				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOuts
}
