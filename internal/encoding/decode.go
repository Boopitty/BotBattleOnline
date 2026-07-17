package encoding

import (
	"encoding/json"
	"log"
	"net/http"
)

// Decodes request into a given struct.
// Handles error internally, returns success state in boolean.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dest any) bool {
	success := true
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dest)
	if err != nil {
		success = false
		log.Printf("Error decoding request: %v", err)
		RespondWithError(w, http.StatusInternalServerError, err)
		return success
	}
	return success
}
