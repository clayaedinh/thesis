package main

import "github.com/clayaedinh/thesis/application/rsa/chaincode"

func main() {
	chaincode.SetConnectionVariables("org1", "user1", "localhost:7051")

	//connection := chaincode.ChaincodeConnect()

}
