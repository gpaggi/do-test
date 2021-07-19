package server

import (
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/tg123/go-htpasswd"

	H "echoapi/handler"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "echoapi_http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})
	httpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "echoapi_http_requests_total",
		Help: "How many HTTP requests processed by status code, method and HTTP path.",
	}, []string{"code", "method", "path"},
	)
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

// Prometheus middleware to handle requests timeing and count
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)
		httpRequests.WithLabelValues(strconv.Itoa(rw.statusCode), r.Method, path).Inc()
		timer.ObserveDuration()
	})
}

// Start binds the webserver and returns an error on failure
func Start(listen, htpasswdPath string) error {
	// Configure the webserver
	r := mux.NewRouter()
	r.Use(prometheusMiddleware)
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	htp, err := htpasswd.New(htpasswdPath, htpasswd.DefaultSystems, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Successfully opened htpasswd file at %s", htpasswdPath)

	// Simple two routes setup
	r.Path("/metrics").Handler(promhttp.Handler()).Methods(http.MethodGet)
	r.HandleFunc("/api/echo", H.BasicAuth(H.EchoHandler, htp)).Methods(http.MethodPost, http.MethodPut)

	// Start the webserver
	log.Infof("Starting up and binding to %s", listen)
	return http.ListenAndServe(listen, loggedRouter)
}
