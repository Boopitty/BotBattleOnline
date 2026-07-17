package encoding

import (
	"encoding/json"
	"log"
)

// Marshal a given struct and return it in []byte form
func MakeJSONResponse(payload interface{}) []byte {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return []byte(`{"error": "Internal Server Error"}`)
	}
	return dat
}

// Place a strint into a pre-made response struct and return it.
func MakeErrorResponse(payload string) any {
	resp := struct {
		message string
	}{
		message: payload,
	}
	return resp
}
