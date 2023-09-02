package ioc

import (
	"github.com/keweiLv/webook/internal/service/sms"
	"github.com/keweiLv/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	return memory.NewService()
}
