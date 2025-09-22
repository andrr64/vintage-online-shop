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

func (s *userService) FindByID(userId uint) (*models.User, error) {
	user, err := s.repo.FindByID(userId)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userService) UpdateAccount(userID uint, data dto.BodyUpdateAccountDTO) (dto.ResponseUserInfoDTO, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return dto.ResponseUserInfoDTO{}, errors.New("user not found")
	}

	// update fullname jika diisi
	if data.Fullname != "" && data.Fullname != user.Fullname {
		user.Fullname = data.Fullname
	}

	// cek username jika diisi
	if data.Username != "" && data.Username != user.Username {
		exists, _ := s.repo.IsUsernameExists(data.Username)
		if exists {
			return dto.ResponseUserInfoDTO{}, errors.New("username already taken")
		}
		user.Username = data.Username
	}

	// cek email jika diisi
	if data.Email != "" && data.Email != user.Email {
		exists, _ := s.repo.IsEmailExists(data.Email)
		if exists {
			return dto.ResponseUserInfoDTO{}, errors.New("email already taken")
		}
		user.Email = data.Email
	}

	// simpan perubahan
	if err := s.repo.Update(user); err != nil {
		return dto.ResponseUserInfoDTO{}, errors.New("failed to update account")
	}

	return dto.ResponseUserInfoDTO{
		Username: user.Username,
		Email:    user.Email,
		Fullname: user.Fullname,
	}, nil
}
