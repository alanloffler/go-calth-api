package auth

import (
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
