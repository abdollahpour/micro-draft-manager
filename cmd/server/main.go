package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/abdollahpour/almaniha-draft/internal/api"
	"github.com/abdollahpour/almaniha-draft/internal/config"
	"github.com/abdollahpour/almaniha-draft/internal/db"
	"github.com/abdollahpour/almaniha-draft/internal/util"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var Version = "development"

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	version := flag.Bool("version", false, "print version version")
	debug := flag.Bool("debug", false, "active debug manager")

	flag.Parse()
	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}
	if *debug {
		log.SetLevel(log.TraceLevel)
	}

	conf := config.NewEnvConfiguration()

	defer util.DisconnectMongoDB()

	mongoDb := db.NewMongoDb(conf, false)
	handler := api.NewDbApi(mongoDb)

	router := mux.NewRouter()

	router.Handle("/api/v1/drafts", handler.CreateDraft()).Methods(http.MethodPost)
	router.Handle("/api/v1/drafts", handler.DraftsPaginatedHandler()).Methods(http.MethodGet)
	router.Handle("/api/v1/drafts/{draftId}", handler.ReadDraft()).Methods(http.MethodGet)
	// router.Handle("/api/v1/business/{businessId}", handler.DeleteDraft()).Methods(http.MethodDelete)
	router.Handle("/live", handler.LivenessHandler()).Methods(http.MethodGet)

	http.Handle("/", router)

	port := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	fmt.Printf("\nListen on port %s\n", port)
	http.ListenAndServe(port, router)
}
