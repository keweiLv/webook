package repository

import (
	"context"
	"github.com/keweiLv/webook/internal/domain"
	"github.com/keweiLv/webook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
	layout                = "2006-01-02"
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) Edit(ctx context.Context, u domain.User) error {
	parseBirthday, err := DateStringToUnixMillis(u.Birthday)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	err = r.dao.UpdateById(ctx, dao.User{
		Id:       u.Id,
		Birthday: parseBirthday,
		Profile:  u.Profile,
		Nickname: u.NickName,
	})
	return err
}

func (r *UserRepository) Profile(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.dao.GetUserDetail(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	tmp := u.Birthday
	var birthdayMillis string
	if tmp != 0 {
		birthdayMillis, err = UnixMillisToDateString(tmp)
	}
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Birthday: birthdayMillis,
		NickName: u.Nickname,
		Profile:  u.Profile,
	}, nil
}

func DateStringToUnixMillis(dateStr string) (int64, error) {
	layout := "2006-01-02"
	parsedTime, err := time.Parse(layout, dateStr)
	if err != nil {
		return 0, err
	}

	unixMillis := parsedTime.UnixNano() / int64(time.Millisecond)
	return unixMillis, nil
}

func UnixMillisToDateString(unixMillis int64) (string, error) {
	unixSeconds := unixMillis / 1000
	t := time.Unix(unixSeconds, 0)
	return t.Format("2006-01-02"), nil
}
