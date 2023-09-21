package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/blubooks/blubook-cli/pkg/app"
	"github.com/blubooks/blubook-cli/pkg/router"
)

func Server() {

	var err error
	application := app.New()
	appRouter := router.New(application)

	address := fmt.Sprintf(":%d", 4080)
	log.Printf("Starting server %v", address)

	srv := &http.Server{
		Addr:    address,
		Handler: appRouter,
	}

	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		log.Println("Shutting down server")

		ctx, cancel := context.WithTimeout(context.Background(), 5000)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Println("server shutdown failure", err)
		}

		if err == nil {
			/*
				if err = dbConn.Close(); err != nil {
					logrus.WithField("error", err).Warn("Db connection closing failure")
				}
			*/
		}

		close(closed)
	}()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Println("server startup failure", err)

	}
}
