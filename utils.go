package lib

import (
	"encoding/json"
	"fmt"
	"time"

	httpEvent "github.com/taubyte/go-sdk/http/event"
)

// setCORSHeaders sets CORS headers for HTTP responses
func setCORSHeaders(h httpEvent.Event) {
	h.Headers().Set("Access-Control-Allow-Origin", "*")
	h.Headers().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	h.Headers().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// handleHTTPError handles HTTP errors and returns appropriate response
func handleHTTPError(h httpEvent.Event, err error, code int) uint32 {
	fmt.Printf("[ERROR] HTTP Error: %v\n", err)
	h.Write([]byte(err.Error()))
	h.Return(code)
	return 1
}

// sendJSONResponse sends a JSON response
func sendJSONResponse(h httpEvent.Event, data interface{}) uint32 {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return handleHTTPError(h, err, 500)
	}
	h.Headers().Set("Content-Type", "application/json")
	h.Write(jsonData)
	h.Return(200)
	return 0
}

// getQueryParam gets a query parameter with default value
func getQueryParam(h httpEvent.Event, key, defaultValue string) string {
	value, err := h.Query().Get(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// getRequiredQueryParam gets a required query parameter
func getRequiredQueryParam(h httpEvent.Event, key string) (string, uint32) {
	value, err := h.Query().Get(key)
	if err != nil {
		h.Write([]byte(fmt.Sprintf("Missing required parameter: %s", key)))
		h.Return(400)
		return "", 1
	}
	return value, 0
}

// getCurrentTimestamp returns current timestamp in milliseconds
func getCurrentTimestamp() int64 {
	return time.Now().UnixMilli()
}
