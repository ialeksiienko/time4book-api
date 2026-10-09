package jwt

import (
	"testing"
	"time"
	outports "time4book/internal/app/core/ports/out"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestGenerateTokens(t *testing.T) {
	manager := NewManager("test-secret", 15*time.Minute, 24*time.Hour)
	userID := uuid.New()
	before := time.Now()
	tokens, err := manager.GenerateTokens(userID, "owner")
	if err != nil {
		t.Fatal(err)
	}
	after := time.Now()

	for _, tc := range []struct {
		name     string
		token    *outports.Token
		duration time.Duration
	}{
		{outports.AccessToken, tokens.AccessToken, 15 * time.Minute},
		{outports.RefreshToken, tokens.RefreshToken, 24 * time.Hour},
	} {
		t.Run(tc.name, func(t *testing.T) {
			parsed, err := jwtlib.Parse(tc.token.Value, func(*jwtlib.Token) (any, error) {
				return []byte("test-secret"), nil
			}, jwtlib.WithValidMethods([]string{"HS256"}))
			if err != nil {
				t.Fatal(err)
			}
			claims := parsed.Claims.(jwtlib.MapClaims)
			if claims["type"] != tc.name || claims["sub"] != userID.String() || claims["role"] != "owner" {
				t.Fatalf("unexpected claims: %v", claims)
			}
			if tc.token.ExpiresAt.Before(before.Add(tc.duration)) || tc.token.ExpiresAt.After(after.Add(tc.duration)) {
				t.Fatalf("unexpected expiry: %v", tc.token.ExpiresAt)
			}
			if claims["exp"] != float64(tc.token.ExpiresAt.Unix()) {
				t.Fatal("token expiry does not match returned expiry")
			}
		})
	}

	if tokens.AccessToken.Value == tokens.RefreshToken.Value {
		t.Error("access and refresh tokens must differ")
	}
	if gotID, role, err := manager.ValidateToken(tokens.AccessToken.Value); err != nil || gotID != userID || role != "owner" {
		t.Fatalf("access token rejected: %v, %v, %v", gotID, role, err)
	}
	if _, _, err := manager.ValidateToken(tokens.RefreshToken.Value); err == nil {
		t.Error("refresh token accepted as an access token")
	}
}
