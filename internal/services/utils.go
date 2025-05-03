package services

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type validationError struct {
	tag     string
	Message string
}

func (v *validationError) Tag() string {
	return v.tag
}

func formatValidationError(err error) *validationError {
	if errs, ok := err.(validator.ValidationErrors); ok {
		e := errs[0]

		field := camelToSnake(e.Field())

		var message string
		switch e.Tag() {
		case "required":
			message = fmt.Sprintf("Field %s is required", field)
		case "gt":
			message = fmt.Sprintf("Field %s must be greater than %s", field, e.Param())
		case "oneof":
			message = fmt.Sprintf("Field %s must be one of [%s]", field, e.Param())
		default:
			message = fmt.Sprintf("Field %s is invalid (%s)", field, e.Tag())
		}

		return &validationError{
			tag:     e.Tag(),
			Message: message,
		}
	}
	return nil
}

func camelToSnake(s string) string {
	snake := regexp.MustCompile("([A-Z]+)([A-Z][a-z])").ReplaceAllString(s, "${1}_${2}")
	snake = regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func comparePassword(hashedPassword string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func generateTokenString(
	jwtSecretKey string,
	jwtExp time.Duration,
	issuer string,
	login string,
) (string, error) {
	claims := struct {
		jwt.RegisteredClaims
		Login string `json:"login"`
	}{
		Login: login,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}
	return signedToken, nil
}
