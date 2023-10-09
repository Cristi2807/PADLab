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

type Status struct {
	errors []int64
	sick   bool
}

var statusMutex sync.Mutex
var status = make(map[string]*Status)

func registerError(sName string, sAddr string) {

	statusMutex.Lock()
	defer statusMutex.Unlock()

	if status[sName+sAddr] == nil {
		status[sName+sAddr] = &Status{sick: false, errors: make([]int64, 0)}
	}

	if status[sName+sAddr].sick == false {
		status[sName+sAddr].errors = append(status[sName+sAddr].errors, time.Now().UnixMilli())

		if len(status[sName+sAddr].errors) == 4 {
			status[sName+sAddr].errors = status[sName+sAddr].errors[1:]
		}

		if len(status[sName+sAddr].errors) == 3 && (status[sName+sAddr].errors[2]-status[sName+sAddr].errors[0] <= 52*time.Second.Milliseconds()) {
			status[sName+sAddr].sick = true
			fmt.Println("Service of type \"", sName, "\" found ", sAddr, " is SICK. 3 errors in <= 52 seconds.")
		}
	}
}

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

		req, _ := http.NewRequest(http.MethodGet, "http://"+serviceDiscoveryURL+"/registry", nil)
		response, err := http.DefaultClient.Do(req)

		if err != nil {
			continue
		}

		registryMutex.Lock()

		json.NewDecoder(response.Body).Decode(&registry)
		//fmt.Println(registry)

		registryMutex.Unlock()
	}
}

//func makeRequest() {
//
//	//req, _ := http.NewRequest(http.MethodPost, "http://localhost:5000/transaction", strings.NewReader("{\"shoesId\":\"d0d24b92-7df8-45ca-a0e7-cee6a05e1ffc\",\"quantity\":\"15\",\"operationType\": -1}"))
//	req, _ := http.NewRequest(http.MethodGet, "http://localhost:5000/shoes", nil)
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

	//time.Sleep(6 * time.Second)
	//
	//fmt.Println("starting request")
	//for i := 0; i < 7; i++ {
	//	go makeRequest()
	//}

	select {}
}
