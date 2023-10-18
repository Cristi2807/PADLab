package main

import (
	"github.com/redis/go-redis/v9"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	req, err := http.NewRequest(http.MethodGet, "/shoes", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getShoes)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "getShoes handler returned wrong status code")

}
