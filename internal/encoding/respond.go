package encoding

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Write an error response and log it in the server
func RespondWithError(w http.ResponseWriter, code int, err error) {
	log.Printf("Error: %v\n", err)
	w.WriteHeader(code)
	w.Write([]byte("Error Code: " + fmt.Sprint(code)))
}

// Write a response in JSON format
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}
