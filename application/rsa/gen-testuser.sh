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
    echo
    echo "${YELLOW}gen-testuser.sh${NC}"
    echo
    echo "Syntax: sh gen-testuser.sh <Mode>"
    echo
    echo "Available Modes"
    echo
    echo "${CYAN}user${NC} - creates single user with standard format"
    echo "Syntax: sh gen-testuser.sh user <user_num> <org_num> <role>"
    echo "Roles: DOCTOR, PATIENT, PHARMA"
    echo "Example: sh gen-testuser.sh user 0001 1 DOCTOR"
    echo 
    echo "${CYAN}samples${NC} - generates three users for use in thesis prescription network."
    echo "Syntax: sh gen-testuser.sh samples"
    echo
    echo "Users generated by samples mode:"
    echo "${ORANGE}user0001${NC}: a doctor under Organization 1."
    echo "${ORANGE}user0002${NC}: a patient under Organization 1."
    echo "${ORANGE}user0003${NC}: a pharmacist under Organization 2."
    echo
    echo "For both modes: In addition to the Fabric CA signatures found in test-network, RSA keys are"
    echo "also generated in keys/user_000x."
    echo
    echo "Other Modes"
    echo
    echo "${CYAN}deletekeys${NC} - delete keys store in /rsakeys"
    echo "Syntax: sh gen-testuser.sh deletekeys"
    echo 
    echo "${CYAN}admin${NC} - creates keys for a user named Admin"
    echo "Syntax: sh gen-testuser.sh admin"
    echo 
}

printError(){
    echo "${RED}Error: unexpected argument. Try 'sh gen-testuser.sh help'.${NC}"
}

createUserFabric () {
    local USER_NUM=$1
    local ORG_NUM=$2
    local ROLE=$3
    export PATH=${PWD}/../bin:${PWD}:$PATH
    export FABRIC_CFG_PATH=$PWD/../config/
    export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/org${ORG_NUM}.example.com/
    fabric-ca-client register --caname ca-org${ORG_NUM} --id.name user${USER_NUM} --id.secret userpw${USER_NUM} --id.type client --id.attrs 'role=${ROLE}:ecert' --tls.certfiles "${PWD}/organizations/fabric-ca/org${ORG_NUM}/tls-cert.pem"

    if [ "$ORG_NUM" = "1" ] || [ $ORG_NUM -eq 1 ]; then
        PORT=$ORG1PORT
    else
        PORT=$ORG2PORT
    fi

    fabric-ca-client enroll -u https://user${USER_NUM}:userpw${USER_NUM}@${PORT} --caname ca-org${ORG_NUM} -M "${PWD}/organizations/peerOrganizations/org${ORG_NUM}.example.com/users/user${USER_NUM}@org${ORG_NUM}.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org${ORG_NUM}/tls-cert.pem"
    cp ${PWD}/organizations/peerOrganizations/org${ORG_NUM}.example.com/msp/config.yaml ${PWD}/organizations/peerOrganizations/org${ORG_NUM}.example.com/users/user${USER_NUM}@org${ORG_NUM}.example.com/msp/config.yaml
}

oldGenKeyPair () {
    USERNAME="user${1}"
    [ ! -d $USERNAME ] && mkdir $USERNAME
    cd $USERNAME
    openssl genpkey -out privkey.pem -quiet -algorithm rsa -pkeyopt rsa_keygen_bits:2048
    openssl pkey -in privkey.pem -out pubkey.pem -pubout
    cd ..
}



storePublicKey () {
    local USER_NUM=$1
    local ORG_NUM=$2
    local PORT=
    if [ "$ORG_NUM" = "1" ] || [ $ORG_NUM -eq 1 ]; then
        PORT=$ORG1PORT_PEER
    else
        PORT=$ORG2PORT_PEER
    fi
    echo "./rsa -user=user${USER_NUM} -org=org${ORG_NUM} -port=${PORT} storekey user${USER_NUM}"
    ./rsa -user=Admin -org=org${ORG_NUM} -port=${PORT} storekey user${USER_NUM}
}

createUser () {
    local USER_NUM=$1
    local ORG_NUM=$2
    local ROLE=$3
    echo "${CYAN}creating user${USER_NUM} in org${ORG_NUM} with role ${ROLE}.${NC}"
    cd ${TEST_NETWORK_PATH}
    createUserFabric $USER_NUM $ORG_NUM $ROLE
    cd ${APP_RSA_PATH}
    [ ! -d rsakeys ] && mkdir rsakeys
    cd rsakeys
    #oldGenKeyPair $USER_NUM
    ./rsa genkey user${USER_NUM}
    cd ..
    storePublicKey $USER_NUM $ORG_NUM
}

oldCreateAdmin() {
    echo "${CYAN}creating admin in org1...${NC}"
    [ ! -d rsakeys ] && mkdir rsakeys
    cd rsakeys
    USERNAME="Admin"
    [ ! -d $USERNAME ] && mkdir $USERNAME
    cd $USERNAME
    openssl genpkey -out privkey.pem -quiet -algorithm rsa -pkeyopt rsa_keygen_bits:2048
    openssl pkey -in privkey.pem -out pubkey.pem -pubout
    cd ../..
    ./rsa -user=Admin -org=org1 -port=localhost:7051 storekey Admin
}

createAdmin(){
    ./rsa genkey Admin
    ./rsa storekey Admin
}

createUsers () {
    createUser "0001" 1 "DOCTOR"
    createUser "0002" 1 "PATIENT"
    createUser "0003" 2 "PHARMA"  
}

deleteKeys () {
    rm -r rsakeys
}

if [ "$1" = "help" ] || [ $# -eq 0 ]; then
    printHelp
elif [ "$1" = "user" ]; then
    createUser $2 $3 $4
elif [ "$1" = "samples" ]; then
    createUsers
elif [ "$1" = "deletekeys" ]; then
    deleteKeys
elif [ "$1" = "admin" ]; then
    createAdmin
elif [ $# -gt 0 ]; then
    printError
fi



