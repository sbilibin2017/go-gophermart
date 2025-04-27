package password

import (
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

func Compare(enteredPassword, storedPasswordHash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(enteredPassword))
	if err != nil {
		logger.Logger.Errorf("Password comparison failed: %v", err)
		return err
	}
	return nil
}
