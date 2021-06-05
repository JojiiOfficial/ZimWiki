package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JojiiOfficial/ZimWiki/zim"
	log "github.com/sirupsen/logrus"
	"github.com/pelletier/go-toml"
)

func main() {
	setupLogger()

	config, err := toml.LoadFile("config.toml")

	// If the configuration file does not exist, return an error
	if err != nil {
		log.Error(err)
	} 
	
	// Define the variables retrieved from the configration file
	configTree := config.Get("Config").(*toml.Tree)
	libPath := configTree.Get("LibraryPath").(string)
	port := configTree.Get("Port").(string)

	if len(os.Args) > 1 {
		libPath = os.Args[1]
	}

	// Verify library path
	s, err := os.Stat(libPath)
	if err != nil {
		log.Errorf("Can't use '%s' as library path. %s", libPath, err)
		return
	}
	if !s.IsDir() {
		log.Error("Library must be a path!")
		return
	}

	service := zim.New(libPath)
	err = service.Start(libPath)
	if err != nil {
		log.Fatalln(err)
		return
	}

	startServer(service, port)
}

func startServer(zimService *zim.Handler, port string) {
	router := NewRouter(zimService)
	server := createServer(router, port)

	// Start server
	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	log.Info("Server started")
	awaitExit(&server)
}

// Build a new Http server
func createServer(router http.Handler, port string) http.Server {
	return http.Server{
		Addr:    port,
		Handler: router,
	}
}

// Shutdown server gracefully
func awaitExit(httpServer *http.Server) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM)

	// await os signal
	<-signalChan

	// Create a deadline for the await
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// Remove that ugly '^C'
	fmt.Print("\r")

	log.Info("Shutting down server")

	if httpServer != nil {
		err := httpServer.Shutdown(ctx)
		if err != nil {
			log.Warn(err)
		}

		log.Info("HTTP server shutdown complete")
	}

	log.Info("Shutting down complete")
	os.Exit(0)
}

func setupLogger() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: false,
		TimestampFormat:  time.Stamp,
		FullTimestamp:    true,
		ForceColors:      true,
		DisableColors:    false,
	})
}
