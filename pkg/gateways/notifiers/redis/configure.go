package redis

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

func (n *RedisNotifier) Configure(log *logrus.Logger) error {
	var config Config
	prefix := ""
	if err := envconfig.Process(prefix, &config); err != nil {
		return err
	}

	n.log = log
	log.WithField("notifier", "redis").Infof("config:[%+v]", config)

	var err error
	n.pool, err = initPool(config)
	if err != nil {
		return err
	}

	return nil
}
