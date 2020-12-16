package proxy

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/stone-co/webhook-consumer/pkg/gateways/http/notifications"
)

func (n ProxyNotifier) Send(ctx context.Context, eventTypeHeader, eventIDHeader, body string) error {
	log := n.log.WithField("notifier", "proxy")

	req, err := http.NewRequest(http.MethodPost, n.serviceURL.String(), strings.NewReader(body))
	if err != nil {
		log.WithError(err).Info("unable to create a request")
		return fmt.Errorf("unable to create a request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(notifications.EventIDHeader, eventIDHeader)
	req.Header.Set(notifications.EventTypeHeader, eventTypeHeader)

	client := &http.Client{
		Timeout: n.timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).Info("unable to send request to service")
		return fmt.Errorf("unable to send request to service: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		log.Infof("unexpected status code when send request to service: %d", resp.StatusCode)
		return fmt.Errorf("unexpected status code when send request to service: %d", resp.StatusCode)
	}

	return nil
}
