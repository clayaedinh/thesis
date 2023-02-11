package src

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	orgName      string
	userId       string
	mspId        string
	orgUrl       string
	userUrl      string
	cryptoPath   string
	certPath     string
	keyPath      string
	tlsCertPath  string
	peerEndpoint string
	gatewayPeer  string
	channelName  string
)

const chaincodeName = "rsa"

func SetConnectionVariables(newOrg string, newUserId string, peerPort string) {
	orgName = newOrg
	userId = newUserId
	mspId = cases.Title(language.Und).String(orgName) + "MSP"
	orgUrl = orgName + ".example.com"
	userUrl = userId + "@" + orgUrl
	cryptoPath = "../../test-network/organizations/peerOrganizations/" + orgUrl
	certPath = cryptoPath + "/users/" + userUrl + "/msp/signcerts/cert.pem"
	keyPath = cryptoPath + "/users/" + userUrl + "/msp/keystore/"
	tlsCertPath = cryptoPath + "/peers/peer0." + orgUrl + "/tls/ca.crt"
	peerEndpoint = peerPort
	gatewayPeer = "peer0." + orgUrl
	channelName = "mychannel"
}

func PrintConnectionVariables() {
	fmt.Println("==========================")
	fmt.Printf("orgName: %v\n", orgName)
	fmt.Printf("userId: %v\n", userId)
	fmt.Printf("mspId: %v\n", mspId)
	fmt.Printf("orgUrl: %v\n", orgUrl)
	fmt.Printf("userUrl: %v\n", userUrl)
	fmt.Printf("cryptoPath: %v\n", cryptoPath)
	fmt.Printf("certPath: %v\n", certPath)
	fmt.Printf("keyPath: %v\n", keyPath)
	fmt.Printf("tlsCertPath: %v\n", tlsCertPath)
	fmt.Printf("peerEndpoint: %v\n", peerEndpoint)
	fmt.Printf("gatewayPeer: %v\n", gatewayPeer)
	fmt.Printf("channelName: %v\n", channelName)
	fmt.Println("==========================")
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func NewGrpcConnection() (*grpc.ClientConn, error) {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate:%v", err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	return connection, nil
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func newIdentity() *identity.X509Identity {
	certificate, err := loadCertificate(certPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspId, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// newSign creates a function that generates a digital signature from a message digest using a private key.

func loadSignature(filepath string) (crypto.PrivateKey, error) {
	files, err := os.ReadDir(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key directory: %w", err)
	}
	privateKeyPEM, err := os.ReadFile(path.Join(filepath, files[0].Name()))
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}
	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func newSign() identity.Sign {
	privateKey, err := loadSignature(keyPath)
	if err != nil {
		panic(err)
	}
	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}
	return sign
}

func DefaultGateway(clientConnection *grpc.ClientConn) (*client.Gateway, error) {
	id := newIdentity()
	sign := newSign()

	// Create a Gateway connection for a specific client identity
	return client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		// Default timeouts for different gRPC calls
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)

}

func SmartContract(gw *client.Gateway) *client.Contract {
	network := gw.GetNetwork(channelName)
	return network.GetContract(chaincodeName)
}

func getCurrentUser() string {
	return userId
}

func getCurrentOrg() string {
	return orgName
}
