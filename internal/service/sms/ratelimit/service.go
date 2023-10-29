package ratelimit

import (
	"context"
	"fmt"
	"github.com/keweiLv/webook/internal/service/sms"
	"github.com/keweiLv/webook/pkg/ratelimit"
)

var limitedErr = fmt.Errorf("触发限流")

type RatelimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSMSService{
		svc:     svc,
		limiter: limiter,
	}

}

func (s *RatelimitSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		return fmt.Errorf("短信服务限流出现问题,%w", err)
	}
	if limited {
		return limitedErr
	}
	err = s.svc.Send(ctx, tpl, args, numbers...)
	return err
}
