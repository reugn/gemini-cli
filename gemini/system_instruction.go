package gemini

import "github.com/google/generative-ai-go/genai"

// SystemInstruction represents a serializable system prompt, a more forceful
// instruction to the language model. The model will prioritize adhering to
// system instructions over regular prompts.
type SystemInstruction string

// ToContent converts the SystemInstruction to [genai.Content].
func (si SystemInstruction) ToContent() *genai.Content {
	return genai.NewUserContent(genai.Text(si))
}
