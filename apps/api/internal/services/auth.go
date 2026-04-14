package services

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"simon/apps/api/internal/auth"
	"simon/apps/api/internal/models"
	"simon/apps/api/internal/repositories"
)

var (
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrEmailTaken         = errors.New("email já cadastrado")
)

type AuthService struct {
	users     *repositories.UserRepository
	jwtSecret string
}

func NewAuthService(db *gorm.DB, jwtSecret string) *AuthService {
	return &AuthService{
		users:     repositories.NewUserRepository(db),
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Register(name, email, password string) (*models.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	existing, err := s.users.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		Name:     strings.TrimSpace(name),
		Email:    email,
		Password: string(hash),
	}
	if err := s.users.Create(user); err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	user, err := s.users.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}
	token, err := auth.SignAccessToken(user.ID, s.jwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}
