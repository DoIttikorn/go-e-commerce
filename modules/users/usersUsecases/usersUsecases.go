package usersUsecases

import (
	"fmt"

	"github.com/Doittikorn/go-e-commerce/config"
	"github.com/Doittikorn/go-e-commerce/modules/users"
	"github.com/Doittikorn/go-e-commerce/modules/users/usersRepositories"
	"github.com/Doittikorn/go-e-commerce/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type UsersUsecasesImpl interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error)
	GetPassport(req *users.UserCredential) (*users.UserPassport, error)
	RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error)
	DeleteOauth(oauthId string) error
	GetUserProfile(userId string) (*users.User, error)

	InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error)
}

type usersUsecase struct {
	cfg             config.ConfigImpl
	usersRepository usersRepositories.UsersRepositoriesImpl
}

func New(cfg config.ConfigImpl, userRepository usersRepositories.UsersRepositoriesImpl) UsersUsecasesImpl {
	return &usersUsecase{
		cfg:             cfg,
		usersRepository: userRepository,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.UserPassport, error) {

	// hash password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	// insert user
	result, err := u.usersRepository.InsertUser(req, false)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *usersUsecase) InsertAdmin(req *users.UserRegisterReq) (*users.UserPassport, error) {

	// hash password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	// insert user
	result, err := u.usersRepository.InsertUser(req, true)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *usersUsecase) GetPassport(req *users.UserCredential) (*users.UserPassport, error) {
	user, err := u.usersRepository.FindOneUserByEmail(req.Email)

	if err != nil {
		return nil, err
	}
	// compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}
	// Sign Token
	accessToken, err := auth.New(auth.Access, u.cfg.JWT(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})
	if err != nil {
		return nil, fmt.Errorf("sign token failed")
	}

	// create refresh token
	refreshToken, err := auth.New(auth.Access, u.cfg.JWT(), &users.UserClaims{
		Id:     user.Id,
		RoleId: user.RoleId,
	})
	if err != nil {
		return nil, fmt.Errorf("sign token failed")
	}
	// Set passport
	passport := &users.UserPassport{
		User: &users.UserResponse{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			RoleId:   user.RoleId,
		},
		Token: &users.UserToken{
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken.SignToken(),
		},
	}
	if err := u.usersRepository.InsertOauth(passport); err != nil {
		return nil, err
	}
	return passport, nil
}

func (u *usersUsecase) RefreshPassport(req *users.UserRefreshCredential) (*users.UserPassport, error) {
	// Parse token
	claims, err := auth.ParseToken(u.cfg.JWT(), req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// check oauth
	oauth, err := u.usersRepository.FindOneOauth(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Find profile
	profile, err := u.usersRepository.GetProfile(oauth.UserId)

	if err != nil {
		return nil, err
	}

	newClaims := &users.UserClaims{
		Id:     profile.Id,
		RoleId: profile.RoleId,
	}

	// create token
	accessToken, err := auth.New(
		auth.Access,
		u.cfg.JWT(),
		newClaims,
	)
	if err != nil {
		return nil, err
	}

	// create refresh token
	refreshToken := auth.RepeatToken(
		u.cfg.JWT(),
		newClaims,
		claims.ExpiresAt.Unix(),
	)

	if err != nil {
		return nil, err
	}

	newPassport := &users.UserPassport{
		User: &users.UserResponse{
			Id:       profile.Id,
			Email:    profile.Email,
			Username: profile.Username,
			RoleId:   profile.RoleId,
		},
		Token: &users.UserToken{
			Id:           oauth.Id,
			AccessToken:  accessToken.SignToken(),
			RefreshToken: refreshToken,
		},
	}

	// update oauth
	if err := u.usersRepository.UpdateOauth(newPassport.Token); err != nil {
		return nil, err
	}

	return newPassport, nil
}

func (u *usersUsecase) DeleteOauth(oauthId string) error {
	if err := u.usersRepository.DeleteOauth(oauthId); err != nil {
		return err
	}
	return nil
}

func (u *usersUsecase) GetUserProfile(userId string) (*users.User, error) {
	profile, err := u.usersRepository.GetProfile(userId)
	if err != nil {
		return nil, err
	}
	return profile, nil
}
