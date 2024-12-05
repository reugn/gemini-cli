package handler

// QuitCommand processes the chat quit system command.
// It implements the MessageHandler interface.
type QuitCommand struct {
	*IO
}

var _ MessageHandler = (*QuitCommand)(nil)

// NewQuitCommand returns a new QuitCommand.
func NewQuitCommand(io *IO) *QuitCommand {
	return &QuitCommand{IO: io}
}

// Handle processes the chat quit command.
func (h *QuitCommand) Handle(_ string) (Response, bool) {
	return dataResponse("Exiting gemini-cli..."), true
}
