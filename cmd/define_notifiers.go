package main

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/webhook-consumer/pkg/common/configuration"
	"github.com/stone-co/webhook-consumer/pkg/domain"
	"github.com/stone-co/webhook-consumer/pkg/gateways/notifiers/proxy"
	"github.com/stone-co/webhook-consumer/pkg/gateways/notifiers/stdout"
)

// Stdout notifier is used only to debug purpose.
var notificationTypes = map[string]domain.Notifier{
	"stdout": stdout.New(),
	"proxy":  proxy.New(),
}

func defineNotifiers(cfg *configuration.Config, log *logrus.Logger) ([]domain.Notifier, error) {

	notifiersToConfig, err := extractNotifiersFromConfig(cfg.NotifierList)
	if err != nil {
		return nil, fmt.Errorf("configure failed when loading notifiers: %v", err)
	}

	result := []domain.Notifier{}

	for _, notifier := range notifiersToConfig {
		impl, _ := notificationTypes[notifier]
		if err := impl.Configure(cfg, log); err != nil {
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
