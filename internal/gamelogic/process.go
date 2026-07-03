package gamelogic

import (
	"encoding/json"
	"fmt"
)

func Process(req string) string {
	// Parse the request string into a Request struct
	request := Request{}
	err := json.Unmarshal([]byte(req), &request)
	if err != nil {
		return fmt.Sprintf("Error parsing command: %s", err)
	}

	report := commands(&request)
	return fmt.Sprintf("Report: %s", report)
}
