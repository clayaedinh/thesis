package main

import (
	"fmt"
	"log"

	"github.com/clayaedinh/thesis/application/rsa/src"
)

func main() {
	//get public key
	pubkey, err := src.ReadUserPubkey("user0001")
	if err != nil {
		log.Panic(err)
	}
	//connection
	src.SetConnectionVariables("org1", "user0001", "localhost:7051")
	client := src.ChaincodeConnect()

	err = src.ChainStoreUserPubkey(client, "user0001", pubkey)
	if err != nil {
		log.Panic(err)
	}
	out, err := src.ChainRetrieveUserPubkey(client, "user0001")
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("out: %v\n", out)
}
