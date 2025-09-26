package lib

// Pixel represents a single pixel on the canvas
type Pixel struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Color  string `json:"color"`
	UserID string `json:"userId"`
	Time   int64  `json:"time"`
}

// ChatMessage represents a chat message
type ChatMessage struct {
	ID      string `json:"id"`
	UserID  string `json:"userId"`
	Message string `json:"message"`
	Time    int64  `json:"time"`
}

// CanvasUpdate represents a batch of pixel updates
type CanvasUpdate struct {
	Pixels []Pixel `json:"pixels"`
	Time   int64   `json:"time"`
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
