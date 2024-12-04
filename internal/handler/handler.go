package handler

// MessageHandler handles chat messages from the user.
type MessageHandler interface {
	// Handle processes the message and returns a response, along with a flag
	// indicating whether the application should terminate.
	Handle(message string) (Response, bool)
}
