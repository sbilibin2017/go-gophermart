package password

import "golang.org/x/crypto/bcrypt"

type Comparer struct{}

func NewComparer() *Comparer {
	return &Comparer{}
}

func (c *Comparer) Compare(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
