package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"

	"github.com/mohammad-hakemi22/blockchain/utility"
)

const walletFile = "./tmp/wallets.data"

type Wallets struct {
	Wallets map[string]*Wallet
}

func (ws *Wallets) SaveFile() {
	var content bytes.Buffer
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	utility.ErrorHandler("can't encode wallets", err)
	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	utility.ErrorHandler("can't write wallets file", err)
}
