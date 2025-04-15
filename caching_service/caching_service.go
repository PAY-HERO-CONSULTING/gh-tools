package caching_service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/PAY-HERO-CONSULTING/gh-tools/apperrs"
	"github.com/PAY-HERO-CONSULTING/gh-tools/internal/caching"
	"github.com/gomodule/redigo/redis"
)

type (
	CachingService interface {
		CacheValue(key string, value interface{}) error
		Exists(key string) (bool, error)
		GetCachedValue(key string, result interface{}) error
		RemoveFromCache(key string) error
		CacheValueWithExpiry(key string, duration time.Duration, value interface{}) error
		GetPool() *redis.Pool
		IsErrNil(err error) bool
		Int64(key string) (int64, error)
	}

	cachingService struct {
		redisProvider caching.Redis
	}
)

func NewCachingService(
	redisURL string,
) CachingService {
	redisConfig := &caching.RedisConfig{
		IdleTimeout: 3 * time.Minute,
		MaxActive:   1000,
		MaxIdle:     10,
	}

	redisProvider := caching.NewRedisCache(redisURL, redisConfig)

	return &cachingService{
		redisProvider: redisProvider,
	}
}

func (s *cachingService) CacheValueWithExpiry(
	key string,
	duration time.Duration,
	value any,
) error {
	cacheData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = s.redisProvider.SetEx(key, duration, cacheData)
	if err != nil {
		return err
	}

	return nil
}

func (s *cachingService) CacheValue(
	key string,
	value any,
) error {
	cacheData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = s.redisProvider.Set(key, cacheData)
	if err != nil {
		return err
	}

	return nil
}

func (s *cachingService) Exists(
	key string,
) (bool, error) {
	exists, err := s.redisProvider.Exists(key)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *cachingService) GetCachedValue(
	key string,
	result any,
) error {
	payload, err := s.redisProvider.Get(key)
	if err != nil {
		return err
	}

	data, ok := payload.([]byte)
	if !ok {
		return apperrs.NewError(
			errors.New("cannot convert bytes to data"),
		)
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return apperrs.NewError(
			err,
		).LogErrorMessage("failed to unmarshal bytes to struct err: [%v]", err)
	}

	return nil
}

func (s *cachingService) RemoveFromCache(
	key string,
) error {
	err := s.redisProvider.Del(key)
	if err != nil {
		return err
	}

	return nil
}

func (s *cachingService) GetPool() *redis.Pool {
	return s.redisProvider.RedisPool()
}

func (s *cachingService) Int64(key string) (int64, error) {
	return s.redisProvider.Int64(key)
}

func (s *cachingService) IsErrNil(err error) bool {
	return s.redisProvider.IsErrNil(err)
}
