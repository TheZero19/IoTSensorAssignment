package Hash

import "golang.org/x/crypto/bcrypt"

func GetHashPSK(psk string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(psk), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPSK(hashedPSK, providedPSK string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPSK), []byte(providedPSK))
	return err == nil
}
