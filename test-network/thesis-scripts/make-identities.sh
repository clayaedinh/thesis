#This script makes user identities for Doctor and Pharmacist, assuming the network is already set up.
cd ..
#CREATE ORG1 USER0001: DOCTOR ROLE
export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/org1.example.com/

fabric-ca-client register --caname ca-org1 --id.name user0001 --id.secret userpw0001 --id.type client --id.attrs 'role=DOCTOR:ecert' --tls.certfiles "${PWD}/organizations/fabric-ca/org1/tls-cert.pem"

fabric-ca-client enroll -u https://user0001:userpw0001@localhost:7054 --caname ca-org1 -M "${PWD}/organizations/peerOrganizations/org1.example.com/users/user0001@org1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org1/tls-cert.pem"

cp ${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org1.example.com/users/user0001@org1.example.com/msp/config.yaml

#CREATE ORG2 USER0002: PHARMA ROLE
export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=$PWD/../config/
export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/org2.example.com/

fabric-ca-client register --caname ca-org2 --id.name user0002 --id.secret userpw0002 --id.type client --id.attrs 'role=PHARMA:ecert' --tls.certfiles "${PWD}/organizations/fabric-ca/org2/tls-cert.pem"

fabric-ca-client enroll -u https://user0002:userpw0002@localhost:8054 --caname ca-org2 -M "${PWD}/organizations/peerOrganizations/org2.example.com/users/user0002@org2.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org2/tls-cert.pem"

cp ${PWD}/organizations/peerOrganizations/org2.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org2.example.com/users/user0002@org2.example.com/msp/config.yaml
