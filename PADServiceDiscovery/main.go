package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"sync"
)

var registry = make(map[string][]string)
var m sync.Mutex

var requestCounts = make(map[int]int)

func containsValue(slice []string, target string) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

func getStatus(w http.ResponseWriter, _ *http.Request) {
	requestCounts[200]++
	w.WriteHeader(200)
	return
}

func registerService(w http.ResponseWriter, r *http.Request) {
	m.Lock()
	defer m.Unlock()

	params := mux.Vars(r)

	if containsValue(registry[params["serviceType"]], params["serviceURL"]) {
		requestCounts[400]++
		w.WriteHeader(400)
		w.Write([]byte("Already exists"))

		return
	}

	fmt.Println("Registering service of type \"", params["serviceType"], "\" on address ", params["serviceURL"])
	registry[params["serviceType"]] = append(registry[params["serviceType"]], params["serviceURL"])

	requestCounts[200]++
	w.WriteHeader(200)
	return
}

func getServices(w http.ResponseWriter, _ *http.Request) {

	data, _ := json.Marshal(registry)

	w.Header().Set("Content-Type", "application/json")
	requestCounts[200]++
	w.WriteHeader(200)
	w.Write(data)

	return
}

func getMetrics(w http.ResponseWriter, _ *http.Request) {

	arr := []string{"# HELP http_requests_total The total number of HTTP requests.", "# TYPE http_requests_total counter"}

	for statusCode, nr := range requestCounts {
		arr = append(arr, fmt.Sprintf("http_requests_total{code=\"%d\"} %d", statusCode, nr))

	}

	result := strings.Join(arr, "\n")

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	w.Write([]byte(result))

	return
}

func runServer() {
	router := mux.NewRouter()

	router.HandleFunc("/status", getStatus).Methods(http.MethodGet)

	router.HandleFunc("/registry/{serviceType}/{serviceURL}", registerService).Methods(http.MethodPost)
	router.HandleFunc("/registry", getServices).Methods(http.MethodGet)
	router.HandleFunc("/metrics", getMetrics).Methods(http.MethodGet)

	fmt.Println("Service Discovery started")
	if err := http.ListenAndServe(":5001", router); err != nil {
		log.Fatal(err)
	}
}

func main() {

	requestCounts[200] = 0
	requestCounts[400] = 0
	//registry["catalog"] = []string{"http://localhost:5050"}
	//registry["inventory"] = []string{"http://localhost:7070"}

	go runServer()

	select {}
}
