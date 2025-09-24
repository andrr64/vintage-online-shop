package services

import (
	"net/http"

	"errors"
	"vintage-server/dto"
	"vintage-server/helpers"
	"vintage-server/models"
	"vintage-server/repositories"
	security_jwt "vintage-server/security/jwt"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(input dto.InputRegisterDTO) (dto.ResponseRegisterUserDTO, error) {
	// cek username/email
	if existing, _ := s.repo.FindByUsernameOrEmail(input.Username, input.Email); existing != nil {
		return dto.ResponseRegisterUserDTO{}, errors.New("username or email already taken")
	}

	// hash password
	hashedPassword, err := helpers.GeneratePasswordHash(input.Password)
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

func (s *userService) Login(input dto.InputLoginDTO) (string, error) {
	user, err := s.repo.FindByEmail(input.Email)
	if err != nil {
		return "", errors.New("email or password salah")
	}

	// cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return "", errors.New("invalid password")
	}

	// generate JWT via helper
	token, err := security_jwt.CreateUserAccessToken(user.ID, user.Username)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (s *userService) FindByID(userId uint) (*models.User, error) {
	user, err := s.repo.FindByID(userId)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userService) UpdateAccount(userID uint, data dto.InputUpdateAccountDTO) (dto.ResponseUserInfoDTO, error) {
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

func (s *userService) UpdatePassword(userID uint, data dto.InputUpdatePasswordDTO) (int, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return http.StatusNotFound, errors.New("user not found")
	}

	if !helpers.ComparePassword(data.OldPassword, user.Password) {
		return http.StatusBadRequest, errors.New("old password is incorrect")
	}

	hashed, err := helpers.GeneratePasswordHash(data.NewPassword)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to hash new password")
	}

	user.Password = hashed
	if err := s.repo.Update(user); err != nil {
		return http.StatusInternalServerError, errors.New("failed to update password")
	}

	return http.StatusOK, nil
}
