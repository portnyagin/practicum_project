package service

import (
	"context"
	"errors"
	"github.com/portnyagin/practicum_project/internal/app/dto"
	"github.com/portnyagin/practicum_project/internal/app/infrastructure"
	"github.com/portnyagin/practicum_project/internal/app/model"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	dbUser model.UserRepository
	log    *infrastructure.Logger
}

func NewAuthService(userRepo model.UserRepository, log *infrastructure.Logger) *AuthService {
	var target AuthService
	target.dbUser = userRepo
	target.log = log
	return &target
}

func (s *AuthService) hashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), 4)
	return string(bytes), err
}

func (s *AuthService) checkPasswordHash(pass string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}

func (s *AuthService) Register(ctx context.Context, user *dto.User) (*dto.User, error) {
	if user == nil {
		s.log.Debug("AuthService: Register. got nil user")
		return nil, dto.ErrBadParam
	}
	if (user.Login == "") || (user.Pass == "") {
		s.log.Warn("AuthService: Register. Validation error", zap.String("user", user.Login))
		return nil, dto.ErrBadParam
	}

	hp, err := s.hashPassword(user.Pass)
	if err != nil {
		s.log.Error("AuthService: Register. Can't calculate hash", zap.String("login", user.Login), zap.Error(err))
		return nil, err
	}
	user.ID, err = s.dbUser.Save(ctx, user.Login, hp)
	if err != nil {
		s.log.Error("AuthService: Register. Can't register user", zap.String("login", user.Login), zap.Error(err))
		return nil, err
	}
	return user, nil
}

// in success case return value is *user, nil. If *user is nil, then access denied
func (s *AuthService) Check(ctx context.Context, user *dto.User) (*dto.User, error) {
	if user == nil {
		s.log.Debug("AuthService: Check. got nil user")
		return nil, dto.ErrBadParam
	}
	if user.Login == "" {
		s.log.Warn("AuthService: Check. Validation error", zap.String("user", user.Login))
		return nil, dto.ErrBadParam
	}
	modelUser, err := s.dbUser.GetUserByLogin(ctx, user.Login)
	if err != nil {
		if errors.Is(err, &model.NoRowFound) {
			s.log.Debug("AuthService: Check. user not found in db", zap.String("login", user.Login))
			return nil, nil
		} else {
			s.log.Error("AuthService: Check.", zap.Error(err))
			return nil, err
		}
	}
	if s.checkPasswordHash(user.Pass, modelUser.Pass) {
		user.ID = modelUser.ID
		return user, nil
	} else {
		return nil, nil
	}
}
