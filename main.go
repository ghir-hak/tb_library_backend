package lib

import (
	"fmt"
	"time"

	"github.com/taubyte/go-sdk/event"
)

//export init
func init() {
	fmt.Printf("[DEBUG] Pixame backend initializing...\n")
	
	// Initialize databases
	if initDatabases() != 0 {
		fmt.Printf("[ERROR] Failed to initialize databases\n")
		return
	}
	
	fmt.Printf("[DEBUG] Pixame backend initialized successfully\n")
}

//export health
func health(e event.Event) uint32 {
	fmt.Printf("[DEBUG] health check called\n")
	
	h, err := e.HTTP()
	if err != nil {
		fmt.Printf("[ERROR] health HTTP error: %v\n", err)
		return 1
	}
	setCORSHeaders(h)
	
	// Return health status
	response := Response{
		Success: true,
		Data: map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().UnixMilli(),
			"service":   "pixame-backend",
		},
	}
	
	fmt.Printf("[DEBUG] health check completed successfully\n")
	return sendJSONResponse(h, response)
}

