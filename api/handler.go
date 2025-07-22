package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jatin9996/go-cache-ttl/cache"
)

type Handler struct {
	Cache *cache.Cache
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	val, ok := h.Cache.Get(key)

	if !ok {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"value": val})

}

func (h *Handler) Set(w http.ResponseWriter, r *http.Request) {
	var body struct {
		key   string      `json:"key"`
		Value interface{} `json:"value"`
		TTL   int         `json:"ttl`
	}

	json.NewDecoder(r.Body).Decode(&body)

	if body.TTL > 0 {
		h.Cache.Set(body.key, body.Value, time.Duration(body.TTL)*time.Second)
	} else {
		h.Cache.Set(body.key, body.Value)
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	h.Cache.Delete(key)
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Stats(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(h.Cache.Stats())
}
