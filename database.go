package lib

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/taubyte/go-sdk/database"
)

var (
	canvasDB database.Database
	chatDB   database.Database
	dbMutex  sync.RWMutex
	dbInit   bool
)

// initDatabases initializes database connections
func initDatabases() uint32 {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if dbInit {
		return 0 // Already initialized
	}

	fmt.Printf("[DEBUG] Initializing databases\n")
	
	var err error
	
	// Initialize canvas database
	canvasDB, err = database.New("/canvas")
	if err != nil {
		fmt.Printf("[ERROR] Failed to initialize canvas database: %v\n", err)
		return 1
	}
	
	// Initialize chat database
	chatDB, err = database.New("/chat")
	if err != nil {
		fmt.Printf("[ERROR] Failed to initialize chat database: %v\n", err)
		return 1
	}
	
	dbInit = true
	fmt.Printf("[DEBUG] Databases initialized successfully\n")
	return 0
}

// getCanvasDB returns the canvas database connection
func getCanvasDB() (database.Database, uint32) {
	if !dbInit {
		if initDatabases() != 0 {
			var emptyDB database.Database
			return emptyDB, 1
		}
	}
	return canvasDB, 0
}

// getChatDB returns the chat database connection
func getChatDB() (database.Database, uint32) {
	if !dbInit {
		if initDatabases() != 0 {
			var emptyDB database.Database
			return emptyDB, 1
		}
	}
	return chatDB, 0
}

// savePixel saves a pixel to the canvas database
func savePixel(pixel Pixel) error {
	fmt.Printf("[DEBUG] Saving pixel: x=%d, y=%d, color=%s\n", pixel.X, pixel.Y, pixel.Color)
	
	db, dbErr := getCanvasDB()
	if dbErr != 0 {
		return fmt.Errorf("database connection failed")
	}
	
	// Create key from coordinates
	key := fmt.Sprintf("pixel_%d_%d", pixel.X, pixel.Y)
	
	// Convert pixel to JSON
	pixelData, err := json.Marshal(pixel)
	if err != nil {
		return err
	}
	
	// Save to database
	err = db.Put(key, pixelData)
	if err != nil {
		fmt.Printf("[ERROR] Failed to save pixel: %v\n", err)
		return err
	}
	
	fmt.Printf("[DEBUG] Pixel saved successfully\n")
	return nil
}

// getPixelFromDatabase retrieves a pixel from the canvas database
func getPixelFromDatabase(x, y int) (*Pixel, error) {
	fmt.Printf("[DEBUG] Getting pixel: x=%d, y=%d\n", x, y)
	
	db, dbErr := getCanvasDB()
	if dbErr != 0 {
		return nil, fmt.Errorf("database connection failed")
	}
	
	// Create key from coordinates
	key := fmt.Sprintf("pixel_%d_%d", x, y)
	
	// Get from database
	pixelData, err := db.Get(key)
	if err != nil {
		// Pixel doesn't exist, return nil
		return nil, nil
	}
	
	// Parse JSON
	var pixel Pixel
	err = json.Unmarshal(pixelData, &pixel)
	if err != nil {
		return nil, err
	}
	
	fmt.Printf("[DEBUG] Pixel retrieved successfully\n")
	return &pixel, nil
}

// saveChatMessage saves a chat message to the chat database
func saveChatMessage(message ChatMessage) error {
	fmt.Printf("[DEBUG] Saving chat message: id=%s, user=%s\n", message.ID, message.UserID)
	
	db, dbErr := getChatDB()
	if dbErr != 0 {
		return fmt.Errorf("database connection failed")
	}
	
	// Convert message to JSON
	messageData, err := json.Marshal(message)
	if err != nil {
		return err
	}
	
	// Save to database
	err = db.Put(message.ID, messageData)
	if err != nil {
		fmt.Printf("[ERROR] Failed to save chat message: %v\n", err)
		return err
	}
	
	fmt.Printf("[DEBUG] Chat message saved successfully\n")
	return nil
}

// getChatMessages retrieves all chat messages
func getChatMessages() ([]ChatMessage, error) {
	fmt.Printf("[DEBUG] Getting chat messages\n")
	
	db, dbErr := getChatDB()
	if dbErr != 0 {
		return nil, fmt.Errorf("database connection failed")
	}
	
	// List all keys in chat database
	keys, err := db.List("")
	if err != nil {
		return nil, err
	}
	
	var messages []ChatMessage
	
	// Get each message
	for _, key := range keys {
		messageData, err := db.Get(key)
		if err != nil {
			continue // Skip invalid messages
		}
		
		var message ChatMessage
		err = json.Unmarshal(messageData, &message)
		if err != nil {
			continue // Skip invalid messages
		}
		
		messages = append(messages, message)
	}
	
	fmt.Printf("[DEBUG] Retrieved %d chat messages\n", len(messages))
	return messages, nil
}
