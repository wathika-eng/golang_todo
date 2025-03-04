package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserServices struct {
	secretKey []byte
}

type Auth interface {
	GenerateToken(UserID uint, email string) (string, string, error)
	ValidateToken(token string, isRefresh bool) (*jwt.Token, error)
	HashPassword(password string) (string, error)
	CheckPassword(userPass string, password string) error
}

func NewUserServices(secretKey []byte) Auth {
	return &UserServices{
		secretKey: secretKey,
	}
}

func (s *UserServices) HashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing the password: %v", err)
	}
	return string(hashedPass), nil
}

func (s *UserServices) CheckPassword(userPass string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(userPass), []byte(password))
}

func (s *UserServices) GenerateToken(userID uint, email string) (string, string, error) {
	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"sub":     email,
		"iss":     "todoApp",
		"exp":     time.Now().Add(time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})
	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": email,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	accessToken, At_err := access_token.SignedString(s.secretKey)
	refreshToken, Rf_err := refresh_token.SignedString(s.secretKey)
	if At_err != nil || Rf_err != nil {
		return "", "", fmt.Errorf("error generating tokens: %v %v", At_err, Rf_err)
	}
	return accessToken, refreshToken, nil
}

func (s *UserServices) ValidateToken(token string, isRefresh bool) (*jwt.Token, error) {
	if isRefresh {
		s.secretKey = []byte(token)
	}
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(s.secretKey), nil
	})
}
