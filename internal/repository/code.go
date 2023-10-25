package repository

import (
	"context"
	"github.com/keweiLv/webook/internal/repository/cache"
)

var (
	ErrCodeSendTooMany        = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
)

type CodeRepository interface {
	Store(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

type CachedCodeRepository struct {
	cache    cache.CodeRedisCache
	lruCache *cache.Cache
}

//func NewCodeRepository(c cache.CodeRedisCache, lru *cache.Cache) CodeRepository {
//	return &CachedCodeRepository{
//		cache:    c,
//		lruCache: lru,
//	}
//}

func NewCodeRepository(c cache.CodeRedisCache) CodeRepository {
	return &CachedCodeRepository{
		cache: c,
	}
}

func (repo *CachedCodeRepository) Store(ctx context.Context, biz string, phone string, code string) error {
	//return repo.cache.Set(ctx, biz, phone, code)
	// 使用 lru 本地缓存
	return repo.lruCache.Set(biz, phone, code, "3", false)
}

func (repo *CachedCodeRepository) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	//return repo.cache.Verify(ctx, biz, phone, inputCode)
	// 使用 lru 本地缓存
	return repo.lruCache.Verify(biz, phone, inputCode)
}
