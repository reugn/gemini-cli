package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"sync"

	"google.golang.org/genai"
)

const DefaultModel = "gemini-2.5-flash"

// ChatSession represents a gemini powered chat session.
type ChatSession struct {
	ctx context.Context

	client *genai.Client
	chat   *genai.Chat
	config *genai.GenerateContentConfig
	model  string

	loadModels sync.Once
	models     []string
}

// NewChatSession returns a new [ChatSession].
func NewChatSession(ctx context.Context, model string,
	contentConfig *genai.GenerateContentConfig) (*ChatSession, error) {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	chat, err := client.Chats.Create(ctx, model, contentConfig, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %w", err)
	}

	return &ChatSession{
		ctx:    ctx,
		client: client,
		chat:   chat,
		config: contentConfig,
		model:  model,
	}, nil
}

// SendMessage sends a request to the model as part of a chat session.
func (c *ChatSession) SendMessage(input string) (*genai.GenerateContentResponse, error) {
	return c.chat.SendMessage(c.ctx, genai.Part{Text: input})
}

// SendMessageStream is like SendMessage, but with a streaming request.
func (c *ChatSession) SendMessageStream(input string) iter.Seq2[*genai.GenerateContentResponse, error] {
	return c.chat.SendMessageStream(c.ctx, genai.Part{Text: input})
}

// ModelInfo returns information about the chat generative model in JSON format.
func (c *ChatSession) ModelInfo() (string, error) {
	modelInfo, err := c.client.Models.Get(c.ctx, c.model, nil)
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
		for model, err := range c.client.Models.All(c.ctx) {
			if err != nil {
				continue
			}
			c.models = append(c.models, model.Name)
		}
	})
	return c.models
}

// SetModel sets the chat generative model.
func (c *ChatSession) SetModel(model string) error {
	chat, err := c.client.Chats.Create(c.ctx, model, c.config, c.GetHistory())
	if err != nil {
		return fmt.Errorf("failed to set model: %w", err)
	}

	c.model = model
	c.chat = chat
	return nil
}

// GetHistory returns the chat session history.
func (c *ChatSession) GetHistory() []*genai.Content {
	return c.chat.History(true)
}

// SetHistory sets the chat session history.
func (c *ChatSession) SetHistory(history []*genai.Content) error {
	chat, err := c.client.Chats.Create(c.ctx, c.model, c.config, history)
	if err != nil {
		return fmt.Errorf("failed to set history: %w", err)
	}

	c.chat = chat
	return nil
}

// ClearHistory clears the chat session history.
func (c *ChatSession) ClearHistory() error {
	return c.SetHistory(nil)
}

// SetSystemInstruction sets the chat session system instruction.
func (c *ChatSession) SetSystemInstruction(systemInstruction *genai.Content) error {
	c.config.SystemInstruction = systemInstruction
	chat, err := c.client.Chats.Create(c.ctx, c.model, c.config, c.GetHistory())
	if err != nil {
		return fmt.Errorf("failed to set system instruction: %w", err)
	}

	c.chat = chat
	return nil
}
