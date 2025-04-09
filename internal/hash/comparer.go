package hash

import "golang.org/x/crypto/bcrypt"

type HashComarer struct{}

func (hc *HashComarer) Compare(hashedPassword, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}
