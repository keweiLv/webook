package service

import (
	"context"
	"fmt"
	"github.com/keweiLv/webook/internal/repository"
	"github.com/keweiLv/webook/internal/service/sms"
	"math/rand"
)

const codeTplID = "1906254"

var (
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrCodeSendTooMant        = repository.ErrCodeSendTooMany
)

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type OldCodeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &OldCodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (svc *OldCodeService) Send(ctx context.Context, biz string, phone string) error {
	// 生产验证码
	code := svc.generateCode()
	// redis 存储
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = svc.smsSvc.Send(ctx, codeTplID, []string{code}, phone)
	//if err != nil {
	return err
}

func (svc *OldCodeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *OldCodeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
