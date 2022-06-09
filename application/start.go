package application

import (
	"encoding/json"
	"fmt"
	"github.com/dulguundd/logError-lib/logger"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func Start() {
	environment := getEnvironment()

	router := mux.NewRouter()

	h := Handler{newHandler(environment.zeebeConfig)}

	//define device routes
	router.HandleFunc("/info", h.TopologyInfo).Methods(http.MethodGet)
	router.HandleFunc("/resource", h.DeployResource).Methods(http.MethodPost)
	router.HandleFunc("/instance", h.DeployInstance).Methods(http.MethodPost)

	//starting server
	log.Fatal(http.ListenAndServe(environment.serviceConfig.address, router))
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

func serviceLatencyLogger(start time.Time) {
	elapsed := time.Since(start)
	logMessage := fmt.Sprintf("response latencie %s", elapsed)
	logger.Info(logMessage)
}
