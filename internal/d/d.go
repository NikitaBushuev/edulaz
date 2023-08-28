package d

import (
	"crypto/ecdsa"
	"log"
	"os"
	"path"
	"time"

	elzbc "blockchain"
	elznt "network"
	elzut "utility"
)

func Run() {
	log.Println("usage:", os.Args[0], "<path>", "<addr>")

	mypath := "."
	myaddress := ":9090"

	if len(os.Args) > 1 {
		mypath = os.Args[1]
	}

	if len(os.Args) > 2 {
		myaddress = os.Args[2]
	}

	chain_path := path.Join(mypath, "blockchain.json")
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

	mychain := &elzbc.Chain{}
	if elzut.Load(chain_path, mychain) != nil {
		mychain = elzbc.NewChain(&key.PublicKey)
	}
	elzut.Save(chain_path, mychain)

	log.Println("address:", elzbc.Address(&key.PublicKey))

	go func() {
		for {
			elznt.Sync(mychain, addresses)
			elzut.Save(chain_path, mychain)
			time.Sleep(time.Minute)
		}
	}()

	myblock := &elzbc.Block{}

	elznt.Serve(myaddress, addresses, key, myblock, mychain)
}
