package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/webhook-consumer/pkg/common/configuration"
	"github.com/stone-co/webhook-consumer/pkg/domain"
	"github.com/stone-co/webhook-consumer/pkg/gateways/http/notifications"
)

var _ domain.Notifier = &ProxyNotifier{}

type ProxyNotifier struct {
	log        *logrus.Logger
	serviceURL *url.URL
	timeout    time.Duration
}

func New(log *logrus.Logger) *ProxyNotifier {
	return &ProxyNotifier{
		log: log,
	}
}

func (n *ProxyNotifier) Configure(config *configuration.Config) error {
	var err error

	n.timeout = config.ProxyNotifier.Timeout
	n.serviceURL, err = url.Parse(config.ProxyNotifier.Url)
	if err != nil || config.ProxyNotifier.Url == "" {
		return fmt.Errorf("failed to parse url '%s': %v", config.ProxyNotifier.Url, err)
	}

	return nil
}

func (n ProxyNotifier) Send(ctx context.Context, input domain.NotificationInput) error {
	log := n.log.WithField("notifier", "proxy")

	req, err := http.NewRequest(http.MethodPost, n.serviceURL.String(), strings.NewReader(input.Body))
	if err != nil {
		log.Infof("unable to create a request: %w", err)
		return fmt.Errorf("unable to create a request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(notifications.EventIDHeader, input.Header.EventID)
	req.Header.Set(notifications.EventTypeHeader, input.Header.EventType)

	client := &http.Client{
		Timeout: n.timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Infof("unable to send request to service: %v", err)
		return fmt.Errorf("unable to send request to service: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Infof("unexpected status code when send request to service: %d", resp.StatusCode)
		return fmt.Errorf("unexpected status code when send request to service: %d", resp.StatusCode)
	}

	return nil
}
