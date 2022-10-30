package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	gracePeriod = 5 * time.Second
)

func (app *application) serve() error {

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.AppPort),
		Handler:      app.routes(),
		ErrorLog:     log.New(app.logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	// shutdown mechanism.
	go func() {

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.logger.Printf("shutting down server on signal %s", s.String())

		// 5 seconds grace period before shutdown
		ctx, cancel := context.WithTimeout(context.Background(), gracePeriod)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.Printf("completing background tasks")

		app.wg.Wait() // wait for background go routines to finish before shutting down
		shutdownError <- nil
	}()

	app.logger.Printf("starting server(%s)", srv.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Printf("graceful shutdown of server(%s) completed", srv.Addr)

	return nil
}
