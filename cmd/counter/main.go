package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"counter/internal/app/statistics"
	"counter/pkg/counter"
)

const (
	storageFileName = "test" // TODO: move to .env
	counterLifeTime = time.Minute
)

func main() {
	storage := counter.NewFileStorage(storageFileName)
	counterProvider := counter.New(storage, counterLifeTime)
	if err := counterProvider.Load(); err != nil {
		fmt.Printf("can not load counter data")
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	go func() {
		select {
		case sig := <-c:
			fmt.Printf("Got %s signal. Aborting...\n", sig)
			if err := counterProvider.Cancel(); err != nil {
				log.Fatal(err)
			}
			os.Exit(1)
		}
	}()

	statHandler := statistics.New(counterProvider)
	router := chi.NewRouter().With(counterProvider.Middleware)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Get("/counts", statHandler.Get)

	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}
