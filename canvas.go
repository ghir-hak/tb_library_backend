package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/taubyte/go-sdk/event"
)

//export drawPixel
func drawPixel(e event.Event) uint32 {
	fmt.Printf("[DEBUG] drawPixel called\n")
	
	h, err := e.HTTP()
	if err != nil {
		fmt.Printf("[ERROR] drawPixel HTTP error: %v\n", err)
		return 1
	}
	setCORSHeaders(h)
	
	// Read request body
	body, err := io.ReadAll(h.Body())
	if err != nil {
		fmt.Printf("[ERROR] drawPixel read body error: %v\n", err)
		return handleHTTPError(h, err, 400)
	}
	
	// Parse pixel data
	var pixel Pixel
	err = json.Unmarshal(body, &pixel)
	if err != nil {
		fmt.Printf("[ERROR] drawPixel JSON parse error: %v\n", err)
		return handleHTTPError(h, err, 400)
	}
	
	// Set timestamp
	pixel.Time = getCurrentTimestamp()
	
	// Save pixel to database
	err = savePixel(pixel)
	if err != nil {
		fmt.Printf("[ERROR] drawPixel save error: %v\n", err)
		return handleHTTPError(h, err, 500)
	}
	
	// Publish pixel update
	err = publishPixelUpdate(pixel)
	if err != nil {
		fmt.Printf("[ERROR] drawPixel publish error: %v\n", err)
		return handleHTTPError(h, err, 500)
	}
	
	// Return success response
	response := Response{
		Success: true,
		Data:    pixel,
	}
	
	fmt.Printf("[DEBUG] drawPixel completed successfully\n")
	return sendJSONResponse(h, response)
}

//export getPixel
func getPixel(e event.Event) uint32 {
	fmt.Printf("[DEBUG] getPixel called\n")
	
	h, err := e.HTTP()
	if err != nil {
		fmt.Printf("[ERROR] getPixel HTTP error: %v\n", err)
		return 1
	}
	setCORSHeaders(h)
	
	// Get x coordinate
	xStr, retCode := getRequiredQueryParam(h, "x")
	if retCode != 0 {
		return retCode
	}
	
	// Get y coordinate
	yStr, retCode := getRequiredQueryParam(h, "y")
	if retCode != 0 {
		return retCode
	}
	
	// Convert to integers
	x, err := strconv.Atoi(xStr)
	if err != nil {
		fmt.Printf("[ERROR] getPixel invalid x coordinate: %v\n", err)
		return handleHTTPError(h, err, 400)
	}
	
	y, err := strconv.Atoi(yStr)
	if err != nil {
		fmt.Printf("[ERROR] getPixel invalid y coordinate: %v\n", err)
		return handleHTTPError(h, err, 400)
	}
	
	// Get pixel from database
	pixel, err := getPixelFromDatabase(x, y)
	if err != nil {
		fmt.Printf("[ERROR] getPixel database error: %v\n", err)
		return handleHTTPError(h, err, 500)
	}
	
	// Return response
	response := Response{
		Success: true,
		Data:    pixel,
	}
	
	fmt.Printf("[DEBUG] getPixel completed successfully\n")
	return sendJSONResponse(h, response)
}

//export getCanvas
func getCanvas(e event.Event) uint32 {
	fmt.Printf("[DEBUG] getCanvas called\n")
	
	h, err := e.HTTP()
	if err != nil {
		fmt.Printf("[ERROR] getCanvas HTTP error: %v\n", err)
		return 1
	}
	setCORSHeaders(h)
	
	// Get all pixels from database
	pixels, err := getAllPixels()
	if err != nil {
		fmt.Printf("[ERROR] getCanvas database error: %v\n", err)
		return handleHTTPError(h, err, 500)
	}
	
	// Return response
	response := Response{
		Success: true,
		Data:    pixels,
	}
	
	fmt.Printf("[DEBUG] getCanvas completed successfully, returning %d pixels\n", len(pixels))
	return sendJSONResponse(h, response)
}

// getPixelFromDB retrieves a pixel from the database
func getPixelFromDB(x, y int) (*Pixel, error) {
	return getPixelFromDatabase(x, y)
}

// getAllPixels retrieves all pixels from the canvas database
func getAllPixels() ([]Pixel, error) {
	fmt.Printf("[DEBUG] Getting all pixels\n")
	
	db, dbErr := getCanvasDB()
	if dbErr != 0 {
		return nil, fmt.Errorf("database connection failed")
	}
	
	// List all keys in canvas database
	keys, err := db.List("")
	if err != nil {
		return nil, err
	}
	
	var pixels []Pixel
	
	// Get each pixel
	for _, key := range keys {
		pixelData, err := db.Get(key)
		if err != nil {
			continue // Skip invalid pixels
		}
		
		var pixel Pixel
		err = json.Unmarshal(pixelData, &pixel)
		if err != nil {
			continue // Skip invalid pixels
		}
		
		pixels = append(pixels, pixel)
	}
	
	fmt.Printf("[DEBUG] Retrieved %d pixels\n", len(pixels))
	return pixels, nil
}
