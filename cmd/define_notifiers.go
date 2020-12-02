package main

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/stone-co/webhook-consumer/pkg/domain"
	"github.com/stone-co/webhook-consumer/pkg/gateways/notifiers/stdout"
)

func defineNotifiers(log *logrus.Logger, notifiers string) ([]domain.NotifierMethod, error) {

	type notifierInfo struct {
		notifier domain.NotifierMethod
		used     bool
	}

	// Stdout method is used only to debug purpose.
	var notificationTypes = map[string]notifierInfo{
		"stdout": {
			notifier: stdout.New(log),
		},
		"proxyapi": {
			notifier: stdout.New(log),
		},
	}

	result := []domain.NotifierMethod{}

	for _, notifier := range strings.Split(notifiers, ";") {
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

		info.used = true
		notificationTypes[notifier] = info

		result = append(result, info.notifier)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("empty notifier list")
	}

	return result, nil
}
