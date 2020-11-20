package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"github.com/stone-co/webhook-consumer/pkg/common/configuration"
	"github.com/stone-co/webhook-consumer/pkg/domain/usecase"
	"github.com/stone-co/webhook-consumer/pkg/gateways/http"
	"github.com/stone-co/webhook-consumer/pkg/gateways/notifiers/stdout"
)

func main() {
	log := logrus.New()
	log.Infoln("starting webhook-consumer service...")

	cfg, err := configuration.LoadConfig()
	if err != nil {
		log.WithError(err).Fatal("unable to load app configuration")
	}

	// Stdout method only in this PR.
	method := stdout.New(log)

	usecase := usecase.NewNotificationUsecase(log, method)

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// NewServer HTTP Server listening for requests.
	httpServer := http.NewHttpServer(*cfg, log, usecase)
	go func() {
		log.Infof("starting http api at %s", httpServer.Addr)
		serverErrors <- httpServer.ListenAndServe()
	}()

	// =================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.WithError(err).Fatal("http server error")

	case sig := <-shutdown:
		log.Infof("server shutdown %v\n", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.HTTPConfig.ShutdownTimeout)
		defer cancel()

		log.Infof("stopping http server %v\n", sig)

		// Asking listener to shutdown and shed load.
		if err := httpServer.Shutdown(ctx); err != nil {
			_ = httpServer.Close()
			log.WithError(err).Error("could not stop server gracefully")
		}
		log.Infof("http server stopped %v\n", sig)
	}
}
