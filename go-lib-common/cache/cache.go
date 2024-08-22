//go:generate mockery --name=Cacher
package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/pkg/errors"
)

type Key string

type Data struct {
	Key   Key `json:"key"`
	Value any `json:"value"`
}

type Cacher interface {
	GetRedisInstance() *redis.Client
	Set(ctx context.Context, data Data, duration time.Duration) error
	SetNx(ctx context.Context, data Data, duration time.Duration) (isSuccessSet bool, err error)
	Get(ctx context.Context, key Key, dest any) error
	Delete(ctx context.Context, key Key) error
	BatchSet(ctx context.Context, datas []Data, duration time.Duration) error
	BatchGet(ctx context.Context, keys []Key, dest any) error
	Incr(ctx context.Context, key string) (*redis.IntCmd, error)
	IncrBy(ctx context.Context, key string, value int64) (*redis.IntCmd, error)
	Expire(ctx context.Context, key string, ttl time.Duration) (*redis.BoolCmd, error)
	Ttl(ctx context.Context, key string) (*redis.DurationCmd, error)
}

type Driver string

// Drivers
const (
	InMemoryDriver = Driver("inMemory")
	RedisDriver    = Driver("redis")
)

type Cache struct {
	driver   *Driver
	host     string
	password string
	database string
	username string
}

type Option func(*Cache)

func WithDriver(driver Driver) Option {
	return func(c *Cache) {
		c.driver = &driver
	}
}

func WithHost(host string) Option {
	return func(c *Cache) {
		c.host = host
	}
}

func WithPassword(password string) Option {
	return func(c *Cache) {
		c.password = password
	}
}

func WithUsername(username string) Option {
	return func(c *Cache) {
		c.username = username
	}
}

func WithDatabase(db string) Option {
	return func(c *Cache) {
		c.database = db
	}
}

var (
	ErrDriverUnavailable = errors.New("cache: driver unavailable")
)

func NewCache(
	options ...Option,
) (Cacher, error) {
	c := Cache{}
	for _, option := range options {
		option(&c)
	}

	if c.driver == nil {
		return nil, ErrDriverUnavailable
	}

	switch *c.driver {
	case RedisDriver:
		db, err := strconv.Atoi(c.database)
		if err != nil {
			return nil, err
		}
		return NewRedis(c.host, c.password, db, c.username), nil
	case InMemoryDriver:
		return NewInMemory(), nil
	default:
		return nil, ErrDriverUnavailable
	}
}
