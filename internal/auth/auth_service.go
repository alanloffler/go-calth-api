package auth

import (
	"errors"
	"time"

	"github.com/alanloffler/go-calth-api/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	cfg *config.Config
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenClaims struct {
	UserID     string `json:"userId"`
	BusinessID string `json:"businessID"`
	RoleID     string `json:"roleId"`
	jwt.RegisteredClaims
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{cfg: cfg}
}

func (s *AuthService) GenerateTokenPair(userID, businessID, roleID string) (*TokenPair, error) {
	accessToken, err := s.generateToken(userID, businessID, roleID, s.cfg.JwtSecret, s.cfg.JwtAccessExpiry)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateToken(userID, businessID, roleID, s.cfg.JwtRefreshSecret, s.cfg.JwtRefreshExpiry)
	if err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (s *AuthService) ValidateAccessToken(tokenStr string) (*TokenClaims, error) {
	return s.validateToken(tokenStr, s.cfg.JwtSecret)
}

func (s *AuthService) ValidateRefreshToken(tokenStr string) (*TokenClaims, error) {
	return s.validateToken(tokenStr, s.cfg.JwtRefreshSecret)
}

func (s *AuthService) generateToken(userID, businessID, roleID, secret, expiry string) (string, error) {
	duration, err := time.ParseDuration(expiry)
	if err != nil {
		return "", err
	}

	claims := TokenClaims{
		UserID:     userID,
		BusinessID: businessID,
		RoleID:     roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func (s *AuthService) validateToken(tokenStr, secret string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
