package src

func BlankPackagedPrescription() string {
	var prescription Prescription
	// Get current user pubkey
	pubkey, err := readLocalPubkey(currentUserObscure())
	if err != nil {
		panic(err)
	}
	packaged, err := packagePrescription(pubkey, &prescription)
	if err != nil {
		panic(err)
	}
	return packaged
}
