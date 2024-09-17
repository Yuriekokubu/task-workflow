package user

import (
	"errors"
	"github.com/Yuriekokubu/workflow/internal/auth"
	"github.com/Yuriekokubu/workflow/internal/model"
	"log"

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

func (service Service) Login(req model.RequestLogin) (string, uint, string, error) {
	user, err := service.Repository.FindOneByUsername(req.Username)
	if err != nil {
		return "", 0, "", errors.New("invalid user or password")
	}

	if ok := checkPasswordHash(req.Password, user.Password); !ok {
		return "", 0, "", errors.New("invalid user or password")
	}

	token, err := auth.CreateToken(user.Username, service.secret)
	if err != nil {
		log.Println("Fail to create token")
		return "", 0, "", errors.New("something went wrong")
	}
	return user.Username, user.ID, token, nil
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
