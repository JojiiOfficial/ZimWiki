package main

import (
	"net"
	"net/http"
	"strings"
	"time"

	"net/http/pprof"
	_ "net/http/pprof"

	"git.jojii.de/jojii/ZimWiki/handlers"
	"git.jojii.de/jojii/ZimWiki/zim"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Route defining a route
type Route struct {
	Name        string
	Methods     []HTTPMethod
	Pattern     string
	HandlerFunc RouteFunction
}

func (r Route) getMethods() []string {
	methods := make([]string, len(r.Methods))
	for i := range r.Methods {
		methods[i] = string(r.Methods[i])
	}
	return methods
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
type RouteFunction func(http.ResponseWriter, *http.Request, handlers.HandlerData) error

// Routes
var (
	globalRoutes = Routes{
		// -- Index routes
		// Main/Home pages and aliases
		{
			Name:        "IndexRoot",
			Pattern:     "/",
			Methods:     []HTTPMethod{GetMethod},
			HandlerFunc: handlers.Index,
		},
		{
			Name:        "IndexHtml",
			Pattern:     "/index.html",
			Methods:     []HTTPMethod{GetMethod},
			HandlerFunc: handlers.Index,
		},

		// -- Assets
		// Requests for static files
		{
			Name:        "",
			Pattern:     "/assets/{type}/{file}",
			Methods:     []HTTPMethod{GetMethod},
			HandlerFunc: handlers.Assets,
		},

		{
			Name:        "Search",
			Pattern:     "/search/{wiki}/",
			Methods:     []HTTPMethod{GetMethod, POSTMethod},
			HandlerFunc: handlers.Search,
		},
	}
)

// WikiRoutes
var (
	// Raw wiki page
	wikiRaw = Route{
		Name:        "",
		Methods:     []HTTPMethod{GetMethod},
		HandlerFunc: handlers.WikiRaw,
	}

	// Raw wiki page
	wikiView = Route{
		Name:        "WikiView",
		Methods:     []HTTPMethod{GetMethod},
		HandlerFunc: handlers.WikiView,
	}
)

// NewRouter create new router and its required components
func NewRouter(zimService *zim.Handler) *mux.Router {
	hd := handlers.HandlerData{
		ZimService: zimService,
	}

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range globalRoutes {
		router.
			Methods(route.getMethods()...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(RouteHandler(route.HandlerFunc, route.Name, hd))
	}

	// Add view handler
	router.Methods(wikiView.getMethods()...).
		PathPrefix("/wiki/view/").
		Name(wikiView.Name).
		Handler(RouteHandler(wikiView.HandlerFunc, wikiView.Name, hd))

	// Add raw handler
	router.Methods(wikiRaw.getMethods()...).
		PathPrefix("/wiki/raw/").
		Name(wikiRaw.Name).
		Handler(RouteHandler(wikiRaw.HandlerFunc, wikiRaw.Name, hd))

	attachProfiler(router)

	return router
}

func attachProfiler(router *mux.Router) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	// Manually add support for paths linked to by index page at /debug/pprof/
	router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/debug/pprof/block", pprof.Handler("block"))
	router.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
}

// RouteHandler logs stuff
func RouteHandler(inner RouteFunction, name string, hd handlers.HandlerData) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Error(err)
			}
		}()

		// Only debug routes which have a names
		needDebug := len(name) > 0
		if needDebug {
			log.Infof("[%s] %s\n", r.Method, name)
		}

		start := time.Now()

		// Set accept gzip
		hd.AcceptGzip = strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

		// Use ResponseProxy
		// This way we can automatically
		// send content compressed using gzip
		rp := handlers.NewResponseProxy(hd.AcceptGzip, w)
		defer rp.Done()

		// Process request and handle its error
		if err := inner(rp, r, hd); err != nil {
			if _, ok := err.(*net.OpError); ok {
				log.Warn(err)
				return
			}

			if err != handlers.ErrNotFound {
				sendServerError(rp)
			}

			log.Error(err)
			return
		}

		// Print duration of processing
		if needDebug {
			printProcessingDuration(start)
		}
	})
}
