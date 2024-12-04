package gemini

import (
	"github.com/google/generative-ai-go/genai"
)

type boxed[T any] struct {
	value T
}

// GenerativeModelBuilder implements the builder pattern for [genai.GenerativeModel].
type GenerativeModelBuilder struct {
	copy *genai.GenerativeModel

	name              *boxed[string]
	generationConfig  *boxed[genai.GenerationConfig]
	safetySettings    *boxed[[]*genai.SafetySetting]
	tools             *boxed[[]*genai.Tool]
	toolConfig        *boxed[*genai.ToolConfig]
	systemInstruction *boxed[*genai.Content]
	cachedContentName *boxed[string]
}

// NewGenerativeModelBuilder returns a new [GenerativeModelBuilder] with empty default values.
func NewGenerativeModelBuilder() *GenerativeModelBuilder {
	return &GenerativeModelBuilder{}
}

// newCopyGenerativeModelBuilder creates a new [GenerativeModelBuilder],
// taking the default values from an existing [genai.GenerativeModel] object.
func newCopyGenerativeModelBuilder(copy *genai.GenerativeModel) *GenerativeModelBuilder {
	return &GenerativeModelBuilder{copy: copy}
}

// WithName sets the model name.
func (b *GenerativeModelBuilder) WithName(
	modelName string,
) *GenerativeModelBuilder {
	b.name = &boxed[string]{modelName}
	return b
}

// WithGenerationConfig sets the generation config.
func (b *GenerativeModelBuilder) WithGenerationConfig(
	generationConfig genai.GenerationConfig,
) *GenerativeModelBuilder {
	b.generationConfig = &boxed[genai.GenerationConfig]{generationConfig}
	return b
}

// WithSafetySettings sets the safety settings.
func (b *GenerativeModelBuilder) WithSafetySettings(
	safetySettings []*genai.SafetySetting,
) *GenerativeModelBuilder {
	b.safetySettings = &boxed[[]*genai.SafetySetting]{safetySettings}
	return b
}

// WithTools sets the tools.
func (b *GenerativeModelBuilder) WithTools(
	tools []*genai.Tool,
) *GenerativeModelBuilder {
	b.tools = &boxed[[]*genai.Tool]{tools}
	return b
}

// WithToolConfig sets the tool config.
func (b *GenerativeModelBuilder) WithToolConfig(
	toolConfig *genai.ToolConfig,
) *GenerativeModelBuilder {
	b.toolConfig = &boxed[*genai.ToolConfig]{toolConfig}
	return b
}

// WithSystemInstruction sets the system instruction.
func (b *GenerativeModelBuilder) WithSystemInstruction(
	systemInstruction *genai.Content,
) *GenerativeModelBuilder {
	b.systemInstruction = &boxed[*genai.Content]{systemInstruction}
	return b
}

// WithCachedContentName sets the name of the [genai.CachedContent] to use.
func (b *GenerativeModelBuilder) WithCachedContentName(
	cachedContentName string,
) *GenerativeModelBuilder {
	b.cachedContentName = &boxed[string]{cachedContentName}
	return b
}

// build builds and returns a new [genai.GenerativeModel] using the given [genai.Client].
// It will panic if the copy and the model name are not set.
func (b *GenerativeModelBuilder) build(client *genai.Client) *genai.GenerativeModel {
	if b.copy == nil && b.name == nil {
		panic("model name is required")
	}

	model := b.copy
	if b.name != nil {
		model = client.GenerativeModel(b.name.value)
		if b.copy != nil {
			model.GenerationConfig = b.copy.GenerationConfig
			model.SafetySettings = b.copy.SafetySettings
			model.Tools = b.copy.Tools
			model.ToolConfig = b.copy.ToolConfig
			model.SystemInstruction = b.copy.SystemInstruction
			model.CachedContentName = b.copy.CachedContentName
		}
	}
	b.configure(model)
	return model
}

// configure configures the given generative model using the builder values.
func (b *GenerativeModelBuilder) configure(model *genai.GenerativeModel) {
	if b.generationConfig != nil {
		model.GenerationConfig = b.generationConfig.value
	}
	if b.safetySettings != nil {
		model.SafetySettings = b.safetySettings.value
	}
	if b.tools != nil {
		model.Tools = b.tools.value
	}
	if b.toolConfig != nil {
		model.ToolConfig = b.toolConfig.value
	}
	if b.systemInstruction != nil {
		model.SystemInstruction = b.systemInstruction.value
	}
	if b.cachedContentName != nil {
		model.CachedContentName = b.cachedContentName.value
	}
}
