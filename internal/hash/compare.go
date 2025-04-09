package hash

import "golang.org/x/crypto/bcrypt"

func CompareHashAndPassword(hashedPassword, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}
