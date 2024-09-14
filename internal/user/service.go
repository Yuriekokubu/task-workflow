package user

import (
	"errors"
	"log"
	"github.com/Yuriekokubu/workflow/internal/auth"
	"github.com/Yuriekokubu/workflow/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	Repository Repository
	secret     string
}

func NewService(db *gorm.DB, secret string) Service {
	return Service{
		Repository: NewRepository(db),
		secret:     secret,
	}
}

func (service Service) Login(req model.RequestLogin) (string, error) {
	user, err := service.Repository.FindOneByUsername(req.Username)
	if err != nil {
		return "", errors.New("Invalid user or password")
	}

	if ok := checkPasswordHash(req.Password, user.Password); !ok {
		return "", errors.New("Invalid user or password")
	}

	token, err := auth.CreateToken(user.Username, service.secret)
	if err != nil {
		log.Println("Fail to create token")
		return "", errors.New("Something went wrong")
	}
	return token, nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Exported function name
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
