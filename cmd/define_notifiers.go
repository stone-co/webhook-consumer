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

func defineNotifiers(cfg *configuration.Config, log *logrus.Logger) ([]domain.Notifier, error) {

	type notifierInfo struct {
		notifier domain.Notifier
		used     bool
	}

	// Stdout notifier is used only to debug purpose.
	var notificationTypes = map[string]notifierInfo{
		"stdout": {
			notifier: stdout.New(log),
		},
		"proxy": {
			notifier: proxy.New(log),
		},
	}

	result := []domain.Notifier{}

	for _, notifier := range strings.Split(cfg.NotifierList, ";") {
		notifier = strings.TrimSpace(notifier)
		if notifier == "" {
			continue
		}

		notifier = strings.ToLower(notifier)
		info, ok := notificationTypes[notifier]
		if !ok {
			return nil, fmt.Errorf("undefined notifier: %v", notifier)
		}

		if info.used {
			return nil, fmt.Errorf("duplicated notifier: %v", notifier)
		}

		if err := info.notifier.Configure(cfg); err != nil {
			return nil, fmt.Errorf("configure failed in [%s] notifier: %v", notifier, err)
		}

		info.used = true
		notificationTypes[notifier] = info

		result = append(result, info.notifier)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("empty notifier list")
	}

	return result, nil
}
