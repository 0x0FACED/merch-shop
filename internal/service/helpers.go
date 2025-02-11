package service

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrFailedComparingHashAndPassword = errors.New("failed to compare hash and pass")
)

func compareHashAndPassword(hash, inputPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(inputPassword)); err != nil {
		return fmt.Errorf("%w: %w", ErrFailedComparingHashAndPassword, err)
	}
	return nil
}

func calcHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
