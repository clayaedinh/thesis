cd ../../test-network
./network.sh deployCC -ccn rsa -ccp ../chaincode/rsa -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')" -cccg ../chaincode/rsa/collections_config.json