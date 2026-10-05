package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"siakad-mini/app/model"
)

var (
	ErrInvalidToken = errors.New("token tidak valid")
	ErrExpiredToken = errors.New("token kedaluwarsa")
)

// claims adalah isi token: user_id, email, role, dan exp (lewat RegisteredClaims).
type claims struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(secret string, ttl time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), ttl: ttl}
}

// ExpiresInSeconds dipakai sebagai nilai expires_in pada response login.
func (m *JWTManager) ExpiresInSeconds() int { return int(m.ttl.Seconds()) }

func (m *JWTManager) Generate(u model.User) (string, error) {
	now := time.Now()
	c := claims{
		UserID: u.ID,
		Email:  u.Email,
		Role:   u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, c).SignedString(m.secret)
}

// Parse memeriksa tanda tangan, algoritma, dan masa berlaku token.
func (m *JWTManager) Parse(tokenString string) (model.AuthUser, error) {
	c := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, c, func(t *jwt.Token) (any, error) {
		// Tolak algoritma selain HMAC (mencegah serangan alg "none"/confusion).
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("algoritma tidak diharapkan: %v", t.Header["alg"])
		}
		return m.secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithExpirationRequired())
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return model.AuthUser{}, ErrExpiredToken
		}
		return model.AuthUser{}, ErrInvalidToken
	}
	if !token.Valid || c.UserID == 0 {
		return model.AuthUser{}, ErrInvalidToken
	}
	return model.AuthUser{UserID: c.UserID, Email: c.Email, Role: c.Role}, nil
}