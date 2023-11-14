package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var serviceDiscoveryURL string

var routingThreshold int

var ss sync.Map
var concurrentTasks = make(chan bool, 5)

var rdb *redis.Client

var registry = make(map[string][]string)
var counter = make(map[string]int)
var registryMutex sync.Mutex

var requestCounts = make(map[int]int)

func makeRequestWithRouting(service string, method string, path string, body io.Reader) (*http.Response, string) {
	var forwards = 0

	for forwards <= routingThreshold {
		addr := roundRobinGetNext(service)
		req, _ := http.NewRequest(method, "http://"+addr+path, body)
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			go circuitBreaker(service, addr)
			forwards++

		} else {
			return resp, addr
		}
	}

	return nil, ""
}

func circuitBreaker(sName string, sAddr string) {

	var stop = time.Now().Add(52 * time.Second)
	var errCount = 0

	for time.Now().Before(stop) {

		req, _ := http.NewRequest(http.MethodGet, "http://"+sAddr+"/status", nil)
		_, err := http.DefaultClient.Do(req)

		if err != nil {
			errCount++
		}

		if errCount == 3 {
			fmt.Println("Service of type \"", sName, "\" found ", sAddr, " is SICK. 3 errors in <= 52 seconds.")
			return
		}

		time.Sleep(300 * time.Millisecond)
	}

	return
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
	requestCounts[200]++
	w.WriteHeader(200)
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
	router.HandleFunc("/metrics", getMetrics).Methods(http.MethodGet)

	router.HandleFunc("/shoes", getShoes).Methods(http.MethodGet)
	router.HandleFunc("/shoes/{id}", getShoesById).Methods(http.MethodGet)
	router.HandleFunc("/shoes", postShoes).Methods(http.MethodPost)
	router.HandleFunc("/shoes/{id}", putShoesById).Methods(http.MethodPut)

	router.HandleFunc("/transaction/{id}", getTransactionsByShoesId).Methods(http.MethodGet)
	router.HandleFunc("/stock/{id}", getStockByShoesId).Methods(http.MethodGet)
	router.HandleFunc("/turnaround/{id}/{opType}", getTurnaround).Methods(http.MethodGet)
	router.HandleFunc("/turnaround/{id}/{opType}/{since}/{until}", getTurnaroundTimePeriod).Methods(http.MethodGet)
	router.HandleFunc("/transaction", postTransaction).Methods(http.MethodPost)
	router.HandleFunc("/transaction/2phase", postTransactionTwoPhase).Methods(http.MethodPost)

	fmt.Println("API Gateway started")
	if err := http.ListenAndServe(":5000", router); err != nil {
		log.Fatal(err)
	}
}

func main() {

	requestCounts[200] = 0
	requestCounts[400] = 0
	requestCounts[404] = 0
	requestCounts[408] = 0
	requestCounts[500] = 0
	requestCounts[429] = 0

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

	routingThresholdString, found1 := os.LookupEnv("ROUTING_THRESHOLD")

	if found1 == false {
		fmt.Println("ROUTING_THRESHOLD ENV variable not set!")
		return
	}

	routingThresholdInt, err := strconv.Atoi(routingThresholdString)

	if err != nil {
		fmt.Println("ROUTING_THRESHOLD must be an INT.")
		return
	}

	routingThreshold = routingThresholdInt

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
