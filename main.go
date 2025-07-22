package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jatin9996/go-cache-ttl/api"
	"github.com/jatin9996/go-cache-ttl/cache"
)

func main() {
	c := cache.NewCache(10*time.Second, 100) // 10s TTL, 100 entries max
	h := &api.Handler{Cache: c}

	http.HandleFunc("/get", h.Get)
	http.HandleFunc("/set", h.Set)
	http.HandleFunc("/delete", h.Delete)
	http.HandleFunc("/stats", h.Stats)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
