package config

import (
	"github.com/google/generative-ai-go/genai"
	"github.com/reugn/gemini-cli/gemini"
)

// ApplicationData encapsulates application state and configuration.
// Note that the chat history is stored in plain text format.
type ApplicationData struct {
	SystemPrompts  map[string]gemini.SystemInstruction
	SafetySettings []*genai.SafetySetting
	History        map[string][]*gemini.SerializableContent
}

// newDefaultApplicationData returns a new ApplicationData with default values.
func newDefaultApplicationData() *ApplicationData {
	defaultSafetySettings := []*genai.SafetySetting{
		{Category: genai.HarmCategoryHarassment, Threshold: genai.HarmBlockLowAndAbove},
		{Category: genai.HarmCategoryHateSpeech, Threshold: genai.HarmBlockLowAndAbove},
		{Category: genai.HarmCategorySexuallyExplicit, Threshold: genai.HarmBlockLowAndAbove},
		{Category: genai.HarmCategoryDangerousContent, Threshold: genai.HarmBlockLowAndAbove},
	}
	return &ApplicationData{
		SystemPrompts:  make(map[string]gemini.SystemInstruction),
		SafetySettings: defaultSafetySettings,
		History:        make(map[string][]*gemini.SerializableContent),
	}
}

// AddHistoryRecord adds a history record to the application data.
func (d *ApplicationData) AddHistoryRecord(label string, content []*genai.Content) {
	serializableContent := make([]*gemini.SerializableContent, len(content))
	for i, c := range content {
		serializableContent[i] = gemini.NewSerializableContent(c)
	}
	d.History[label] = serializableContent
}
