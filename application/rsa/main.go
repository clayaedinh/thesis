package main

import "github.com/clayaedinh/thesis/application/rsa/application"

func main() {
	application.SetConnectionVariables("org1", "user1", "localhost:7051")

	//connection := chaincode.ChaincodeConnect()

}
