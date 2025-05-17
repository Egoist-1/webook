package dao

import (
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
	"webook/pkg/er"
)

type UserDAO interface {
	Insert(ctx context.Context, u User) (int64, error)
	FindById(ctx context.Context, uid int64) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	UpdateProfile(ctx context.Context, u User) error
	CreateOrFind(ctx context.Context, u User) (int64, error)
}

type userDaoImpl struct {
	db *gorm.DB
}

func (dao userDaoImpl) CreateOrFind(ctx context.Context, u User) (int64, error) {
	err := dao.db.WithContext(ctx).FirstOrCreate(&u, User{Phone: u.Phone}).Error
	return u.Id, err
}

func (dao userDaoImpl) UpdateProfile(ctx context.Context, u User) error {
	u.UTime = time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Where("id = ?", u.Id).Updates(u).Error
}

func (dao userDaoImpl) FindByEmail(ctx context.Context, email string) (u User, err error) {
	err = dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err == gorm.ErrRecordNotFound {
		return User{}, err.NewErr(err.UserAuthFailed, "账号密码错误", err.Error())
	}
	return
}

func (dao userDaoImpl) FindById(ctx context.Context, uid int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", uid).First(&u).Error
	return u, err
}

func (dao userDaoImpl) Insert(ctx context.Context, u User) (int64, error) {
	u.CTime = time.Now().Unix()
	u.UTime = time.Now().Unix()
	err := dao.db.WithContext(ctx).Model(&User{}).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		switch mysqlErr.Number {
		case 1062:
			return 0, err.NewErr(err.UserExist, "dao insert 用户已存在", "")
		}
	}
	return u.Id, err
}

func NewUserDao(db *gorm.DB) UserDAO {
	return &userDaoImpl{db: db}
}

type User struct {
	Id       int64 `gorm:"primaryKey"`
	Name     string
	Email    sql.NullString `gorm:"unique"`
	Phone    sql.NullString `gorm:"unique"`
	Avatar   string
	Password string
	AboutMe  string
	CTime    int64
	UTime    int64
}
