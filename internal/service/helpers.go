package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func compareHashAndPassword(hash, inputPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(inputPassword)); err != nil {
		return fmt.Errorf("failed to compare hash and input password, err: %w", err)
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(inputPassword))
}
