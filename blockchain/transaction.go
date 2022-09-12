package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/mohammad-hakemi22/blockchain/utility"
)

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte
	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	utility.ErrorHandler("can't encode the transaction", err)
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Input) == 1 && len(tx.Input[0].ID) == 0 && tx.Input[0].Out == -1
}

func CoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}

	txin := TxInput{ID: []byte{}, Out: -1, Sig: data}
	txout := TxOutput{Value: 100, PubKey: to}

	tx := Transaction{nil, []TxInput{txin}, []TxOutput{txout}}
	tx.SetID()
	return &tx
}

func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput
	acc, validOutputs := chain.FindSpendableOutputs(from, amount)
	if acc < amount {
		utility.ErrorHandler("not enough funds", errors.New("not enough funds"))
	}
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		utility.ErrorHandler("can't decode tx id", err)
		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}
	outputs = append(outputs, TxOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, to})
	}
	tx := Transaction{nil, inputs, outputs}
	tx.SetID()
	return &tx
}