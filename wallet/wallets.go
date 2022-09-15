package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"os"

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

func (ws *Wallets) LoadFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}
	var wallets Wallets
	fileContent, err := ioutil.ReadFile(walletFile)
	utility.ErrorHandler("can't load wallets file", err)
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	gob.Register(elliptic.P256())
	err = decoder.Decode(&wallets)
	utility.ErrorHandler("can't decode wallets content", err)
	ws.Wallets = wallets.Wallets
	return nil
}