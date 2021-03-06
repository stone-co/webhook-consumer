package main

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/webhook-consumer/pkg/domain"
	"github.com/stone-co/webhook-consumer/pkg/gateways/notifiers/proxy"
	"github.com/stone-co/webhook-consumer/pkg/gateways/notifiers/redis"
	"github.com/stone-co/webhook-consumer/pkg/gateways/notifiers/stdout"
)

// Stdout notifier is used only to debug purpose.
var notificationTypes = map[string]domain.Notifier{
	"stdout": stdout.New(),
	"proxy":  proxy.New(),
	"redis":  redis.New(),
}

func defineNotifiers(notifierList string, log *logrus.Logger) ([]domain.Notifier, error) {
	notifiersToConfig, err := extractNotifiersFromConfig(notifierList)
	if err != nil {
		return nil, fmt.Errorf("configure failed when loading notifiers: %v", err)
	}

	result := []domain.Notifier{}
	for _, notifier := range notifiersToConfig {
		impl := notificationTypes[notifier]
		if err := impl.Configure(log); err != nil {
			return nil, fmt.Errorf("configure failed in [%s] notifier: %v", notifier, err)
		}

		result = append(result, impl)
	}

	return result, nil
}

func extractNotifiersFromConfig(notifiers string) ([]string, error) {
	usedNotifications := map[string]bool{}

	result := []string{}
	for _, notifier := range strings.Split(notifiers, ";") {
		notifier = strings.TrimSpace(notifier)
		if notifier == "" {
			continue
		}

		notifier = strings.ToLower(notifier)
		_, ok := notificationTypes[notifier]
		if !ok {
			return nil, fmt.Errorf("undefined notifier: %v", notifier)
		}

		_, used := usedNotifications[notifier]
		if used {
			return nil, fmt.Errorf("duplicated notifier: %v", notifier)
		}

		usedNotifications[notifier] = true
		result = append(result, notifier)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("empty notifier list")
	}

	return result, nil
}
