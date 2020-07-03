package main

import (
	"net/http"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	setupLogger()
	startServer()
}

func startServer() {
	router := NewRouter()
	server := createServer(router)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		log.Fatal(server.ListenAndServe())
		wg.Done()
	}()

	log.Info("Server started")

	wg.Wait()

	log.Info("Server stopped")
}

func createServer(router http.Handler) http.Server {
	return http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}

func setupLogger() {
	// Set debug level
	log.SetLevel(log.DebugLevel)

	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: false,
		TimestampFormat:  time.Stamp,
		FullTimestamp:    true,
		ForceColors:      true,
		DisableColors:    false,
	})
}
