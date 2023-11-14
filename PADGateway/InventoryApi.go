package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"sync"
)

func getTransactionsByShoesId(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	params := mux.Vars(r)

	val, err := rdb.Get(context.Background(), "/transaction/"+params["id"]).Bytes()
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
		addr := roundRobinGetNext("inventory")
		req, _ := http.NewRequest(r.Method, "http://"+addr+"/transaction/"+params["id"], r.Body)
		resp, err1 := http.DefaultClient.Do(req)

		if err1 != nil {
			go circuitBreaker("inventory", addr)
			forwards++

		} else {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode == http.StatusOK {
				//fmt.Println("Saving in cache")
				rdb.Set(context.Background(), "/transaction/"+params["id"], body, 0)
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

func getStockByShoesId(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	params := mux.Vars(r)

	val, err := rdb.Get(context.Background(), "/stock/"+params["id"]).Bytes()
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
		addr := roundRobinGetNext("inventory")
		req, _ := http.NewRequest(r.Method, "http://"+addr+"/stock/"+params["id"], r.Body)
		resp, err1 := http.DefaultClient.Do(req)

		if err1 != nil {
			go circuitBreaker("inventory", addr)
			forwards++

		} else {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode == http.StatusOK {
				//fmt.Println("Saving in cache")
				rdb.Set(context.Background(), "/stock/"+params["id"], body, 0)
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

func getTurnaround(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	params := mux.Vars(r)

	val, err := rdb.Get(context.Background(), "/turnaround/"+params["id"]+"/"+params["opType"]).Bytes()
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
		addr := roundRobinGetNext("inventory")
		req, _ := http.NewRequest(r.Method, "http://"+addr+"/turnaround/"+params["id"]+"/"+params["opType"], r.Body)
		resp, err1 := http.DefaultClient.Do(req)

		if err1 != nil {
			go circuitBreaker("inventory", addr)
			forwards++

		} else {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode == http.StatusOK {
				//fmt.Println("Saving in cache")
				rdb.Set(context.Background(), "/turnaround/"+params["id"]+"/"+params["opType"], body, 0)
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

func getTurnaroundTimePeriod(w http.ResponseWriter, r *http.Request) {

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

	params := mux.Vars(r)

	var forwards = 0

	for forwards <= routingThreshold {
		addr := roundRobinGetNext("inventory")
		req, _ := http.NewRequest(r.Method, "http://"+addr+"/turnaround/"+params["id"]+"/"+params["opType"]+"/"+params["since"]+"/"+params["until"], r.Body)
		resp, err1 := http.DefaultClient.Do(req)

		if err1 != nil {
			go circuitBreaker("inventory", addr)
			forwards++

		} else {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

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

func postTransaction(w http.ResponseWriter, r *http.Request) {

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
	resp1, _ := makeRequestWithRouting("catalog", http.MethodGet, "/shoes/"+cell.ShoesId, nil)

	if resp1 == nil {
		w.Header().Set("Content-Type", "application/json")
		requestCounts[500]++
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("\"An internal error happened. Try again later\""))

		return
	}

	if resp1.StatusCode == http.StatusNotFound {
		w.Header().Set("Content-Type", "application/json")
		requestCounts[404]++
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

	resp, _ := makeRequestWithRouting("inventory", r.Method, "/transaction", rBody)

	if resp == nil {

		w.Header().Set("Content-Type", "application/json")
		requestCounts[500]++
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("\"An internal error happened. Try again later\""))

		return
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		//fmt.Println("Deleting from cache")
		rdb.Del(context.Background(), "/transaction/"+cell.ShoesId)
		rdb.Del(context.Background(), "/stock/"+cell.ShoesId)
		rdb.Del(context.Background(), "/turnaround/"+cell.ShoesId+"/"+strconv.Itoa(cell.OperationType))
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	requestCounts[resp.StatusCode]++
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	return
}

func postTransactionTwoPhase(w http.ResponseWriter, r *http.Request) {

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

	type Shoes struct {
		Id       string `json:"id,omitempty"`
		Color    string `json:"color"`
		Size     string `json:"size"`
		Price    string `json:"price"`
		Brand    string `json:"brand"`
		Category string `json:"category"`
		Model    string `json:"model"`
	}

	type Transaction struct {
		Id            string `json:"id,omitempty"`
		ShoesId       string `json:"shoesId"`
		OperationType int    `json:"operationType"`
		Quantity      string `json:"quantity"`
	}

	type TwoPhase struct {
		Shoes       Shoes       `json:"shoes"`
		Transaction Transaction `json:"transaction"`
	}

	var cell TwoPhase

	err := json.NewDecoder(r.Body).Decode(&cell)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if cell.Transaction.OperationType == -1 {
		http.Error(w, "Invalid operationType for a new product.", http.StatusBadRequest)
		return
	}

	body1, _ := json.Marshal(cell.Shoes)
	body2 := io.NopCloser(bytes.NewReader(body1))
	resp1, catalogAddr := makeRequestWithRouting("catalog", http.MethodPost, "/shoes/2phase", body2)

	if resp1 == nil {
		w.Header().Set("Content-Type", "application/json")
		requestCounts[500]++
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("\"An internal error happened. Try again later\""))

		return
	}

	if resp1.StatusCode != http.StatusOK {
		w.Header().Set("Content-Type", resp1.Header.Get("Content-Type"))
		requestCounts[resp1.StatusCode]++
		w.WriteHeader(resp1.StatusCode)
		body, _ := io.ReadAll(resp1.Body)
		w.Write(body)

		return
	}

	var shoes Shoes
	json.NewDecoder(resp1.Body).Decode(&shoes)

	cell.Shoes.Id = shoes.Id
	cell.Transaction.ShoesId = shoes.Id
	body3, _ := json.Marshal(cell.Transaction)
	body4 := io.NopCloser(bytes.NewReader(body3))
	resp, inventoryAddr := makeRequestWithRouting("inventory", http.MethodPost, "/transaction/2phase", body4)

	if resp == nil {
		req, _ := http.NewRequest(http.MethodPost, "http://"+catalogAddr+"/rollback/"+shoes.Id, nil)
		http.DefaultClient.Do(req)

		w.Header().Set("Content-Type", "application/json")
		requestCounts[500]++
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("\"An internal error happened. Try again later\""))

		return
	}

	if resp.StatusCode != http.StatusOK {
		req, _ := http.NewRequest(http.MethodPost, "http://"+catalogAddr+"/rollback/"+shoes.Id, nil)
		http.DefaultClient.Do(req)

		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		requestCounts[resp.StatusCode]++
		w.WriteHeader(resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		w.Write(body)

		return
	}

	var transaction Transaction
	json.NewDecoder(resp.Body).Decode(&transaction)
	cell.Transaction.Id = transaction.Id

	req, _ := http.NewRequest(http.MethodPost, "http://"+catalogAddr+"/commit/"+shoes.Id, nil)
	_, err1 := http.DefaultClient.Do(req)

	if err1 != nil {
		req1, _ := http.NewRequest(http.MethodPost, "http://"+inventoryAddr+"/rollback/"+transaction.Id, nil)
		http.DefaultClient.Do(req1)

		w.Header().Set("Content-Type", "application/json")
		requestCounts[500]++
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("\"An internal error happened. Try again later\""))

		return
	}

	req1, _ := http.NewRequest(http.MethodPost, "http://"+inventoryAddr+"/commit/"+transaction.Id, nil)
	http.DefaultClient.Do(req1)

	w.Header().Set("Content-Type", "application/json")
	requestCounts[200]++
	w.WriteHeader(http.StatusOK)
	finalBody, _ := json.Marshal(cell)
	w.Write(finalBody)

	return
}
