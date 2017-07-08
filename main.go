// ReShifter enables backing up and restoring OpenShift clusters.
// The ReShifter app launches an API and a UI at port 8080.
// The API is instrumented, exposing Prometheus metrics.
// When launching the app with the defaults, the backups are created in the
// current directory and the temporary work files are placed in the /tmp directory.
package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/mhausenblas/reshifter/pkg/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	if envd := os.Getenv("DEBUG"); envd != "" {
		log.SetLevel(log.DebugLevel)
	}
	port := "8080"
	host, _ := util.ExternalIP()
	r := mux.NewRouter()
	// the HTTP API:
	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/v1/version", versionHandler).Methods("GET")
	r.HandleFunc("/v1/explorer", explorerHandler).Methods("GET")
	r.HandleFunc("/v1/epstats", epstatsHandler).Methods("GET")
	r.HandleFunc("/v1/backup", backupCreateHandler).Methods("POST")
	r.HandleFunc("/v1/backup/all", backupListHandler).Methods("GET")
	r.HandleFunc("/v1/backup/{backupid:[0-9]+}", backupRetrieveHandler).Methods("GET")
	r.HandleFunc("/v1/restore", restoreHandler).Methods("POST")
	log.Printf("Serving API from: %s:%s/v1", host, port)
	// the Web UI:
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/")))
	log.Printf("Serving UI from: %s:%s/", host, port)
	http.Handle("/", r)
	// the app server:
	srv := &http.Server{Handler: r, Addr: "0.0.0.0:" + port}
	log.Fatal(srv.ListenAndServe())
}
