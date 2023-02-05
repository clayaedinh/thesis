package scripts

import "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"

const (
	USER_DOCTOR     string = "DOCTOR"
	USER_PATIENT           = "PATIENT"
	USER_PHARMACIST        = "PHARMA"
)

// Creates user for given username, secret, Certificate Authority, and role.
// Role: (USER_DOCTOR, USER_PATIENT, USER_PHARMACIST)
func CreateUser(username string, secret string, caname string, role string) {
	// Step 1: create MSP client
	//client, err := msp.New()

	// Step 2: Registration

	// create attribute for registration
	attribute := msp.Attribute{
		Name:  "role",
		Value: role,
		ECert: true,
	}
	// create registration request
	regRequest := msp.RegistrationRequest{
		Name:       username,
		Type:       "user",
		Attributes: []msp.Attribute{attribute},
		CAName:     caname,
		Secret:     secret,
	}

}
