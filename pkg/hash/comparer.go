package hash

import "golang.org/x/crypto/bcrypt"

type Comparer struct{}

func NewComarer() *Comparer {
	return &Comparer{}
}

func (hc *Comparer) Compare(hashedPassword, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}
