package password

import (
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"golang.org/x/crypto/bcrypt"
)

type Hasher struct {
	cost int
}

func NewHasher(cost int) *Hasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &Hasher{
		cost: cost,
	}
}

func (h *Hasher) Hash(password string) *string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		logger.Logger.Errorf("Error hashing password: %v", err)
		return nil
	}
	s := string(hash)
	return &s
}
