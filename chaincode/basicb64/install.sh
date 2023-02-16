cd ../../test-network
./network.sh deployCC -ccn basicb64 -ccp ../chaincode/basicb64 -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')" -cccg ../chaincode/basicb64/collections_config.json