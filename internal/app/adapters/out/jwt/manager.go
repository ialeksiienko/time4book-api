package jwt

import (
	"errors"
	"fmt"
	"time"
	"time4book/internal/app/core/ports"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTManager struct {
	secretKey            string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewManager(
	secret string,
	accessTokenDuration time.Duration,
	refreshTokenDuration time.Duration,
) *JWTManager {
	return &JWTManager{
		secretKey:            secret,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

func (j *JWTManager) GenerateToken(userID uuid.UUID, role string, tokenType string) (*ports.Token, error) {
	dur := j.accessTokenDuration
	if tokenType == ports.RefreshToken {
		dur = j.refreshTokenDuration
	}

	expiresAt := time.Now().Add(dur)
	claims := jwt.MapClaims{
		"sub":  userID.String(),
		"role": role,
		"type": string(tokenType),
		"exp":  expiresAt.Unix(),
		"iat":  time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	val, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return nil, fmt.Errorf("sign token: %w", err)
	}

	return &ports.Token{
		Value:     val,
		ExpiresAt: expiresAt,
	}, nil
}

func (j *JWTManager) ValidateToken(tokenStr string) (uuid.UUID, string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return uuid.Nil, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["type"] != string(ports.AccessToken) {
			return uuid.Nil, "", errors.New("invalid token type")
		}

		sub, ok := claims["sub"].(string)
		if !ok || sub == "" {
			return uuid.Nil, "", errors.New("invalid token subject")
		}

		userID, err := uuid.Parse(sub)
		if err != nil {
			return uuid.Nil, "", fmt.Errorf("invalid token subject: %w", err)
		}

		role, _ := claims["role"].(string)
		return userID, role, nil
	}

	return uuid.Nil, "", errors.New("invalid token")
}
