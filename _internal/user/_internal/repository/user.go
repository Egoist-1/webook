package repository

import (
	"context"
	"database/sql"
	"github.com/jinzhu/copier"
	"webook/_internal/user/_internal/domain"
	"webook/_internal/user/_internal/repository/dao"
)

type UserRepo interface {
	Create(ctx context.Context, user domain.User) (int64, error)
	Profile(ctx context.Context, uid int64) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	Edit(ctx context.Context, u domain.User) error
	LoginByPhone(ctx context.Context, phone string) (int64, error)
}

type userRepoImpl struct {
	dao dao.UserDAO
}

func (repo userRepoImpl) LoginByPhone(ctx context.Context, phone string) (int64, error) {
	var u dao.User
	u.Phone = sql.NullString{
		String: phone,
		Valid:  true,
	}
	return repo.dao.CreateOrFind(ctx, u)
}

func (repo userRepoImpl) Edit(ctx context.Context, u domain.User) error {
	return repo.dao.UpdateProfile(ctx, repo.domainToEntity(u))
}

func (repo userRepoImpl) FindByEmail(ctx context.Context, email string) (u domain.User, err error) {
	res, err := repo.dao.FindByEmail(ctx, email)
	copier.Copy(&u, &res)
	return

}

func (repo userRepoImpl) Profile(ctx context.Context, uid int64) (domain.User, error) {
	user, err := repo.dao.FindById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	var res domain.User
	err = copier.Copy(&res, &user)
	return res, err
}

func NewUserRepo(dao dao.UserDAO) UserRepo {
	return &userRepoImpl{
		dao: dao,
	}
}

func (repo userRepoImpl) Create(ctx context.Context, user domain.User) (int64, error) {
	return repo.dao.Insert(ctx, repo.domainToEntity(user))
}
func (repo userRepoImpl) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id:   u.Id,
		Name: u.Name,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		AboutMe:  u.AboutMe,
		CTime:    u.CTime,
		UTime:    0,
	}
}
