package gemini

import (
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

// SerializableContent is the data type containing multipart text message content.
// It is a serializable equivalent of [genai.Content], where message content parts
// are represented as strings.
type SerializableContent struct {
	// Ordered parts that constitute a single message.
	Parts []string
	// The producer of the content. Must be either 'user' or 'model'.
	Role string
}

// NewSerializableContent instantiates and returns a new SerializableContent from
// the given [genai.Content].
// It will panic if the content type is not supported.
func NewSerializableContent(c *genai.Content) *SerializableContent {
	parts := make([]string, len(c.Parts))
	for i, part := range c.Parts {
		parts[i] = partToString(part)
	}
	return &SerializableContent{
		Parts: parts,
		Role:  c.Role,
	}
}

// ToContent converts the SerializableContent into a [genai.Content].
func (c *SerializableContent) ToContent() *genai.Content {
	parts := make([]genai.Part, len(c.Parts))
	for i, part := range c.Parts {
		parts[i] = genai.Text(part)
	}
	return &genai.Content{
		Parts: parts,
		Role:  c.Role,
	}
}

func partToString(part genai.Part) string {
	switch p := part.(type) {
	case genai.Text:
		return string(p)
	default:
		panic(fmt.Errorf("unsupported part type: %T", part))
	}
}
