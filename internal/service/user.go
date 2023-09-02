package service

import (
	"context"
	"errors"
	"github.com/keweiLv/webook/internal/domain"
	"github.com/keweiLv/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")
)

type UserService interface {
	Login(ctx context.Context, email, password string) (domain.User, error)
	SignUp(ctx context.Context, u domain.User) error
	Edit(ctx context.Context, u domain.User) error
	Profile(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}

type OldUserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &OldUserService{
		repo: repo,
	}
}

func (svc *OldUserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	// 判断用户是否存在
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *OldUserService) SignUp(ctx context.Context, u domain.User) error {
	// 密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *OldUserService) Edit(ctx context.Context, u domain.User) error {
	err := svc.repo.Edit(ctx, u)
	return err
}

func (svc *OldUserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}

func (svc *OldUserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, phone)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	// 明确知道是新用户
	// 此时的 u 没有 ID
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil {
		return u, err
	}
	// 可能有主从延迟问题
	return svc.repo.FindByPhone(ctx, phone)
}
