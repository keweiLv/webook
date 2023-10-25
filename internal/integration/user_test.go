package integration

import (
	"bytes"
	"encoding/json"
	"github.com/keweiLv/webook/internal/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_e2e_SendLoginSMSCode(t *testing.T) {
	server := InitWebServer()
	//rdb := ioc.InitRedis()
	testCases := []struct {
		name string

		// 准备数据
		before func(t *testing.T)
		// 验证数据
		after   func(t *testing.T)
		reqBody string

		wantCode int
		wantBody web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {

			},
			reqBody: `{
			"phone":15212345678					
			}`,
			wantCode: 200,
			wantBody: web.Result{
				Msg: "发送成功",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodPost,
				"/users/login_sms/code/send", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()

			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			var webRes web.Result
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			assert.Equal(t, tc.wantBody, webRes)
		})
	}
}
