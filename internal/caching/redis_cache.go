package caching

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type (
	Redis interface {
		Del(key string) error
		Exists(key string) (bool, error)
		Expire(key string, ttl time.Duration) (any, error)
		Float64(key string) (float64, error)
		Get(key string) (any, error)
		Int64(key string) (int64, error)
		Set(key string, val any) (any, error)
		SetEx(key string, ttl time.Duration, val any) (any, error)
		String(key string) (string, error)
		RedisPool() *redis.Pool
		IsErrNil(err error) bool
	}

	RedisConfig struct {
		IdleTimeout time.Duration
		MaxActive   int
		MaxIdle     int
	}

	redisCache struct {
		pool *redis.Pool
	}
)

func NewRedisCache(
	redisURL string,
	config *RedisConfig,
) Redis {
	return NewRedisCacheWithParams(
		redisURL,
		config,
	)
}

func NewRedisCacheWithParams(
	redisURL string,
	config *RedisConfig,
) *redisCache {

	redisPool := &redis.Pool{
		IdleTimeout: config.IdleTimeout,
		MaxActive:   config.MaxActive,
		MaxIdle:     config.MaxIdle,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(redisURL)
		},
	}

	return &redisCache{
		pool: redisPool,
	}
}

func (p *redisCache) Del(key string) error {
	_, err := p.do("DEL", key)
	return err
}

func (p *redisCache) Exists(key string) (bool, error) {
	return redis.Bool(p.do("EXISTS", key))
}

func (p *redisCache) Expire(key string, ttl time.Duration) (any, error) {
	return p.do("EXPIRE", key, int(ttl.Seconds()))
}

func (p *redisCache) Float64(
	key string,
) (float64, error) {
	return redis.Float64(p.Get(key))
}

func (p *redisCache) Get(
	key string,
) (any, error) {
	return p.do("GET", key)
}

func (p *redisCache) Int64(
	key string,
) (int64, error) {
	return redis.Int64(p.Get(key))
}

func (p *redisCache) RedisPool() *redis.Pool {
	return p.pool
}

func (p *redisCache) Set(
	key string,
	val any,
) (any, error) {
	return p.do("SET", key, val)
}

func (p *redisCache) SetEx(
	key string,
	ttl time.Duration,
	val any,
) (any, error) {
	return p.do("SETEX", key, int(ttl.Seconds()), val)
}

func (p *redisCache) String(
	key string,
) (string, error) {
	return redis.String(p.Get(key))
}

func (p *redisCache) do(
	commandName string,
	args ...any,
) (any, error) {
	conn := p.pool.Get()
	defer conn.Close()

	return conn.Do(commandName, args...)
}

func (p *redisCache) IsErrNil(
	err error,
) bool {
	return err == redis.ErrNil
}
