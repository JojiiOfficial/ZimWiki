package main

import (
	"net/http"
	"time"

	"git.jojii.de/jojii/zimserver/handlers"
	"git.jojii.de/jojii/zimserver/zim"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Route defining a route
type Route struct {
	Name        string
	Method      HTTPMethod
	Pattern     string
	HandlerFunc RouteFunction
}

// HTTPMethod http method. GET, POST, DELETE, HEADER, etc...
type HTTPMethod string

// HTTP methods
const (
	GetMethod    HTTPMethod = "GET"
	POSTMethod   HTTPMethod = "POST"
	PUTMethod    HTTPMethod = "PUT"
	DeleteMethod HTTPMethod = "DELETE"
)

// Routes all HTTP routes
type Routes []Route

// RouteFunction function for handling a route
type RouteFunction func(http.ResponseWriter, *http.Request, *handlers.HandlerData) error

// Routes
var (
	routes = Routes{
		// -- Index routes
		// Main/Home pages and aliases
		{
			Name:        "IndexRoot",
			Pattern:     "/",
			Method:      GetMethod,
			HandlerFunc: handlers.Index,
		},
		{
			Name:        "IndexHtml",
			Pattern:     "/index.html",
			Method:      GetMethod,
			HandlerFunc: handlers.Index,
		},

		// -- Assets
		// Requests for static files
		{
			Name:        "",
			Pattern:     "/assets/{type}/{file}",
			Method:      GetMethod,
			HandlerFunc: handlers.Assets,
		},

		// -- Wiki Raw
		// Raw wiki page
		{
			Name:        "WikiRaw",
			Pattern:     "/wiki/raw/{wikiID}/{namespace}/{file}",
			Method:      GetMethod,
			HandlerFunc: handlers.WikiRaw,
		},
	}
)

// NewRouter create new router and its required components
func NewRouter(zimService *zim.Handler) *mux.Router {
	hd := handlers.HandlerData{
		ZimService: zimService,
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(string(route.Method)).
			Path(route.Pattern).
			Name(route.Name).
			Handler(RouteHandler(route.HandlerFunc, route.Name, &hd))
	}

	return router
}

// RouteHandler logs stuff
func RouteHandler(inner RouteFunction, name string, hd *handlers.HandlerData) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Info(err)
			}
		}()

		// Only debug routes which have a names
		needDebug := len(name) > 0
		if needDebug {
			log.Infof("[%s] %s\n", r.Method, name)
		}

		start := time.Now()

		// Process request and handle its error
		if err := inner(w, r, hd); err != nil {
			sendServerError(w)
			log.Error(err)
			return
		}

		// Print duration of processing
		if needDebug {
			printProcessingDuration(start)
		}
	})
}
