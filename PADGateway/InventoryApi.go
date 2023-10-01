package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"sync"
)

func getTransactionsByShoesId(w http.ResponseWriter, r *http.Request) {

	select {
	case concurrentTasks <- true:
		defer takeFromChannel()
	default:
		//fmt.Println("All resources taken. Not serving your request 429")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("\"Concurrency limit achieved.\""))
		return
	}

	defer r.Body.Close()
	params := mux.Vars(r)

	val, err := rdb.Get(context.Background(), "/transaction/"+params["id"]).Bytes()
	if err == nil {
		fmt.Println("Using cache")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(val)

		return
	}

	req, _ := http.NewRequest(r.Method, "http://localhost:7070/transaction/"+params["id"], r.Body)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("\"An internal error happened. Try again later\""))

		return
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Saving in cache")
		rdb.Set(context.Background(), "/transaction/"+params["id"], body, 0)
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	return
}

func getStockByShoesId(w http.ResponseWriter, r *http.Request) {

	select {
	case concurrentTasks <- true:
		defer takeFromChannel()
	default:
		//fmt.Println("All resources taken. Not serving your request 429")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("\"Concurrency limit achieved.\""))
		return
	}

	defer r.Body.Close()
	params := mux.Vars(r)

	val, err := rdb.Get(context.Background(), "/stock/"+params["id"]).Bytes()
	if err == nil {
		fmt.Println("Using cache")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(val)

		return
	}

	req, _ := http.NewRequest(r.Method, "http://localhost:7070/stock/"+params["id"], r.Body)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("\"An internal error happened. Try again later\""))

		return
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Saving in cache")
		rdb.Set(context.Background(), "/stock/"+params["id"], body, 0)
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	return
}

func getTurnaround(w http.ResponseWriter, r *http.Request) {

	select {
	case concurrentTasks <- true:
		defer takeFromChannel()
	default:
		//fmt.Println("All resources taken. Not serving your request 429")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("\"Concurrency limit achieved.\""))
		return
	}

	defer r.Body.Close()
	params := mux.Vars(r)

	val, err := rdb.Get(context.Background(), "/turnaround/"+params["id"]+"/"+params["opType"]).Bytes()
	if err == nil {
		fmt.Println("Using cache")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(val)

		return
	}

	req, _ := http.NewRequest(r.Method, "http://localhost:7070/turnaround/"+params["id"]+"/"+params["opType"], r.Body)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("\"An internal error happened. Try again later\""))

		return
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Saving in cache")
		rdb.Set(context.Background(), "/turnaround/"+params["id"]+"/"+params["opType"], body, 0)
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	return
}

func getTurnaroundTimePeriod(w http.ResponseWriter, r *http.Request) {

	select {
	case concurrentTasks <- true:
		defer takeFromChannel()
	default:
		//fmt.Println("All resources taken. Not serving your request 429")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("\"Concurrency limit achieved.\""))
		return
	}

	defer r.Body.Close()

	params := mux.Vars(r)
	req, _ := http.NewRequest(r.Method, "http://localhost:7070/turnaround/"+params["id"]+"/"+params["opType"]+"/"+params["since"]+"/"+params["until"], r.Body)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("\"An internal error happened. Try again later\""))

		return
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	return
}

func postTransaction(w http.ResponseWriter, r *http.Request) {

	select {
	case concurrentTasks <- true:
		defer takeFromChannel()
	default:
		//fmt.Println("All resources taken. Not serving your request 429")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("\"Concurrency limit achieved.\""))
		return
	}

	reqBytes, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	type Transaction struct {
		ShoesId       string `json:"shoesId"`
		OperationType int    `json:"operationType"`
	}

	var cell Transaction

	rBody := io.NopCloser(bytes.NewReader(reqBytes))
	err := json.NewDecoder(rBody).Decode(&cell)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//checking if such product exists
	req1, _ := http.NewRequest(http.MethodGet, "http://localhost:5050/shoes/"+cell.ShoesId, nil)
	resp1, err2 := http.DefaultClient.Do(req1)

	if err2 != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("\"An internal error happened. Try again later\""))

		return
	}

	if resp1.StatusCode == http.StatusNotFound {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("\"Cannot write transaction for non-existing in catalog product\""))

		return
	}

	if cell.OperationType == -1 {
		res, _ := ss.LoadOrStore(cell.ShoesId, &sync.Mutex{})
		mutex := res.(*sync.Mutex)
		//fmt.Println("Trying to lock mutex ", cell.ShoesId, " ", time.Now())
		mutex.Lock()
		//fmt.Println("Mutex ", cell.ShoesId, "locked at", time.Now())
		defer mutex.Unlock()
	}

	rBody = io.NopCloser(bytes.NewReader(reqBytes))
	req, _ := http.NewRequest(r.Method, "http://localhost:7070/transaction", rBody)
	resp, err1 := http.DefaultClient.Do(req)

	if err1 != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("\"An internal error happened. Try again later\""))

		return
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Deleting from cache")
		rdb.Del(context.Background(), "/transaction/"+cell.ShoesId)
		rdb.Del(context.Background(), "/stock/"+cell.ShoesId)
		rdb.Del(context.Background(), "/turnaround/"+cell.ShoesId+"/"+strconv.Itoa(cell.OperationType))
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	//if cell.OperationType == -1 {
	//	mutex.Unlock()
	//	//time.Sleep(time.Millisecond)
	//	//if mutex.TryLock() == true {
	//	//	mutex.Unlock()
	//	//	ss.Delete("1")
	//	//	fmt.Println("No other thread in queue, deleting mutex from sync.Map")
	//	//}
	//}

	//if cell.OperationType == -1 {
	//	fmt.Println("Exiting mutex ", cell.ShoesId, " at ", time.Now())
	//}

	return
}