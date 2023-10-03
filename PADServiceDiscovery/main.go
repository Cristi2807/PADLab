package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"sync"
)

var registry = make(map[string][]string)
var m sync.Mutex

func containsValue(slice []string, target string) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

func getStatus(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
	return
}

func registerService(w http.ResponseWriter, r *http.Request) {
	m.Lock()
	defer m.Unlock()

	params := mux.Vars(r)

	if containsValue(registry[params["serviceType"]], params["serviceURL"]) {
		w.WriteHeader(400)
		w.Write([]byte("Already exists"))

		return
	}

	registry[params["serviceType"]] = append(registry[params["serviceType"]], params["serviceURL"])

	w.WriteHeader(200)
	return
}

func getServices(w http.ResponseWriter, _ *http.Request) {

	data, _ := json.Marshal(registry)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)

	return
}

func runServer() {
	router := mux.NewRouter()

	router.HandleFunc("/status", getStatus).Methods(http.MethodGet)

	router.HandleFunc("/registry/{serviceType}/{serviceURL}", registerService).Methods(http.MethodPost)
	router.HandleFunc("/registry", getServices).Methods(http.MethodGet)

	fmt.Println("Service Discovery started")
	if err := http.ListenAndServe(":5001", router); err != nil {
		log.Fatal(err)
	}
}

func main() {

	go runServer()

	select {}
}
