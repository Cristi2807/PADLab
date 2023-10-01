package main

import (
	"context"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func getShoes(w http.ResponseWriter, r *http.Request) {

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

	val, err := rdb.Get(context.Background(), "/shoes").Bytes()
	if err == nil {
		//fmt.Println("Using cache")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(val)

		return
	}

	req, _ := http.NewRequest(r.Method, "http://localhost:5050/shoes", r.Body)
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
		//fmt.Println("Saving in cache")
		rdb.Set(context.Background(), "/shoes", body, 0)
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	return
}

func getShoesById(w http.ResponseWriter, r *http.Request) {

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

	val, err := rdb.Get(context.Background(), "/shoes/"+params["id"]).Bytes()
	if err == nil {
		//fmt.Println("Using cache")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(val)

		return
	}

	req, _ := http.NewRequest(r.Method, "http://localhost:5050/shoes/"+params["id"], r.Body)
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
		//fmt.Println("Saving in cache")
		rdb.Set(context.Background(), "/shoes/"+params["id"], body, 0)
	}

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	return
}

func postShoes(w http.ResponseWriter, r *http.Request) {

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

	req, _ := http.NewRequest(r.Method, "http://localhost:5050/shoes", r.Body)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("\"An internal error happened. Try again later\""))

		return
	}

	if resp.StatusCode == http.StatusOK {
		//fmt.Println("Deleting from cache")
		rdb.Del(context.Background(), "/shoes")
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	return
}

func putShoesById(w http.ResponseWriter, r *http.Request) {

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
	req, _ := http.NewRequest(r.Method, "http://localhost:5050/shoes/"+params["id"], r.Body)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("\"An internal error happened. Try again later\""))

		return
	}

	if resp.StatusCode == http.StatusOK {
		//fmt.Println("Deleting from cache")
		rdb.Del(context.Background(), "/shoes")
		rdb.Del(context.Background(), "/shoes/"+params["id"])
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	w.Write(body)

	return
}
