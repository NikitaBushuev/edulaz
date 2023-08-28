package network

import (
	"bytes"
	"crypto/ecdsa"
	"log"
	"net"

	elzbc "blockchain"
	elzut "utility"
)

const (
	BUFF_SIZE = 1024 * 1024
)

const (
	_MINING_DIFFICULTY = 2
)

const (
	LENGTH = iota + 1
	CHAIN
	TX
	PROVE
	BALANCE
	HASH
	LAST
	ADDR
)

func Read(conn net.Conn) []byte {
	var buf [BUFF_SIZE]byte
	n, err := conn.Read(buf[:])
	elzut.LogError(err)
	return buf[:n]
}

func Write(conn net.Conn, data []byte) {
	_, err := conn.Write(data)
	elzut.LogError(err)
}

func WriteAndRead(conn net.Conn, data []byte) []byte {
	Write(conn, data)
	return Read(conn)
}

func Request(address string, option byte, data []byte) []byte {
	conn, err := net.Dial(_NETWORK, address)
	if elzut.LogError(err) != nil {
		return nil
	}
	defer conn.Close()
	return WriteAndRead(conn, bytes.Join([][]byte{{option}, data}, nil))
}

func Listen(address string, handle func(conn net.Conn)) {
	ln, err := net.Listen(_NETWORK, address)
	if elzut.LogError(err) != nil {
		return
	}
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if elzut.LogError(err) != nil {
			break
		}
		go handle(conn)
	}
}

func Length(address string) int {
	length := int32(0)
	elzut.Deserialize(&length, Request(address, LENGTH, nil))
	return int(length)
}

func Chain(address string, chain *elzbc.Chain) {
	elzut.Deserialize(chain, Request(address, CHAIN, nil))
}

func Tx(address string, tx *elzbc.Transaction) bool {
	return Request(address, TX, elzut.Serialize(tx)) != nil
}

func Prove(address string, block *elzbc.Block) bool {
	return Request(address, PROVE, elzut.Serialize(block)) != nil
}

func Balance(address string, addr string) int64 {
	balance := int64(0)
	elzut.Deserialize(&balance, Request(address, BALANCE, elzut.Serialize(addr)))
	return balance
}

func Hash(address string, index int) string {
	var hash string
	elzut.Deserialize(&hash, Request(address, HASH, elzut.Serialize(int32(index))))
	return hash
}

func Last(address string) int {
	index := int32(0)
	elzut.Deserialize(&index, Request(address, LAST, nil))
	return int(index)
}

func Addr(address string) string {
	var addr string
	elzut.Deserialize(&addr, Request(address, ADDR, nil))
	return addr
}

func Sync(chain *elzbc.Chain, addresses []string) {
	Chain(Choose(addresses), chain)
}

func Choose(addresses []string) string {
	var best string
	if len(addresses) != 0 {
		best = addresses[0]
	}
	lmax := 0
	for _, address := range addresses {
		length := Length(address)
		if length > lmax {
			best = address
		}
	}
	return best
}

func Serve(address string, addresses []string,
	key *ecdsa.PrivateKey, myblock *elzbc.Block, chain *elzbc.Chain) {

	mining := false
	myaddr := elzbc.Address(&key.PublicKey)

	mapping := map[byte]func(data []byte) []byte{
		LENGTH: func(data []byte) []byte {
			return elzut.Serialize(int32(chain.Len()))
		},

		CHAIN: func(data []byte) []byte {
			return elzut.Serialize(chain)
		},

		TX: func(data []byte) []byte {
			var tx elzbc.Transaction
			elzut.Deserialize(&tx, data)

			if tx.Invalid() {
				elzut.LogError(elzbc.ErrInvalidTransaction)
				return nil
			}

			myblock.Append(tx)

			if myblock.Full() {
				go func() {
					if myblock.Mine(myaddr, _MINING_DIFFICULTY, &mining) {
						myblock.Sign(key)
						for _, address := range addresses {
							log.Println("accept", address, Prove(address, myblock))
						}
					}
					*myblock = elzbc.NewBlock(chain.Previous())
				}()
			}

			return []byte{1}
		},

		PROVE: func(data []byte) []byte {
			var block elzbc.Block
			elzut.Deserialize(&block, data)

			if block.Invalid() {
				elzut.LogError(elzbc.ErrInvalidBlock)
				return nil
			}

			mining = false

			chain.Append(block)

			*myblock = elzbc.NewBlock(chain.Previous())

			return []byte{1}
		},

		BALANCE: func(data []byte) []byte {
			var addr string
			elzut.Deserialize(&addr, data)
			return elzut.Serialize(int64(chain.Balance(addr) + myblock.Balance(addr)))
		},

		HASH: func(data []byte) []byte {
			var index int32
			elzut.Deserialize(&index, data)
			return elzut.Serialize(chain.Hash(int(index)))
		},

		LAST: func(data []byte) []byte {
			return elzut.Serialize(int32(chain.Last()))
		},

		ADDR: func(data []byte) []byte {
			return elzut.Serialize(myaddr)
		},
	}

	Listen(address, func(conn net.Conn) {
		request := Read(conn)
		if handle, ok := mapping[request[0]]; ok {
			Write(conn, handle(request[1:]))
		}
		conn.Close()
	})
}

var (
	_NETWORK = "tcp"
)
