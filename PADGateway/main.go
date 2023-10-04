package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var serviceDiscoveryURL string

var ss sync.Map
var concurrentTasks = make(chan bool, 5)

var rdb *redis.Client

var registry = make(map[string][]string)
var counter = make(map[string]int)
var registryMutex sync.Mutex

func roundRobinGetNext(service string) string {
	registryMutex.Lock()
	defer registryMutex.Unlock()

	if len(registry[service]) == 0 {
		return ""
	}

	counter[service] = (counter[service] + 1) % len(registry[service])
	return registry[service][counter[service]]
}

func getRegistryFromServiceDiscovery() {
	for {
		time.Sleep(5 * time.Second)

		req, _ := http.NewRequest(http.MethodGet, serviceDiscoveryURL+"/registry", nil)
		response, err := http.DefaultClient.Do(req)

		if err != nil {
			continue
		}

		registryMutex.Lock()

		json.NewDecoder(response.Body).Decode(&registry)
		//fmt.Println(registry)

		for service, _ := range registry {
			counter[service] = 0
		}

		registryMutex.Unlock()
	}
}

//func makeRequest() {
//
//	//req, _ := http.NewRequest(http.MethodPost, "http://localhost:5000/transaction", strings.NewReader("{\"shoesId\":\"d0d24b92-7df8-45ca-a0e7-cee6a05e1ffc\",\"quantity\":\"15\",\"operationType\": -1}"))
//	req, _ := http.NewRequest(http.MethodPost, "http://localhost:5001/register/catalog/1", nil)
//	response, err := http.DefaultClient.Do(req)
//	if err != nil {
//		fmt.Println("Error:", err.Error())
//		return
//	}
//	defer response.Body.Close()
//
//	fmt.Println(response.StatusCode)
//}

func getStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	return
}

func takeFromChannel() {
	select {
	case <-concurrentTasks:
		return
	default:
		return
	}
}

func runServer() {
	router := mux.NewRouter()

	router.HandleFunc("/status", getStatus).Methods(http.MethodGet)

	router.HandleFunc("/shoes", getShoes).Methods(http.MethodGet)
	router.HandleFunc("/shoes/{id}", getShoesById).Methods(http.MethodGet)
	router.HandleFunc("/shoes", postShoes).Methods(http.MethodPost)
	router.HandleFunc("/shoes/{id}", putShoesById).Methods(http.MethodPut)

	router.HandleFunc("/transaction/{id}", getTransactionsByShoesId).Methods(http.MethodGet)
	router.HandleFunc("/stock/{id}", getStockByShoesId).Methods(http.MethodGet)
	router.HandleFunc("/turnaround/{id}/{opType}", getTurnaround).Methods(http.MethodGet)
	router.HandleFunc("/turnaround/{id}/{opType}/{since}/{until}", getTurnaroundTimePeriod).Methods(http.MethodGet)
	router.HandleFunc("/transaction", postTransaction).Methods(http.MethodPost)

	fmt.Println("API Gateway started")
	if err := http.ListenAndServe(":5000", router); err != nil {
		log.Fatal(err)
	}
}

func main() {
	redisURL, found := os.LookupEnv("REDIS_URL")

	if found == false {
		fmt.Println("REDIS_URL ENV variable not set!")
		return
	}

	serviceDiscoveryURL, found = os.LookupEnv("SERVICE_DISCOVERY_URL")

	if found == false {
		fmt.Println("SERVICE_DISCOVERY_URL ENV variable not set!")
		return
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	go runServer()

	go getRegistryFromServiceDiscovery()

	//time.Sleep(3 * time.Second)
	//
	//fmt.Println("starting request")
	//for i := 0; i < 5; i++ {
	//	go makeRequest()
	//}

	select {}
}
