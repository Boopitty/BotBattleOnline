package encoding

import (
	"encoding/json"
	"log"
)

func MakeJSONResponse(payload interface{}) []byte {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return []byte(`{"error": "Internal Server Error"}`)
	}
	return dat
}

func MakeErrorResponse(payload string) any {
	resp := struct {
		message string
	}{
		message: payload,
	}
	return resp
}
