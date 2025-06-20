package utils

import (
	"fmt"
	"payslip-generator-service/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtUtil struct {
	Config *config.Config
}

func NewJwtUtil(config *config.Config) *JwtUtil {
	return &JwtUtil{
		Config: config,
	}
}

// Claims structure for JWT
type Claims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates a new access token
func (j *JwtUtil) GenerateAccessToken(ID string) (string, error) {
	claims := &Claims{
		ID: ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.Config.Security.Jwt.AccessTokenLifetime) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Config.Security.Jwt.AccessTokenSecret))
}

// GenerateRefreshToken generates a new refresh token
func (j *JwtUtil) GenerateRefreshToken(ID string) (string, error) {
	claims := &Claims{
		ID: ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.Config.Security.Jwt.RefreshTokenLifetime) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Config.Security.Jwt.RefreshTokenSecret))
}

// ValidateToken validates the token string and returns the claims
func (j *JwtUtil) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Config.Security.Jwt.AccessTokenSecret), nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

// ValidateRefreshToken validates the refresh token string and returns the claims
func (j *JwtUtil) ValidateRefreshToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.Config.Security.Jwt.RefreshTokenSecret), nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

// RefreshAccessToken generates a new access token using a valid refresh token
func (j *JwtUtil) RefreshAccessToken(refreshToken string) (string, error) {
	claims, err := j.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token")
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return "", fmt.Errorf("refresh token expired")
	}

	// Generate a new access token
	return j.GenerateAccessToken(claims.ID)
}
