package cli

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"os"
	"path"

	elzbc "blockchain"
	elznt "network"
	elzut "utility"
)

func Run() {
	mypath := "."

	if len(os.Args) > 1 {
		mypath = os.Args[1]
	}

	key_path := path.Join(mypath, "private_key.json")
	addrs_path := path.Join(mypath, "addresses.json")

	addresses := []string{}
	elzut.Load(addrs_path, &addresses)
	elzut.Save(addrs_path, &addresses)

	if len(addresses) < 1 {
		log.Println("addresses is empty, configure addresses.json")
		return
	}

	key := &ecdsa.PrivateKey{PublicKey: ecdsa.PublicKey{Curve: elzbc.CURVE}}
	if elzut.Load(key_path, key) != nil {
		key = elzbc.GenerateKey()
	}
	elzut.Save(key_path, key)

	myaddr := elzbc.Address(&key.PublicKey)
	log.Println("address:", myaddr)

	for {
		fmt.Print("> ")

		var cmd string
		fmt.Scan(&cmd)

		switch cmd {
		case "exit":
			return

		case "myaddress":
			log.Println(myaddr)

		case "mybalance":
			best := elznt.Choose(addresses)
			log.Println("choose", best)

			log.Println(elznt.Balance(best, myaddr))

		case "balance":
			var addr string
			fmt.Scan(&addr)

			best := elznt.Choose(addresses)
			log.Println("choose", best)

			log.Println(elznt.Balance(best, addr))

		case "tx":
			var addr string
			fmt.Scan(&addr)

			var amount int64
			fmt.Scan(&amount)

			best := elznt.Choose(addresses)
			log.Println("choose", best)

			prev := elznt.Hash(best, elznt.Last(best))

			tx := elzbc.NewTransaction(prev, myaddr, addr, amount)
			tx.Sign(key)

			for _, address := range addresses {
				log.Println("accept", address, elznt.Tx(address, &tx))
			}
		}
	}
}
