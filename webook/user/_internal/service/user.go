package service

import (
	"context"
	"golang.org/x/crypto/bcrypt"
	"start/webook/pkg/e"
	"start/webook/user/_internal/domain"
	"start/webook/user/_internal/repository"
)

type UserService interface {
	Signup(ctx context.Context, user domain.User) (id int, err error)
	Profile(ctx context.Context, uid int) (domain.User, error)
	LoginEmail(ctx context.Context, email string, password string) (int, error)
	Edit(ctx context.Context, u domain.User) error
	LoginByPhone(ctx context.Context, phone string) (int, error)
}

type userServiceImpl struct {
	repo repository.UserRepo
}

func (svc userServiceImpl) LoginByPhone(ctx context.Context, phone string) (int, error) {
	return svc.repo.LoginByPhone(ctx, phone)
}

func (svc userServiceImpl) Edit(ctx context.Context, u domain.User) error {
	return svc.repo.Edit(ctx, u)
}

func (svc userServiceImpl) LoginEmail(ctx context.Context, email string, password string) (uid int, err error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err != nil {
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return 0, e.NewErr(e.UserAuthFailed, "svc 账号密码错误", err.Error())
	}
	return 0, err
}

func NewUserServiceImpl(repo repository.UserRepo) UserService {
	return &userServiceImpl{
		repo: repo,
	}
}
func (svc userServiceImpl) Profile(ctx context.Context, uid int) (domain.User, error) {
	return svc.repo.Profile(nil, uid)
}
func (svc userServiceImpl) Signup(ctx context.Context, user domain.User) (id int, err error) {
	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user.Password = string(password)
	return svc.repo.Create(ctx, user)
}
