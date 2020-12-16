package redis

import (
	"fmt"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Config struct {
	Address            string        `envconfig:"REDIS_ADDR" required:"true"`
	Port               string        `envconfig:"REDIS_PORT" required:"true"`
	Password           string        `envconfig:"REDIS_PASSWORD"`
	UseTLS             bool          `envconfig:"REDIS_USE_TLS" default:"false"`
	MaxIdle            int           `envconfig:"REDIS_MAX_IDLE" default:"100"`
	MaxActive          int           `envconfig:"REDIS_MAX_ACTIVE" default:"1000"`
	IdleTimeout        time.Duration `envconfig:"REDIS_IDLE_TIMEOUT" default:"1m"`
	DialConnectTimeout time.Duration `envconfig:"REDIS_CONNECT_TIMEOUT" default:"1s"`
	DialReadTimeout    time.Duration `envconfig:"REDIS_READ_TIMEOUT" default:"300ms"`
	DialWriteTimeout   time.Duration `envconfig:"REDIS_WRITE_TIMEOUT" default:"300ms"`
}

func (c Config) Addr() string {
	return strings.Join([]string{c.Address, c.Port}, ":")
}

func (c Config) URL() string {
	scheme := "redis"
	if c.UseTLS {
		scheme = "rediss"
	}
	addr := c.Addr()
	if c.Password != "" {
		addr = fmt.Sprintf(":%s@%s", c.Password, addr)
	}
	return fmt.Sprintf("%s://%s", scheme, addr)
}

func initPool(cfg Config) (*redis.Pool, error) {
	redisPool := &redis.Pool{
		MaxIdle:     cfg.MaxIdle,
		MaxActive:   cfg.MaxActive,
		IdleTimeout: cfg.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", cfg.Addr(),
				redis.DialPassword(cfg.Password),
				redis.DialUseTLS(cfg.UseTLS),
				redis.DialConnectTimeout(cfg.DialConnectTimeout),
				redis.DialReadTimeout(cfg.DialReadTimeout),
				redis.DialWriteTimeout(cfg.DialWriteTimeout))
			if err != nil {
				return nil, fmt.Errorf("could not connect to redis: %w", err)
			}
			return conn, nil
		},
	}

	err := ping(redisPool)
	if err != nil {
		return nil, err
	}

	return redisPool, nil
}

func ping(pool *redis.Pool) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		return err
	}

	return nil
}
