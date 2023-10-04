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

var statusMutex sync.Mutex

type Status struct {
	sick       bool
	errorCount int
	okCount    int
}

var status = make(map[string]*Status)

func checkRegistryHealth() {
	for {
		time.Sleep(2 * time.Second)
		//fmt.Println("Checking health of microservices...")

		registryMutex.Lock()
		reg := registry
		registryMutex.Unlock()

		statusMutex.Lock()

		for service, urlList := range reg {
			for i := 0; i < len(urlList); i++ {
				req, _ := http.NewRequest(http.MethodGet, "http://"+urlList[i]+"/status", nil)
				_, err := http.DefaultClient.Do(req)

				if status[service+urlList[i]] == nil {
					status[service+urlList[i]] = &Status{sick: false, okCount: 0, errorCount: 0}
				}

				// if now error, not sick then increment error count
				if err != nil && status[service+urlList[i]].sick == false {
					status[service+urlList[i]].errorCount++
				}

				// reset error count , if now not error and not sick
				if err == nil && status[service+urlList[i]].sick == false {
					status[service+urlList[i]].errorCount = 0
				}

				// if now not error, but sick increase okCount
				if err == nil && status[service+urlList[i]].sick == true {
					status[service+urlList[i]].okCount++
					status[service+urlList[i]].errorCount = 0
				}

				// reset okCount, if now error and sick
				if err != nil && status[service+urlList[i]].sick == true {
					status[service+urlList[i]].okCount = 0
					//status[service+urlList[i]].errorCount++
				}

				if status[service+urlList[i]].sick == false && status[service+urlList[i]].errorCount == 3 {
					fmt.Println(service, urlList[i], "didn't respond 3 times in a row. Marking it as sick.")
					status[service+urlList[i]].sick = true
					status[service+urlList[i]].errorCount = 0
					status[service+urlList[i]].okCount = 0
				}

				if status[service+urlList[i]].sick == true && status[service+urlList[i]].okCount == 10 {
					fmt.Println(service, urlList[i], "responded 10 times in a row. Marking it as healthy again.")
					status[service+urlList[i]].sick = false
					status[service+urlList[i]].errorCount = 0
					status[service+urlList[i]].okCount = 0
				}

				//if status[service+urlList[i]].sick == true && status[service+urlList[i]].errorCount == 20 {
				//	fmt.Println("delete at all", service, urlList)
				//	status[service+urlList[i]].errorCount = 0
				//	status[service+urlList[i]].okCount = 0
				//}

			}
		}
		statusMutex.Unlock()
	}
}

func roundRobinGetNext(service string) string {
	registryMutex.Lock()
	defer registryMutex.Unlock()

	if len(registry[service]) == 0 {
		return ""
	}

	counter[service] = (counter[service] + 1) % len(registry[service])

	statusMutex.Lock()

	// get only healthy instance, if available
	i := 0
	for status[service+registry[service][counter[service]]] != nil &&
		status[service+registry[service][counter[service]]].sick == true &&
		i < len(registry[service]) {
		counter[service] = (counter[service] + 1) % len(registry[service])
		i++
	}

	statusMutex.Unlock()

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

	go checkRegistryHealth()

	//time.Sleep(6 * time.Second)
	//
	//fmt.Println("starting request")
	//for i := 0; i < 7; i++ {
	//	go makeRequest()
	//}

	select {}
}
