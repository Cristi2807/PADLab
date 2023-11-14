package main

import (
	"context"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func getShoes(w http.ResponseWriter, r *http.Request) {

	//defer r.Body.Close()

	val, err := rdb.Get(context.Background(), "/shoes").Bytes()
	if err == nil {
		//fmt.Println("Using cache")
		w.Header().Set("Content-Type", "application/json")
		requestCounts[200]++
		w.WriteHeader(http.StatusOK)
		w.Write(val)

		return
	}

	select {
	case concurrentTasks <- true:
		defer takeFromChannel()
	default:
		//fmt.Println("All resources taken. Not serving your request 429")
		w.Header().Set("Content-Type", "application/json")
		requestCounts[429]++
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("\"Concurrency limit achieved.\""))
		return
	}

	var forwards = 0

	for forwards <= routingThreshold {
		addr := roundRobinGetNext("catalog")
		req, _ := http.NewRequest(r.Method, "http://"+addr+"/shoes", r.Body)
		resp, err1 := http.DefaultClient.Do(req)

		if err1 != nil {
			go circuitBreaker("catalog", addr)
			forwards++

		} else {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode == http.StatusOK {
				//fmt.Println("Saving in cache")
				rdb.Set(context.Background(), "/shoes", body, 0)
			}

			w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
			requestCounts[resp.StatusCode]++
			w.WriteHeader(resp.StatusCode)
			w.Write(body)

			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	requestCounts[500]++
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("\"An internal error happened. Try again later\""))

	return

}

func getShoesById(w http.ResponseWriter, r *http.Request) {

	//defer r.Body.Close()
	params := mux.Vars(r)

	val, err := rdb.Get(context.Background(), "/shoes/"+params["id"]).Bytes()
	if err == nil {
		//fmt.Println("Using cache")
		w.Header().Set("Content-Type", "application/json")
		requestCounts[200]++
		w.WriteHeader(http.StatusOK)
		w.Write(val)

		return
	}

	select {
	case concurrentTasks <- true:
		defer takeFromChannel()
	default:
		//fmt.Println("All resources taken. Not serving your request 429")
		w.Header().Set("Content-Type", "application/json")
		requestCounts[429]++
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("\"Concurrency limit achieved.\""))
		return
	}

	var forwards = 0

	for forwards <= routingThreshold {
		addr := roundRobinGetNext("catalog")
		req, _ := http.NewRequest(r.Method, "http://"+addr+"/shoes/"+params["id"], r.Body)
		resp, err1 := http.DefaultClient.Do(req)

		if err1 != nil {
			go circuitBreaker("catalog", addr)
			forwards++

		} else {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode == http.StatusOK {
				//fmt.Println("Saving in cache")
				rdb.Set(context.Background(), "/shoes/"+params["id"], body, 0)
			}

			w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
			requestCounts[resp.StatusCode]++
			w.WriteHeader(resp.StatusCode)
			w.Write(body)

			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	requestCounts[500]++
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("\"An internal error happened. Try again later\""))

	return

}

func postShoes(w http.ResponseWriter, r *http.Request) {

	select {
	case concurrentTasks <- true:
		defer takeFromChannel()
	default:
		//fmt.Println("All resources taken. Not serving your request 429")
		w.Header().Set("Content-Type", "application/json")
		requestCounts[429]++
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("\"Concurrency limit achieved.\""))
		return
	}

	defer r.Body.Close()

	var forwards = 0

	for forwards <= routingThreshold {
		addr := roundRobinGetNext("catalog")
		req, _ := http.NewRequest(r.Method, "http://"+addr+"/shoes", r.Body)
		resp, err1 := http.DefaultClient.Do(req)

		if err1 != nil {
			go circuitBreaker("catalog", addr)
			forwards++

		} else {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode == http.StatusOK {
				//fmt.Println("Saving in cache")
				rdb.Del(context.Background(), "/shoes")
			}

			w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
			requestCounts[resp.StatusCode]++
			w.WriteHeader(resp.StatusCode)
			w.Write(body)

			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	requestCounts[500]++
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("\"An internal error happened. Try again later\""))

	return
}

func putShoesById(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	select {
	case concurrentTasks <- true:
		defer takeFromChannel()
	default:
		//fmt.Println("All resources taken. Not serving your request 429")
		w.Header().Set("Content-Type", "application/json")
		requestCounts[429]++
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("\"Concurrency limit achieved.\""))
		return
	}

	defer r.Body.Close()

	var forwards = 0

	for forwards <= routingThreshold {
		addr := roundRobinGetNext("catalog")
		req, _ := http.NewRequest(r.Method, "http://"+addr+"/shoes/"+params["id"], r.Body)
		resp, err1 := http.DefaultClient.Do(req)

		if err1 != nil {
			go circuitBreaker("catalog", addr)
			forwards++

		} else {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode == http.StatusOK {
				//fmt.Println("Deleting from cache")
				rdb.Del(context.Background(), "/shoes")
				rdb.Del(context.Background(), "/shoes/"+params["id"])
			}

			w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
			requestCounts[resp.StatusCode]++
			w.WriteHeader(resp.StatusCode)
			w.Write(body)

			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	requestCounts[500]++
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("\"An internal error happened. Try again later\""))

	return
}
