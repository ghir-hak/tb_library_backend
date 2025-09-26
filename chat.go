package lib

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/taubyte/go-sdk/event"
	pubsub "github.com/taubyte/go-sdk/pubsub/node"
)

//export sendMessage
func sendMessage(e event.Event) uint32 {
	fmt.Printf("[DEBUG] sendMessage called\n")
	
	h, err := e.HTTP()
	if err != nil {
		fmt.Printf("[ERROR] sendMessage HTTP error: %v\n", err)
		return 1
	}
	setCORSHeaders(h)
	
	// Read request body
	body, err := io.ReadAll(h.Body())
	if err != nil {
		fmt.Printf("[ERROR] sendMessage read body error: %v\n", err)
		return handleHTTPError(h, err, 400)
	}
	
	// Parse message data
	var message ChatMessage
	err = json.Unmarshal(body, &message)
	if err != nil {
		fmt.Printf("[ERROR] sendMessage JSON parse error: %v\n", err)
		return handleHTTPError(h, err, 400)
	}
	
	// Set timestamp and generate ID if not provided
	if message.Time == 0 {
		message.Time = getCurrentTimestamp()
	}
	if message.ID == "" {
		message.ID = fmt.Sprintf("msg_%d_%s", message.Time, message.UserID)
	}
	
	// Save message to database
	err = saveChatMessage(message)
	if err != nil {
		fmt.Printf("[ERROR] sendMessage save error: %v\n", err)
		return handleHTTPError(h, err, 500)
	}
	
	// Publish message
	err = publishChatMessage(message)
	if err != nil {
		fmt.Printf("[ERROR] sendMessage publish error: %v\n", err)
		return handleHTTPError(h, err, 500)
	}
	
	// Return success response
	response := Response{
		Success: true,
		Data:    message,
	}
	
	fmt.Printf("[DEBUG] sendMessage completed successfully\n")
	return sendJSONResponse(h, response)
}

//export getMessages
func getMessages(e event.Event) uint32 {
	fmt.Printf("[DEBUG] getMessages called\n")
	
	h, err := e.HTTP()
	if err != nil {
		fmt.Printf("[ERROR] getMessages HTTP error: %v\n", err)
		return 1
	}
	setCORSHeaders(h)
	
	// Get all messages from database
	messages, err := getChatMessages()
	if err != nil {
		fmt.Printf("[ERROR] getMessages database error: %v\n", err)
		return handleHTTPError(h, err, 500)
	}
	
	// Return response
	response := Response{
		Success: true,
		Data:    messages,
	}
	
	fmt.Printf("[DEBUG] getMessages completed successfully, returning %d messages\n", len(messages))
	return sendJSONResponse(h, response)
}

//export getWebSocketURL
func getWebSocketURL(e event.Event) uint32 {
	fmt.Printf("[DEBUG] getWebSocketURL called\n")
	
	h, err := e.HTTP()
	if err != nil {
		fmt.Printf("[ERROR] getWebSocketURL HTTP error: %v\n", err)
		return 1
	}
	setCORSHeaders(h)
	
	// Get channel type from query
	channelType := getQueryParam(h, "type", "canvas")
	
	var channelName string
	switch channelType {
	case "canvas":
		channelName = "canvas"
	case "chat":
		channelName = "chat"
	default:
		fmt.Printf("[ERROR] getWebSocketURL invalid channel type: %s\n", channelType)
		h.Write([]byte("Invalid channel type. Use 'canvas' or 'chat'"))
		h.Return(400)
		return 1
	}
	
	// Create channel and get WebSocket URL
	channel, err := pubsub.Channel(channelName)
	if err != nil {
		fmt.Printf("[ERROR] getWebSocketURL channel creation error: %v\n", err)
		return handleHTTPError(h, err, 500)
	}
	
	// Subscribe to channel
	err = channel.Subscribe()
	if err != nil {
		fmt.Printf("[ERROR] getWebSocketURL subscription error: %v\n", err)
		return handleHTTPError(h, err, 500)
	}
	
	// Get WebSocket URL
	url, err := channel.WebSocket().Url()
	if err != nil {
		fmt.Printf("[ERROR] getWebSocketURL URL error: %v\n", err)
		return handleHTTPError(h, err, 500)
	}
	
	// Return WebSocket URL
	response := Response{
		Success: true,
		Data: map[string]string{
			"websocketUrl": url.Path,
			"channel":      channelName,
		},
	}
	
	fmt.Printf("[DEBUG] getWebSocketURL completed successfully\n")
	return sendJSONResponse(h, response)
}
