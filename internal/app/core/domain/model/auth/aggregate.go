package auth

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	userID       uuid.UUID
	email        string
	passwordHash string
}

func NewCredentials(userID uuid.UUID, email string, plainPassword string) (*Credentials, error) {
	hash, err := hashPassword(plainPassword)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	return &Credentials{
		userID:       userID,
		email:        email,
		passwordHash: hash,
	}, nil
}

func (c *Credentials) UserID() uuid.UUID    { return c.userID }
func (c *Credentials) Email() string        { return c.email }
func (c *Credentials) PasswordHash() string { return c.passwordHash }

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", fmt.Errorf("generate hash from password: %w", err)
	}
	return string(bytes), nil
}

func (c *Credentials) VerifyPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(c.passwordHash), []byte(plainPassword))
	return err == nil
}

type CredentialsProps struct {
	UserID       uuid.UUID
	Email        string
	PasswordHash string
}

func ReconstituteCredentials(props *CredentialsProps) *Credentials {
	return &Credentials{
		userID:       props.UserID,
		email:        props.Email,
		passwordHash: props.PasswordHash,
	}
}
