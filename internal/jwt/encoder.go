package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type JWTEncoderConfig interface {
	GetSecretKey() string
	GetExp() time.Duration
}

type Encoder struct {
	c JWTEncoderConfig
}

func NewEncoder(c JWTEncoderConfig) *Encoder {
	return &Encoder{
		c: c,
	}
}

var (
	ErrTokenIsNotSigned = errors.New("token is not signed")
)

func (e *Encoder) Encode(claims *types.Claims) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": int64(claims.UserID.UserID),
			"exp":     (time.Now().Add(1 * e.c.GetExp())).Unix(),
			"iat":     time.Now().Unix(),
		},
	)
	signedToken, err := token.SignedString([]byte(e.c.GetSecretKey()))
	if err != nil {
		return "", ErrTokenIsNotSigned
	}
	return signedToken, nil
}
