package gemini

import (
	"context"
	"sync"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const DefaultModel = "gemini-pro"

// ChatSession represents a gemini powered chat session.
type ChatSession struct {
	ctx context.Context

	client  *genai.Client
	session *genai.ChatSession

	loadModels sync.Once
	models     []string
}

// NewChatSession returns a new [ChatSession].
func NewChatSession(ctx context.Context, model, apiKey string) (*ChatSession, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	return &ChatSession{
		ctx:     ctx,
		client:  client,
		session: client.GenerativeModel(model).StartChat(),
	}, nil
}

// SendMessage sends a request to the model as part of a chat session.
func (c *ChatSession) SendMessage(input string) (*genai.GenerateContentResponse, error) {
	return c.session.SendMessage(c.ctx, genai.Text(input))
}

// SendMessageStream is like SendMessage, but with a streaming request.
func (c *ChatSession) SendMessageStream(input string) *genai.GenerateContentResponseIterator {
	return c.session.SendMessageStream(c.ctx, genai.Text(input))
}

// SetGenerativeModel sets the name of the generative model for the chat.
// It preserves the history from the previous chat session.
func (c *ChatSession) SetGenerativeModel(model string) {
	history := c.session.History
	c.session = c.client.GenerativeModel(model).StartChat()
	c.session.History = history
}

// ListModels returns a list of the supported generative model names.
func (c *ChatSession) ListModels() []string {
	c.loadModels.Do(func() {
		c.models = []string{DefaultModel}
		iter := c.client.ListModels(c.ctx)
		for {
			modelInfo, err := iter.Next()
			if err != nil {
				break
			}
			c.models = append(c.models, modelInfo.Name)
		}
	})
	return c.models
}

// ClearHistory clears chat history.
func (c *ChatSession) ClearHistory() {
	c.session.History = make([]*genai.Content, 0)
}

// Close closes the chat session.
func (c *ChatSession) Close() error {
	return c.client.Close()
}
