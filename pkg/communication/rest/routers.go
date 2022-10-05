package rest

import (
	"encoding/json"
	"net/http"

	"github.com/bitly/go-simplejson"

	"github.com/lokidb/engine"
)

type keyValue struct {
	Key   string
	Value string
}

func writeError(errorMessage string, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response := simplejson.New()
	response.Set("error", errorMessage)

	payload, err := response.MarshalJSON()
	if err != nil {
		panic(err)
	}

	w.Write(payload)
}

func handleSet(w http.ResponseWriter, r *http.Request, engine *engine.KeyValueStore) {
	switch r.Method {
	case "POST":
		// extraxt key and value from body
		decoder := json.NewDecoder(r.Body)
		var kv keyValue
		err := decoder.Decode(&kv)
		if err != nil {
			writeError(err.Error(), w)
			return
		}

		// Get value from storage engine
		err = (*engine).Set(kv.Key, []byte(kv.Value))
		if err != nil {
			writeError(err.Error(), w)
			return
		}

		// create json response
		json := simplejson.New()
		json.Set("status", "ok")
		response, _ := json.MarshalJSON()

		// send response
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		w.Write(response)
	default:
		writeError("Method not allowed", w)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request, engine *engine.KeyValueStore) {
	switch r.Method {
	case "GET":
		// extraxt key from query
		key := r.URL.Query().Get("key")

		// Get value from storage engine
		value := (*engine).Get(key, nil)

		// create json response
		json := simplejson.New()
		json.Set("value", string(value))
		response, _ := json.MarshalJSON()

		// send response
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		w.Write(response)
	default:
		writeError("Method not allowed", w)
	}
}

func handleDel(w http.ResponseWriter, r *http.Request, engine *engine.KeyValueStore) {
	switch r.Method {
	case "DELETE":
		// extraxt key from query
		key := r.URL.Query().Get("key")

		// Get value from storage engine
		deleted := (*engine).Del(key)

		// create json response
		json := simplejson.New()
		json.Set("deleted", deleted)
		response, _ := json.MarshalJSON()

		// send response
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		w.Write(response)
	default:
		writeError("Method not allowed", w)
	}
}

func handleFlush(w http.ResponseWriter, r *http.Request, engine *engine.KeyValueStore) {
	switch r.Method {
	case "DELETE":
		// Get value from storage engine
		(*engine).Flush()

		// send response
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte{})
	default:
		writeError("Method not allowed", w)
	}
}

func handleKeys(w http.ResponseWriter, r *http.Request, engine *engine.KeyValueStore) {
	switch r.Method {
	case "GET":
		// Get value from storage engine
		keys := (*engine).Keys()

		// create json response
		json := simplejson.New()
		json.Set("keys", keys)
		response, _ := json.MarshalJSON()

		// send response
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		w.Write(response)
	default:
		writeError("Method not allowed", w)
	}
}
