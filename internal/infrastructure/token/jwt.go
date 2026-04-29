package token

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hogiabao7725/go-auth-playground/internal/domain/user"
)

var _ user.TokenProvider = (*JWT)(nil)

type accessClaims struct {
	jwt.RegisteredClaims
	Role string `json:"role"`
}

type refreshClaims struct {
	jwt.RegisteredClaims
}

// Technical errors
var (
	errEmptySecret = errors.New("secret cannot be empty")
	errInvalidTTL  = errors.New("TTL must be positive")
)

type JWT struct {
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
	signingMethod jwt.SigningMethod
}

func NewJWT(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration) *JWT {
	return &JWT{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
		signingMethod: jwt.SigningMethodHS256,
	}
}

func (j *JWT) GenerateAccessToken(userID, role string) (string, error) {
	if err := validateTokenInput(userID, j.accessSecret, j.accessTTL); err != nil {
		return "", err
	}

	now := time.Now()
	claims := accessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTTL)),
		},
		Role: role,
	}

	token := jwt.NewWithClaims(j.signingMethod, claims)
	return token.SignedString([]byte(j.accessSecret))
}

func (j *JWT) GenerateRefreshToken(userID string) (string, error) {
	if err := validateTokenInput(userID, j.refreshSecret, j.refreshTTL); err != nil {
		return "", err
	}

	now := time.Now()
	claims := refreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshTTL)),
		},
	}

	token := jwt.NewWithClaims(j.signingMethod, claims)
	return token.SignedString([]byte(j.refreshSecret))
}

func (j *JWT) ParseAccessToken(tokenString string) (*user.AccessTokenData, error) {
	token, err := jwt.ParseWithClaims(tokenString, &accessClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("infrastructure.token.jwt.ParseAccessToken: unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.accessSecret), nil
	})
	if err != nil {
		return nil, mapJWTError(err)
	}

	claims, ok := token.Claims.(*accessClaims)
	if !ok || !token.Valid {
		return nil, user.ErrTokenInvalid
	}

	return &user.AccessTokenData{
		UserID:    claims.Subject,
		Role:      claims.Role,
		ExpiresAt: claims.ExpiresAt.Time,
		TokenID:   claims.ID,
	}, nil
}

func (j *JWT) ParseRefreshToken(tokenString string) (*user.RefreshTokenData, error) {
	token, err := jwt.ParseWithClaims(tokenString, &refreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("infrastructure.token.jwt.ParseRefreshToken: unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.refreshSecret), nil
	})
	if err != nil {
		return nil, mapJWTError(err)
	}

	claims, ok := token.Claims.(*refreshClaims)
	if !ok || !token.Valid {
		return nil, user.ErrTokenInvalid
	}

	return &user.RefreshTokenData{
		UserID:    claims.Subject,
		ExpiresAt: claims.ExpiresAt.Time,
		TokenID:   claims.ID,
	}, nil
}

func (j *JWT) ValidateAccessToken(tokenString string) bool {
	_, err := j.ParseAccessToken(tokenString)
	return err == nil
}

func (j *JWT) ValidateRefreshToken(tokenString string) bool {
	_, err := j.ParseRefreshToken(tokenString)
	return err == nil
}

func (j *JWT) AccessTTL() time.Duration {
	return j.accessTTL
}

func (j *JWT) RefreshTTL() time.Duration {
	return j.refreshTTL
}

// --- Helper functions ---

// mapJWTError: Technical (jwt lib) → Domain errors
func mapJWTError(err error) error {
	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return user.ErrTokenExpired
	case errors.Is(err, jwt.ErrTokenMalformed),
		errors.Is(err, jwt.ErrSignatureInvalid):
		return user.ErrTokenInvalid
	default:
		return err
	}
}

// validateTokenInput: Mix of domain + technical validation
func validateTokenInput(userID, secret string, ttl time.Duration) error {
	if strings.TrimSpace(userID) == "" {
		return user.ErrEmptyID // Business rule → domain error
	}
	if strings.TrimSpace(secret) == "" {
		return errEmptySecret
	}
	if ttl <= 0 {
		return errInvalidTTL
	}
	return nil
}
