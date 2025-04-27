package password

import (
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

func Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Logger.Errorf("Error hashing password: %v", err)
		return "", err
	}
	return string(hash), nil
}
