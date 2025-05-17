package auth

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"links-service/internal/user"
	"log"
)

type AuthService struct {
	userRepository *user.UserRepository
}

func NewAuthService(repo *user.UserRepository) *AuthService {
	return &AuthService{
		userRepository: repo,
	}
}

func (s *AuthService) Login(email, password string) (string, error) {
	byEmail, err := s.userRepository.FindByEmail(email)
	if err != nil {
		log.Println(err)
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(byEmail.Password), []byte(password))
	if err != nil {
		log.Println(err)
		return "", errors.New(ErrWrongCredential)
	}

	return byEmail.Email, nil
}

func (s *AuthService) Register(name, email, password string) (string, error) {
	existedUser, err := s.userRepository.FindByEmail(email)
	if err == nil || existedUser != nil {
		return "", errors.New(ErrUserAlreadyExists + email)
	}

	passwordCrypt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &user.User{
		Email:    email,
		Password: string(passwordCrypt),
		Name:     name,
	}
	s.userRepository.Create(user)
	return user.Email, nil
}
