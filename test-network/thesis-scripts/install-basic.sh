cd ..
./network.sh deployCC -ccn basic -ccp ../chaincode/basic -ccl go -ccep "OR('Org1MSP.peer','Org2MSP.peer')" -cccg ../chaincode/basic/collections_config.json