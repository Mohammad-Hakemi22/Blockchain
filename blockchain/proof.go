package blockchain

import (
	"bytes"
	"math/big"

	"github.com/mohammad-hakemi22/blockchain/utility"
)

const Difficulty = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	return &ProofOfWork{b, target}
}

func (pow *ProofOfWork) InitDate(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash, pow.Block.Data, utility.ToHex(int64(nonce)), utility.ToHex(int64(Difficulty))},
		[]byte{})
	return data
}
