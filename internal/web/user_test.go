package web

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/keweiLv/webook/internal/domain"
	"github.com/keweiLv/webook/internal/service"
	svcmocks "github.com/keweiLv/webook/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/genproto/googleapis/cloud/visionai/v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncrypt(t *testing.T) {
	password := "hello#123"
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(encrypted, []byte("hello#123"))
	assert.NoError(t, err)
}

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		visionai.AppPlatformEventBody
		name string

		mock func(ctrl *gomock.Controller) service.UserService

		reqBody string

		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "mock@mock.com",
					Password: "Hello@123",
				}).Return(nil)
				return usersvc
			},
			reqBody: `
{
	"email": "mock@mock.com",
	"password": "Hello@123",
	"confirmPassword": "Hello@123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "参数错误，bind失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `
{
	"email": "mock@mock.com"
	"password": "Hello@123"
}
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `
{
	"email": "mockoccom",
	"password": "Hello@123",
	"confirmPassword": "Hello@123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "你的邮箱格式不对",
		},
		{
			name: "两次输入的密码不一致",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `
{
	"email": "mock@mock.com",
	"password": "Hello@123",
	"confirmPassword": "Hello@1234"
}
`,
			wantCode: http.StatusOK,
			wantBody: "两次输入的密码不一致",
		},
		{
			name: "密码必须大于8位，包含数字、特殊字符",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `
{
	"email": "mock@mock.com",
	"password": "Hel111",
	"confirmPassword": "Hel111"
}
`,
			wantCode: http.StatusOK,
			wantBody: "密码必须大于8位，包含数字、特殊字符",
		},
		{
			name: "该邮箱已注册",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "mock@mock.com",
					Password: "Hello@123",
				}).Return(service.ErrUserDuplicateEmail)
				return usersvc
			},
			reqBody: `
{
	"email": "mock@mock.com",
	"password": "Hello@123",
	"confirmPassword": "Hello@123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "该邮箱已注册",
		},
		{
			name: "数据执行系统异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "mock@mock.com",
					Password: "Hello@123",
				}).Return(errors.New("任意的数据库异常"))
				return usersvc
			},
			reqBody: `
{
	"email": "mock@mock.com",
	"password": "Hello@123",
	"confirmPassword": "Hello@123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "系统异常",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			// 用不上 codesvc
			h := NewUserHandler(tc.mock(ctrl), nil)
			h.RegisterRoutes(server)

			req, err := http.NewRequest(http.MethodPost,
				"/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			t.Log(resp)

			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}
