package repository

import (
	"context"
	"database/sql"
	"github.com/keweiLv/webook/internal/domain"
	"github.com/keweiLv/webook/internal/repository/cache"
	"github.com/keweiLv/webook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
	layout                = "2006-01-02"
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
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

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}
	// 数据不存在
	//if err == cache.ErrKeyNotExist {
	//}

	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	tmp := ue.Birthday
	var _ string
	if tmp != 0 {
		_, err = UnixMillisToDateString(tmp)
	}
	if err != nil {
		return domain.User{}, err
	}
	r.entityToDomain(ue)
	err = r.cache.Set(ctx, u)
	// 这里的 err 可以日志记录
	return u, err
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

func (r *UserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}

func (r *UserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.Password,
		Ctime:    u.Ctime.UnixMilli(),
	}
}
