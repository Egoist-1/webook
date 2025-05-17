package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"webook/_internal/user/_internal/domain"
	"webook/_internal/user/_internal/repository"
	"webook/pkg/er"
)

type UserService interface {
	Signup(ctx context.Context, user domain.User) (id int64, err error)
	Profile(ctx context.Context, uid int64) (domain.User, error)
	LoginEmail(ctx context.Context, email string, password string) (int64, error)
	Edit(ctx context.Context, u domain.User) error
	LoginByPhone(ctx context.Context, phone string) (int64, error)
}

type userServiceImpl struct {
	repo repository.UserRepo
}

func (svc userServiceImpl) LoginByPhone(ctx context.Context, phone string) (int64, error) {
	return svc.repo.LoginByPhone(ctx, phone)
}

func (svc userServiceImpl) Edit(ctx context.Context, u domain.User) error {
	return svc.repo.Edit(ctx, u)
}

func (svc userServiceImpl) LoginEmail(ctx context.Context, email string, password string) (uid int64, err error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err != nil {
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return 0, err.NewErr(err.UserAuthFailed, "svc 账号密码错误", err.Error())
	}
	return 0, err
}

func NewUserServiceImpl(repo repository.UserRepo) UserService {
	return &userServiceImpl{
		repo: repo,
	}
}
func (svc userServiceImpl) Profile(ctx context.Context, uid int64) (domain.User, error) {
	return svc.repo.Profile(nil, uid)
}
func (svc userServiceImpl) Signup(ctx context.Context, user domain.User) (id int64, err error) {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user.Password = string(password)
	return svc.repo.Create(ctx, user)
}
