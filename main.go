package main

import (
	"github.com/havlan/searchproxy/cmd"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {

	// if we can't access api key, don't bother starting up
	if os.Getenv("ocp_apim_subscription_key") == "" {
		log.Fatalln("ApiKey not found")
	}

	metricsMiddleware := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	router := mux.NewRouter()
	router.Use(std.HandlerProvider("MetricsMiddleware", metricsMiddleware))
	router.Use(loggingMiddleware)

	// api/
	router.HandleFunc("/api/beta", cmd.HandleBing).Methods("GET")

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
	}

	// Serve our handler.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Panicf("error while serving: %s", err)
		}
	}()

	// Serve our metrics.
	go func() {
		log.Printf("metrics listening at %s", "127.0.0.1:8001")
		if err := http.ListenAndServe("127.0.0.1:8001", promhttp.Handler()); err != nil {
			log.Panicf("error while serving metrics: %s", err)
		}
	}()
	
	// Wait until some signal is captured.
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)
	<-sigC
}
