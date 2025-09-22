package services

import (
	"errors"
	"vintage-server/dto"
	"vintage-server/models"
	"vintage-server/repositories"

	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(input dto.RegisterUserDTO) (dto.ResponseRegisterUserDTO, error) {
	// cek username/email
	if existing, _ := s.repo.FindByUsernameOrEmail(input.Username, input.Email); existing != nil {
		return dto.ResponseRegisterUserDTO{}, errors.New("username or email already taken")
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.ResponseRegisterUserDTO{}, errors.New("failed to hash password")
	}

	// simpan user
	user := models.User{
		Username: input.Username,
		Fullname: input.Fullname,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := s.repo.Create(&user); err != nil {
		return dto.ResponseRegisterUserDTO{}, errors.New("failed to create user")
	}

	
	return dto.ResponseRegisterUserDTO{
		ID:       user.ID,
		Username: user.Username,
		Fullname: user.Fullname,
		Email:    user.Email,
	}, nil
}

func (s *userService) Login(input dto.LoginUserDTO) (string, error) {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return "", errors.New("email or password salah")
	}

	// cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", errors.New("invalid password")
	}

	// ambil secret JWT
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT secret not set")
	}

	// generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return tokenString, nil
}
