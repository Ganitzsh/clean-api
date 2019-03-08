package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
)

var config = NewAPIConfig()

func Routes() http.Handler {
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   config.Cors.AllowedOrigins,
		AllowedMethods:   config.Cors.AllowedMethods,
		AllowedHeaders:   config.Cors.AllowedHeaders,
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)
	r.Use(middleware.AllowContentType(strings.Split(APIV1ContentTypes, ",")...))
	if !config.DevMode {
		r.Use(middleware.Logger)
	}
	return r
}

func Start() error {
	srv := http.Server{
		Addr:    config.GetFullHost(),
		Handler: Routes(),
	}
	done := make(chan bool)
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigint
		if err := srv.Shutdown(context.Background()); err != nil {
			logrus.Fatalf("Could not shutdown: %v", err)
		}
		close(done)
	}()
	var err error
	if config.GetTLS() {
		err = srv.ListenAndServeTLS(config.TLSCert, config.TLSKey)
	} else {
		err = srv.ListenAndServe()
	}
	if err != http.ErrServerClosed {
		return fmt.Errorf("Error: %s", err)
	}
	logrus.Info("Shutting down...")
	<-done
	return nil
}
