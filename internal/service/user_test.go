package service

import (
	"context"
	"errors"
	"github.com/keweiLv/webook/internal/domain"
	"github.com/keweiLv/webook/internal/repository"
	repomocks "github.com/keweiLv/webook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func TestOldUserService_Login(t *testing.T) {
	// 公共时间
	now := time.Now()

	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository

		ctx      context.Context
		email    string
		password string

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "vip@qq.com").
					Return(domain.User{
						Email:    "vip@qq.com",
						Password: "$2a$10$ha/KhE0itj.xEOiJ0JER4eNm6atJScAkUGoH8JBWYoBdIw0BNmrwK",
						Phone:    "18112345678",
						Ctime:    now,
					}, nil)
				return repo
			},
			email:    "vip@qq.com",
			password: "hello@123",

			wantUser: domain.User{
				Email:    "vip@qq.com",
				Password: "$2a$10$ha/KhE0itj.xEOiJ0JER4eNm6atJScAkUGoH8JBWYoBdIw0BNmrwK",
				Phone:    "18112345678",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "vip@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "vip@qq.com",
			password: "hello@123",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "DB错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "vip@qq.com").
					Return(domain.User{}, errors.New("mock db error"))
				return repo
			},
			email:    "vip@qq.com",
			password: "hello@123",

			wantUser: domain.User{},
			wantErr:  errors.New("mock db error"),
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "vip@qq.com").
					Return(domain.User{
						Email:    "vip@qq.com",
						Password: "$2a$10$ha/KhE0itj.xEOiJ0JER4eNm6atJScAkUGoH8JBWYoBdIw0BNmrwK",
						Phone:    "18112345678",
						Ctime:    now,
					}, nil)
				return repo
			},
			email:    "vip@qq.com",
			password: "111hello@123",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// 具体测试代码
			svc := NewUserService(tc.mock(ctrl))
			u, err := svc.Login(tc.ctx, tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

func TestEncrypted(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("hello@123"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}
