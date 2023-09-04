package cache

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type CodeLruCache interface {
	Set(biz, phone, code string) error
	Verify(biz, phone, inputCode string) (bool, error)
}

type Cache struct {
	data    map[string]cacheEntry
	mu      sync.Mutex
	cleanup chan string
}

type cacheEntry struct {
	value     string
	expiry    time.Time
	createdAt time.Time
}

func NewCache() *Cache {
	cache := &Cache{
		data:    make(map[string]cacheEntry),
		cleanup: make(chan string),
	}

	go cache.startCleanupWorker()

	return cache
}

func (c *Cache) Set(biz, phone, code, cnt string, isInside bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := c.generateKey(biz, phone)
	cntKey := c.generateCntKey(biz, phone)
	tmp := c.data[key]
	if time.Now().Before(tmp.expiry) && !isInside {
		return ErrCodeSendTooMany
	}
	c.data[key] = cacheEntry{
		value:     code,
		expiry:    time.Now().Add(5 * time.Minute),
		createdAt: time.Now(),
	}
	c.data[cntKey] = cacheEntry{
		value:     cnt,
		expiry:    time.Now().Add(5 * time.Minute),
		createdAt: time.Now(),
	}
	return nil
}

func (c *Cache) Verify(biz, phone, code string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	key := c.generateKey(biz, phone)
	cntKey := c.generateCntKey(biz, phone)
	entry, ok := c.data[key]
	cntEntry, ok := c.data[cntKey]
	if !ok {
		return false, ErrUnknownForCode
	}
	check, _ := strconv.Atoi(cntEntry.value)
	if check <= 0 {
		return false, ErrCodeVerifyTooManyTimes
	}

	if time.Now().Before(entry.expiry) {
		if entry.value == code {
			delete(c.data, key)
			delete(c.data, cntKey)
			return true, nil
		} else {
			newCnt := strconv.Itoa(check - 1)
			c.Set(biz, phone, newCnt, entry.value, true)
			c.cleanup <- key
			return false, nil
		}
	}
	c.cleanup <- key
	return false, ErrUnknownForCode
}

func (c *Cache) startCleanupWorker() {
	for key := range c.cleanup {
		c.mu.Lock()
		entry, ok := c.data[key]
		if !ok {
			c.mu.Unlock()
			continue
		}

		if time.Since(entry.createdAt) >= 5*time.Minute {
			fmt.Printf("Cleaning up expired entry: Phone %s\n", key)
			delete(c.data, key)
		} else {
			remainingTime := 5*time.Minute - time.Since(entry.createdAt)
			c.mu.Unlock()
			time.Sleep(remainingTime)
			c.cleanup <- key
		}

		c.mu.Unlock()
	}
}

func (c *Cache) generateKey(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (c *Cache) generateCntKey(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s:%s", biz, phone, "cnt")
}
