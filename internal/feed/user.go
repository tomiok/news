package feed

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/mail"
)

const tokenMinLength = 40

type UserStorage interface {
	Save(email, token string) (string, error)
	Validate(token string) error

	SaveSource(sourceURL, place string, userID int) error
}

type UserService struct {
	UserStorage
}

func NewUserService(storage UserStorage) *UserService {
	return &UserService{storage}
}

func (u *UserService) UserLogin(emailOrToken string) (string, error) {
	if !isValidEmail(emailOrToken) {
		if len(emailOrToken) > tokenMinLength {
			token := emailOrToken
			return token, u.Validate(token)
		} else {
			return "", errors.New("not email/not token")
		}
	}

	token, err := generateToken()
	if err != nil {
		return "", err
	}

	return u.Save(emailOrToken, token)
}

func (u *UserService) SaveNewSource(sourceURL, place string, userID int) error {
	return u.SaveSource(sourceURL, place, userID)
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func generateToken() (string, error) {
	// Generate 32 bytes of random data
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Convert to base64 for easier handling
	return base64.URLEncoding.EncodeToString(b), nil
}
