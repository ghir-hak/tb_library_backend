package lib

import (
	"encoding/json"
	"fmt"

	"github.com/taubyte/go-sdk/event"
	pubsub "github.com/taubyte/go-sdk/pubsub/node"
)

// publishPixelUpdate publishes a pixel update to the canvas channel
func publishPixelUpdate(pixel Pixel) error {
	fmt.Printf("[DEBUG] Publishing pixel update: x=%d, y=%d, color=%s\n", pixel.X, pixel.Y, pixel.Color)
	
	// Create canvas channel
	channel, err := pubsub.Channel("canvas")
	if err != nil {
		fmt.Printf("[ERROR] Failed to create canvas channel: %v\n", err)
		return err
	}
	
	// Convert pixel to JSON
	pixelData, err := json.Marshal(pixel)
	if err != nil {
		return err
	}
	
	// Publish to channel
	err = channel.Publish(pixelData)
	if err != nil {
		fmt.Printf("[ERROR] Failed to publish pixel update: %v\n", err)
		return err
	}
	
	fmt.Printf("[DEBUG] Pixel update published successfully\n")
	return nil
}

// publishChatMessage publishes a chat message to the chat channel
func publishChatMessage(message ChatMessage) error {
	fmt.Printf("[DEBUG] Publishing chat message: id=%s, user=%s\n", message.ID, message.UserID)
	
	// Create chat channel
	channel, err := pubsub.Channel("chat")
	if err != nil {
		fmt.Printf("[ERROR] Failed to create chat channel: %v\n", err)
		return err
	}
	
	// Convert message to JSON
	messageData, err := json.Marshal(message)
	if err != nil {
		return err
	}
	
	// Publish to channel
	err = channel.Publish(messageData)
	if err != nil {
		fmt.Printf("[ERROR] Failed to publish chat message: %v\n", err)
		return err
	}
	
	fmt.Printf("[DEBUG] Chat message published successfully\n")
	return nil
}

//export handleCanvasEvent
func handleCanvasEvent(e event.Event) uint32 {
	fmt.Printf("[DEBUG] handleCanvasEvent called\n")
	
	// Handle incoming canvas events from pub/sub
	channel, err := e.PubSub()
	if err != nil {
		fmt.Printf("[ERROR] handleCanvasEvent pub/sub error: %v\n", err)
		return 1
	}
	
	// Get channel data
	data, err := channel.Data()
	if err != nil {
		fmt.Printf("[ERROR] handleCanvasEvent data error: %v\n", err)
		return 1
	}
	
	// Parse pixel data
	var pixel Pixel
	err = json.Unmarshal(data, &pixel)
	if err != nil {
		fmt.Printf("[ERROR] handleCanvasEvent JSON parse error: %v\n", err)
		return 1
	}
	
	// Save pixel to database
	err = savePixel(pixel)
	if err != nil {
		fmt.Printf("[ERROR] handleCanvasEvent save pixel error: %v\n", err)
		return 1
	}
	
	fmt.Printf("[DEBUG] handleCanvasEvent completed successfully\n")
	return 0
}

//export handleChatEvent
func handleChatEvent(e event.Event) uint32 {
	fmt.Printf("[DEBUG] handleChatEvent called\n")
	
	// Handle incoming chat events from pub/sub
	channel, err := e.PubSub()
	if err != nil {
		fmt.Printf("[ERROR] handleChatEvent pub/sub error: %v\n", err)
		return 1
	}
	
	// Get channel data
	data, err := channel.Data()
	if err != nil {
		fmt.Printf("[ERROR] handleChatEvent data error: %v\n", err)
		return 1
	}
	
	// Parse chat message data
	var message ChatMessage
	err = json.Unmarshal(data, &message)
	if err != nil {
		fmt.Printf("[ERROR] handleChatEvent JSON parse error: %v\n", err)
		return 1
	}
	
	// Save message to database
	err = saveChatMessage(message)
	if err != nil {
		fmt.Printf("[ERROR] handleChatEvent save message error: %v\n", err)
		return 1
	}
	
	fmt.Printf("[DEBUG] handleChatEvent completed successfully\n")
	return 0
}
