package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	secret            []byte
	issuer            string
	accessExpireMin   time.Duration
	refreshExpireHour time.Duration
}

type Claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

func NewManager(secret string, issuer string, accessExpireMin int, refreshExpireHour int) *Manager {
	return &Manager{
		secret:            []byte(secret),
		issuer:            issuer,
		accessExpireMin:   time.Duration(accessExpireMin) * time.Minute,
		refreshExpireHour: time.Duration(refreshExpireHour) * time.Hour,
	}
}

func (m *Manager) GeneratePair(userID int64, email string) (string, string, error) {
	accessToken, err := m.generateToken(userID, email, "access", m.accessExpireMin)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := m.generateToken(userID, email, "refresh", m.refreshExpireHour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (m *Manager) Parse(token string) (*Claims, error) {
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (m *Manager) generateToken(userID int64, email string, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   email,
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}
