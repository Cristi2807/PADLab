package main

import (
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Initialize(t *testing.T) {
	redisURL, found := os.LookupEnv("REDIS_URL")

	if found == false {
		t.Fatal("REDIS_URL ENV variable not set!")
	}

	serviceDiscoveryURL, found = os.LookupEnv("SERVICE_DISCOVERY_URL")

	if found == false {
		t.Fatal("SERVICE_DISCOVERY_URL ENV variable not set!")
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	req, _ := http.NewRequest(http.MethodGet, "http://"+serviceDiscoveryURL+"/registry", nil)
	response, err := http.DefaultClient.Do(req)

	if err == nil {
		registryMutex.Lock()

		json.NewDecoder(response.Body).Decode(&registry)
		//fmt.Println(registry)

		registryMutex.Unlock()
	}
}

func TestGetStatus(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getStatus)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "getStatus handler returned wrong status code")
}

func TestGetShoes(t *testing.T) {
	Initialize(t)

	req, err := http.NewRequest(http.MethodGet, "/shoes", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getShoes)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "getShoes handler returned wrong status code")

}

func TestCreateAndGetShoes(t *testing.T) {
	Initialize(t)

	req, err := http.NewRequest(http.MethodPost, "/shoes",
		strings.NewReader("{\"color\":\"red\",\"size\":\"38\",\"price\": \"123.5\", \"brand\": \"Gucci\", "+
			"\"category\": \"sport\", \"model\": \"ab-46\",}"))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(postShoes)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "postShoes handler returned wrong status code")

	type ShoesResponse struct {
		ShoesId string `json:"id"`
	}

	var cell ShoesResponse

	json.NewDecoder(rr.Body).Decode(&cell)

	req1, err1 := http.NewRequest(http.MethodGet, "/shoes/"+cell.ShoesId, nil)
	if err1 != nil {
		t.Fatal(err1)
	}

	rr1 := httptest.NewRecorder()
	handler1 := http.HandlerFunc(getShoes)

	handler1.ServeHTTP(rr1, req1)

	assert.Equal(t, http.StatusOK, rr.Code, "getShoesById handler returned wrong status code")
}
