RED='\033[0;31m'
ORANGE='\033[0;33m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m'

TEST_NETWORK_PATH="../../test-network"
APP_RSA_PATH="../application/rsa"

ORG1PORT="localhost:7054"
ORG2PORT="localhost:8054"
ORG1PORT_PEER="localhost:7051"
ORG2PORT_PEER="localhost:9051"

printHelp () {
    echo "HELP"
}

createUser () {
    local USER_NUM=$1
    local ORG_NUM=$2
    local ROLE=$3
    echo "${CYAN}creating user${USER_NUM} in org${ORG_NUM} with role ${ROLE}.${NC}"
    cd ${TEST_NETWORK_PATH}
    createUserFabric $USER_NUM $ORG_NUM $ROLE
}

createUserFabric () {
    local USER_NUM=$1
    local ORG_NUM=$2
    local ROLE=$3
    export PATH=${PWD}/../bin:${PWD}:$PATH
    export FABRIC_CFG_PATH=$PWD/../config/
    export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/org${ORG_NUM}.example.com/
    fabric-ca-client register --caname ca-org${ORG_NUM} --id.name user${USER_NUM} --id.secret userpw${USER_NUM} --id.type client --id.attrs "role=${ROLE}:ecert" --tls.certfiles "${PWD}/organizations/fabric-ca/org${ORG_NUM}/tls-cert.pem"

    if [ "$ORG_NUM" = "1" ] || [ $ORG_NUM -eq 1 ]; then
        PORT=$ORG1PORT
    else
        PORT=$ORG2PORT
    fi

    fabric-ca-client enroll -u https://user${USER_NUM}:userpw${USER_NUM}@${PORT} --caname ca-org${ORG_NUM} -M "${PWD}/organizations/peerOrganizations/org${ORG_NUM}.example.com/users/user${USER_NUM}@org${ORG_NUM}.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org${ORG_NUM}/tls-cert.pem"
    cp ${PWD}/organizations/peerOrganizations/org${ORG_NUM}.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org${ORG_NUM}.example.com/users/user${USER_NUM}@org${ORG_NUM}.example.com/msp/config.yaml
}

createUsers () {
    createUser "0001" 1 "DOCTOR"
    createUser "0002" 1 "PATIENT"
    createUser "0003" 2 "PHARMA"  
    createUser "0004" 1 "READER"
}

if [ "$1" = "help" ] || [ $# -eq 0 ]; then
    printHelp
elif [ "$1" = "user" ]; then
    createUser $2 $3 $4
elif [ "$1" = "samples" ]; then
    createUsers
elif [ $# -gt 0 ]; then
    printError
fi