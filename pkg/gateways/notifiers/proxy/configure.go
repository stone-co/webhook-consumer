package proxy

import (
	"fmt"
	"net/url"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Url     string        `envconfig:"PROXY_NOTIFIER_URL"`
	Timeout time.Duration `envconfig:"PROXY_NOTIFIER_TIMEOUT" default:"10s"`
}

func (n *ProxyNotifier) Configure(log *logrus.Logger) error {
	var config Config
	prefix := ""
	if err := envconfig.Process(prefix, &config); err != nil {
		return err
	}

	n.log = log
	log.WithField("notifier", "proxy").Infof("config:[%+v]", config)

	var err error

	n.timeout = config.Timeout
	n.serviceURL, err = url.Parse(config.Url)
	if err != nil || config.Url == "" {
		return fmt.Errorf("failed to parse url '%s': %v", config.Url, err)
	}
	n.log = log

	n.log.WithField("notifier", "proxy").Infof("url:[%s] timeout:[%s]", config.Url, n.timeout.String())

	return nil

}
