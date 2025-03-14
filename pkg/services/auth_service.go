package services

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/resend/resend-go/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserServices struct {
	secretKey    []byte
	refreshKey   []byte
	resendApiKey string
}

type Auth interface {
	GenerateToken(UserID uuid.UUID, email string, userRole string, isRefresh bool) (string, string, error)
	ValidateToken(token string, isRefresh bool) (*jwt.Token, error)
	HashPassword(password string) (string, error)
	CheckPassword(userPass string, password string) error
	SendEmail(email string) (bool, error)
}

func NewUserServices(secretKey, refreshKey []byte, resendApiKey string) Auth {
	return &UserServices{
		secretKey:    secretKey,
		refreshKey:   refreshKey,
		resendApiKey: resendApiKey,
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

func (s *UserServices) GenerateToken(userID uuid.UUID, email string, userRole string, isRefresh bool) (string, string, error) {
	secret := s.secretKey
	if isRefresh {
		secret = s.refreshKey
	}
	fmt.Println(userID)
	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    userRole,
		"sub":     email,
		"iss":     "todoApp",
		"exp":     time.Now().Add(time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})
	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": userRole,
		"sub":  email,
		"exp":  time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	accessToken, AtErr := access_token.SignedString(secret)
	refreshToken, RfErr := refresh_token.SignedString(secret)
	if AtErr != nil || RfErr != nil {
		return "", "", fmt.Errorf("error generating tokens: %v %v", AtErr, RfErr)
	}
	return accessToken, refreshToken, nil
}

func (s *UserServices) ValidateToken(token string, isRefresh bool) (*jwt.Token, error) {
	secret := s.secretKey
	if isRefresh {
		secret = s.refreshKey
	}
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return secret, nil
	})
}

func (s *UserServices) SendEmail(email string) (bool, error) {
	client := resend.NewClient(s.resendApiKey)
	fmt.Println([]string{email})
	fmt.Println(s.resendApiKey)
	params := &resend.SendEmailRequest{
		From:    "Notyz <notyz@resend.dev>",
		To:      []string{email},
		Html:    "<strong>hello world</strong>",
		Subject: "Hello from Golang",
		// Cc:      []string{"cc@example.com"},
		// Bcc:     []string{"bcc@example.com"},
		ReplyTo: "replyto@example.com",
	}
	sent, err := client.Emails.Send(params)
	if err != nil {
		return false, fmt.Errorf("%s", err.Error())
	}
	log.Printf("email sent to: %v", sent.Id)
	return true, nil
}
