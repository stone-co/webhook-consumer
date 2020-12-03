package http

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"

	"github.com/stone-co/webhook-consumer/pkg/common/configuration"
	"github.com/stone-co/webhook-consumer/pkg/common/keys"
	"github.com/stone-co/webhook-consumer/pkg/common/validator"
	"github.com/stone-co/webhook-consumer/pkg/domain"
	"github.com/stone-co/webhook-consumer/pkg/gateways/http/healthcheck"
	"github.com/stone-co/webhook-consumer/pkg/gateways/http/notifications"
)

func NewHttpServer(config configuration.Config, keys *keys.Config, log *logrus.Logger, usecase domain.NotificationUsecase) *http.Server {
	validator := validator.NewJSONValidator()

	notificationsHandler := notifications.NewHandler(log, validator, keys, usecase)

	api := NewApi(log, notificationsHandler)
	return api.NewServer("0.0.0.0", config.HTTPConfig)
}

type Api struct {
	log           *logrus.Logger
	Healthcheck   healthcheck.Handler
	notifications *notifications.Handler
}

func NewApi(log *logrus.Logger, notifications *notifications.Handler) *Api {
	return &Api{
		log:           log,
		notifications: notifications,
	}
}

func (a *Api) NewServer(host string, cfg configuration.HTTPConfig) *http.Server {
	// Router
	r := mux.NewRouter()

	// Handlers
	r.HandleFunc("/healthcheck", a.Healthcheck.Get).Methods(http.MethodGet)
	r.HandleFunc("/metrics", promhttp.Handler().ServeHTTP).Methods(http.MethodGet)
	r.HandleFunc("/api/v0/notifications", a.notifications.New).Methods(http.MethodPost)

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())

	n.UseHandler(r)

	endpoint := fmt.Sprintf("%s:%d", host, cfg.Port)

	srv := &http.Server{
		Handler: n,
		Addr:    endpoint,
	}

	return srv
}
