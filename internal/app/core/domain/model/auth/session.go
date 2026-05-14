package auth

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	id           uuid.UUID
	userID       uuid.UUID
	refreshToken string
	expiresAt    time.Time
	createdAt    time.Time
	updatedAt    *time.Time
}

func NewSession(userID uuid.UUID, token string, expiresAt time.Time) *Session {
	return &Session{
		id:           uuid.New(),
		userID:       userID,
		refreshToken: token,
		expiresAt:    expiresAt,
		createdAt:    time.Now().UTC(),
	}
}

func RestoreSession(
	id uuid.UUID,
	userID uuid.UUID,
	token string,
	expiresAt time.Time,
	createdAt time.Time,
	updatedAt *time.Time,
) *Session {
	return &Session{
		id:           id,
		userID:       userID,
		refreshToken: token,
		expiresAt:    expiresAt,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

func (s *Session) ID() uuid.UUID         { return s.id }
func (s *Session) UserID() uuid.UUID     { return s.userID }
func (s *Session) RefreshToken() string  { return s.refreshToken }
func (s *Session) ExpiresAt() time.Time  { return s.expiresAt }
func (s *Session) CreatedAt() time.Time  { return s.createdAt }
func (s *Session) UpdatedAt() *time.Time { return s.updatedAt }

func (s *Session) IsExpired() bool {
	return time.Now().After(s.expiresAt)
}
