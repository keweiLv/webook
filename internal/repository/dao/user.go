package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱已注册")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	UpdateById(ctx context.Context, u User) error
	FindById(ctx context.Context, id int64) (User, error)
}

type GORMUserDao struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDao{
		db: db,
	}
}

func (dao *GORMUserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrUserDuplicateEmail
		}
	}
	return err
}

func (dao *GORMUserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *GORMUserDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}

func (dao *GORMUserDao) UpdateById(ctx context.Context, u User) error {
	err := dao.db.WithContext(ctx).Model(&u).Updates(User{Id: u.Id, Nickname: u.Nickname, Birthday: u.Birthday, Profile: u.Profile}).Error
	return err
}

func (dao *GORMUserDao) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("Id = ?", id).First(&u).Error
	return u, err
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Birthday == 0 {
		tx.Statement.SetColumn("Birthday", nil)
	}
	return nil
}

// 对应数据库表 PO（persistent object）
type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string         `gorm:"not null"`
	Birthday int64
	Nickname string         `gorm:"null"`
	Profile  string         `gorm:"null"`
	Phone    sql.NullString `gorm:"unique"`

	// 毫秒时间
	Ctime int64
	Utime int64
}
