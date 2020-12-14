package proxy

import (
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/webhook-consumer/pkg/domain"
)

var _ domain.Notifier = &ProxyNotifier{}

type ProxyNotifier struct {
	log        *logrus.Logger
	serviceURL *url.URL
	timeout    time.Duration
}

func New() *ProxyNotifier {
	return &ProxyNotifier{}
}
