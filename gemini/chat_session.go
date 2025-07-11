package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const DefaultModel = "gemini-2.5-flash"

// ChatSession represents a gemini powered chat session.
type ChatSession struct {
	ctx context.Context

	client  *genai.Client
	model   *genai.GenerativeModel
	session *genai.ChatSession

	loadModels sync.Once
	models     []string
}

// NewChatSession returns a new [ChatSession].
func NewChatSession(
	ctx context.Context, modelBuilder *GenerativeModelBuilder, apiKey string,
) (*ChatSession, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	generativeModel := modelBuilder.build(client)
	return &ChatSession{
		ctx:     ctx,
		client:  client,
		model:   generativeModel,
		session: generativeModel.StartChat(),
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

// SetModel sets a new generative model configured with the builder and starts
// a new chat session. It preserves the history of the previous chat session.
func (c *ChatSession) SetModel(modelBuilder *GenerativeModelBuilder) {
	history := c.session.History
	c.model = modelBuilder.build(c.client)
	c.session = c.model.StartChat()
	c.session.History = history
}

// CopyModelBuilder returns a copy builder for the chat generative model.
func (c *ChatSession) CopyModelBuilder() *GenerativeModelBuilder {
	return newCopyGenerativeModelBuilder(c.model)
}

// ModelInfo returns information about the chat generative model in JSON format.
func (c *ChatSession) ModelInfo() (string, error) {
	modelInfo, err := c.model.Info(c.ctx)
	if err != nil {
		return "", err
	}
	encoded, err := json.MarshalIndent(modelInfo, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error encoding model info: %w", err)
	}
	return string(encoded), nil
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

// GetHistory returns the chat session history.
func (c *ChatSession) GetHistory() []*genai.Content {
	return c.session.History
}

// SetHistory sets the chat session history.
func (c *ChatSession) SetHistory(content []*genai.Content) {
	c.session.History = content
}

// ClearHistory clears the chat session history.
func (c *ChatSession) ClearHistory() {
	c.session.History = make([]*genai.Content, 0)
}

// Close closes the chat session.
func (c *ChatSession) Close() error {
	return c.client.Close()
}
