package auth

import (
	"fmt"
	"time"
	"errors"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "skill-calculator-access-issuer"
)

func HashPassword(password string) (string, error) {
	byte_pass_hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(byte_pass_hash), nil
}

func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		//return errors.New("Passowrd not matching")
		return fmt.Errorf("Passowrd not matching %v", err)
	}
	return nil
}

func MakeJWT(user_id int, token_secret string, expires_in time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expires_in)),
		Issuer:    string(TokenTypeAccess),
		Subject:   strconv.Itoa(user_id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	result_string, err := token.SignedString([]byte(token_secret))
	if err != nil {
		return "", err
	}

	return result_string, nil
}

func ValidateJWT(token_string, token_secret string) (int, error) {
	claimsStruct := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(token_secret, &claimsStruct, func(token *jwt.Token) (interface{}, error) {
		return []byte(token_secret), nil
	})
	if err != nil {
		return 0, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return 0, err
	}
	if issuer != string(TokenTypeAccess) {
		return 0, errors.New("invalid issuer")
	}

	user_id_str, err := token.Claims.GetSubject()
	if err != nil {
		return 0, fmt.Errorf("token did not have user ID: %w", err)
	}

	return_user_id, err := strconv.Atoi(user_id_str)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID: %w", err)
	}
	return return_user_id, nil
}
