package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/JojiiOfficial/ZimWiki/handlers"
	"github.com/JojiiOfficial/ZimWiki/zim"
	"github.com/pelletier/go-toml"
	log "github.com/sirupsen/logrus"
)

var (
	//go:embed html/*
	WebFS embed.FS

	//go:embed locale.zip
	LocaleByte []byte
)

type configStruct struct {
	libPath             string
	address             string
	EnableSearchCache   bool
	SearchCacheDuration int
}

func main() {
	setupLogger()

	handlers.WebFS = WebFS

	handlers.LocaleByte = LocaleByte

	// Default configuration of ZimWiki
	defaultConfig, _ := toml.Load(`
	[Config]
	LibraryPath = "./library"
	Address = ":8080"
	EnableSearchCache = "true"
	SearchCacheDuration = "2"`)

	// Load default configuration
	libPath := defaultConfig.Get("Config.LibraryPath").(string)
	address := defaultConfig.Get("Config.Address").(string)
	EnableSearchCache, _ := strconv.ParseBool(defaultConfig.Get("Config.EnableSearchCache").(string))
	SearchCacheDuration, _ := strconv.Atoi(defaultConfig.Get("Config.SearchCacheDuration").(string))

	// Load configuration file
	configData, err := toml.LoadFile("config.toml")

	// If the configuration file has been successfully loaded
	if err == nil {
		// Load the configuration from the configuration file
		configDataTree := configData.Get("Config").(*toml.Tree)
		libPath = configDataTree.Get("LibraryPath").(string)
		address = configDataTree.Get("Address").(string)
		EnableSearchCache, _ = strconv.ParseBool(configDataTree.Get("Config.EnableSearchCache").(string))
		SearchCacheDuration, _ = strconv.Atoi((configDataTree.Get("Config.SearchCacheDuration")).(string))
	} else {
		log.Error("Config.toml not found, default configuration will be used.")
	}

	config := configStruct{libPath: libPath, address: address, EnableSearchCache: EnableSearchCache, SearchCacheDuration: SearchCacheDuration}

	handlers.EnableSearchCache = EnableSearchCache
	handlers.SearchCacheDuration = SearchCacheDuration

	if len(os.Args) > 1 {
		config.libPath = os.Args[1]
	}

	// Verify library path
	s, err := os.Stat(config.libPath)
	if err != nil {
		log.Errorf("Can't use '%s' as library path. %s", config.libPath, err)
		return
	}
	if !s.IsDir() {
		log.Error("Library must be a path!")
		return
	}

	service := zim.New(config.libPath)
	err = service.Start(config.libPath)
	if err != nil {
		log.Fatalln(err)
		return
	}

	startServer(service, config)
}

func startServer(zimService *zim.Handler, config configStruct) {
	router := NewRouter(zimService)
	server := createServer(router, config)

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
func createServer(router http.Handler, config configStruct) http.Server {
	return http.Server{
		Addr:    config.address,
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
