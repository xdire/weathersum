package service

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var logger = zerolog.New(os.Stdout).With().Timestamp().Logger()

func StartWeatherService(routes *mux.Router, port int) error {
	srv := &http.Server{
		Handler:      logRequest(routes),
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  time.Second * 60,
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	go func() {
		logger.Info().Msg("Starting HTTP Weather Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
			c <- os.Interrupt
		}
	}()

	// Await for system signal
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	// Shutdown with the deadline
	err := srv.Shutdown(ctx)
	if err != nil {
		return err
	}

	logger.Info().Msg("weather http service shutting down")
	return nil
}

func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info().Str("method", r.Method).Str("path", r.URL.Path).Msg("Request received")

		next.ServeHTTP(w, r)

		logger.Info().Str("method", r.Method).Str("path", r.URL.Path).Msg("Response sent")
	})
}
