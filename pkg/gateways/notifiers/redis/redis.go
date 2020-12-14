package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/stone-co/webhook-consumer/pkg/domain"
)

var _ domain.Notifier = &RedisNotifier{}

type RedisNotifier struct {
	log  *logrus.Logger
	pool *redis.Pool
}

func New() *RedisNotifier {
	return &RedisNotifier{}
}
