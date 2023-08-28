package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"log"
	"math"
	"math/big"

	mrand "math/rand"

	elzut "utility"
)

const (
	CREATOR_REWARD = 1024

	BLOCK_MAXLEN = 2
	CHAIN_MAXLEN = 1024
)

var (
	CURVE = elliptic.P256()
)

var (
	ErrInvalidTransaction = errors.New("invalid transaction")
	ErrInvalidBlock       = errors.New("invalid block")
)

type Creature struct {
	Verifier  string
	Signature string
	Hash      string
}

type TransactionData struct {
	Id       uint64
	Previous string
	Sender   string
	Receiver string
	Amount   int64
}

type Transaction struct {
	Creature
	TransactionData
}

type BlockData struct {
	Nonce        uint32
	Previous     string
	Miner        string
	Transactions []Transaction
}

type Block struct {
	Creature
	BlockData
}

type Chain struct {
	Creator string

	Blocks []Block
}

func Address(pub *ecdsa.PublicKey) string {
	data := elliptic.MarshalCompressed(pub.Curve, pub.X, pub.Y)
	return elzut.EncodeToString([]byte(elzut.HashSum(data)))
}

func (block *Block) Reward() int64 {
	reward := int64(1)
	for _, h := range block.Hash {
		if h != 'A' {
			break
		}
		reward *= 2
	}
	return reward
}

func (block *Block) Balance(addr string) int64 {
	balance := int64(0)
	for _, tx := range block.Transactions {
		switch addr {
		case tx.Sender:
			balance -= tx.Amount
		case tx.Receiver:
			balance += tx.Amount
		}
	}
	switch addr {
	case block.Miner:
		balance += block.Reward()
	}
	return balance
}

func (chain *Chain) Balance(addr string) int64 {
	balance := int64(0)
	switch addr {
	case chain.Creator:
		balance += CREATOR_REWARD
	}
	for _, block := range chain.Blocks {
		balance += block.Balance(addr)
	}
	return balance
}

func (chain *Chain) Hash(index int) string {
	block := chain.Block(index)
	if block != nil {
		return block.Hash
	}
	return elzut.HASH_NULL
}

func (block *Block) Len() int {
	return len(block.Transactions)
}

func (chain *Chain) Len() int {
	return len(chain.Blocks)
}

func (block *Block) Append(tx Transaction) {
	if block.Len() >= BLOCK_MAXLEN {
		panic(errors.New("block overflow"))
	}
	block.Transactions = append(block.Transactions, tx)
}

func (chain *Chain) Append(block Block) {
	if chain.Len() >= CHAIN_MAXLEN {
		panic(errors.New("chain overflow"))
	}
	chain.Blocks = append(chain.Blocks, block)
}

func (creature *Creature) Sign(priv *ecdsa.PrivateKey) {
	creature.Verifier = elzut.EncodeToString(CompressPub(&priv.PublicKey))
	creature.Signature = Sign(priv, creature.Hash)
}

func (creature *Creature) Verify() bool {
	return Verify(DecompressPub(elzut.DecodeString(creature.Verifier)),
		creature.Hash, creature.Signature)
}

func CompressPub(pub *ecdsa.PublicKey) []byte {
	return elliptic.MarshalCompressed(CURVE, pub.X, pub.Y)
}

func DecompressPub(data []byte) *ecdsa.PublicKey {
	X, Y := elliptic.UnmarshalCompressed(CURVE, data)
	return &ecdsa.PublicKey{Curve: CURVE, X: X, Y: Y}
}

func CompressKey(priv *ecdsa.PrivateKey) []byte {
	return priv.D.Bytes()
}

func DecompressKey(data []byte) *ecdsa.PrivateKey {
	D := new(big.Int).SetBytes(data)
	X, Y := CURVE.ScalarBaseMult(D.Bytes())
	return &ecdsa.PrivateKey{PublicKey: ecdsa.PublicKey{Curve: CURVE, X: X, Y: Y}, D: D}
}

func (block *Block) Full() bool {
	return block.Len() >= BLOCK_MAXLEN
}

func (block *Block) Mine(miner string, dfc int, mining *bool) bool {
	block.Miner = miner
	block.Hash = elzut.HashSum(block.BlockData)

	log.Println("start mining address", block.Miner, "...")

	*mining = true

	for *mining && block.Nonce < math.MaxUint32 && block.Hash[:dfc] != elzut.HASH_ZERO[:dfc] {
		block.Nonce++
		block.Hash = elzut.HashSum(block.BlockData)
	}
	log.Println(block.Hash)

	return *mining
}

func (creature *Creature) Invalid() bool {
	return !creature.Verify()
}

func (chain *Chain) Empty() bool {
	return chain.Len() < 1
}

func (chain *Chain) Last() int {
	return chain.Len() - 1
}

func (chain *Chain) Block(index int) *Block {
	if index >= 0 && index < chain.Len() {
		return &chain.Blocks[index]
	}
	return nil
}

func (chain *Chain) Previous() string {
	return chain.Hash(chain.Last())
}

func Verify(pub *ecdsa.PublicKey, hash string, sig string) bool {
	return ecdsa.VerifyASN1(pub, elzut.DecodeString(hash), elzut.DecodeString(sig))
}

func Sign(priv *ecdsa.PrivateKey, hash string) string {
	sig, err := ecdsa.SignASN1(rand.Reader, priv, elzut.DecodeString(hash))
	elzut.LogError(err)
	return elzut.EncodeToString(sig)
}

func NewTransaction(previous string, sender string, receiver string, amount int64) Transaction {
	tx := Transaction{
		TransactionData: TransactionData{
			Id:       mrand.Uint64(),
			Previous: previous,
			Sender:   sender,
			Receiver: receiver,
			Amount:   amount,
		},
	}
	tx.Hash = elzut.HashSum(tx)
	return tx
}

func GenerateKey() *ecdsa.PrivateKey {
	priv, err := ecdsa.GenerateKey(CURVE, rand.Reader)
	elzut.LogError(err)
	return priv
}

func NewBlock(previous string) Block {
	return Block{
		BlockData: BlockData{
			Previous: previous,
		},
	}
}

func NewChain(pub *ecdsa.PublicKey) *Chain {
	return &Chain{
		Creator: Address(pub),
	}
}
