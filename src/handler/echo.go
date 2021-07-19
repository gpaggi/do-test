package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	models "echoapi/models"

	log "github.com/sirupsen/logrus"
)

// doEcho validates an input JSON and adds the top level echoded: true field to it.
// it returns a Response containing the HttpCode and output.
func doEcho(r io.ReadCloser) *models.Response {
	decoder := json.NewDecoder(r)
	var t map[string]interface{}
	err := decoder.Decode(&t)

	// Validate input
	// Check for unmarshalling errors
	if err != nil {
		log.Errorf("Error while unmarshalling POST input: %v", err.Error())
		return &models.Response{
			HttpCode: 400,
			Output:   json.RawMessage(`{"error": "Malformed request, input must be valid JSON"}`),
		}
	}
	// Check if the top-level field echoed is already set to true
	if val, ok := t["echoed"]; ok && val == "true" {
		log.Debugf("Top level echoed field is already set to true, returning error")
		return &models.Response{
			HttpCode: 400,
			Output:   json.RawMessage(`{"error": "Top level echoed field is already set to true"}`),
		}
	}

	// Set top level echoed field to true
	t["echoed"] = "true"

	tStr, err := json.Marshal(t)
	if err != nil {
		log.Error(err)
	}

	// Return the response
	return &models.Response{
		HttpCode: 200,
		Output:   json.RawMessage(tStr),
	}
}

// HTTP Handler which wraps the doEcho function to simplify unit testing
func EchoHandler(w http.ResponseWriter, r *http.Request) {
	ret := doEcho(r.Body)
	// Set response header to json
	w.Header().Set("Content-Type", "application/json")
	// set HTTP code
	w.WriteHeader(ret.HttpCode)
	// Write out the response
	json.NewEncoder(w).Encode(ret.Output)
}
