package config

import (
	"github.com/reugn/gemini-cli/gemini"
	"google.golang.org/genai"
)

const (
	thresholdLow    = "LOW"
	thresholdMedium = "MEDIUM"
	thresholdHigh   = "HIGH"
	thresholdOff    = "OFF"
)

// Threshold is a custom type that wraps genai.HarmBlockThreshold
// and uses the custom string for serialization.
type Threshold string

func (t Threshold) toGenai() genai.HarmBlockThreshold {
	switch t {
	case thresholdLow:
		return genai.HarmBlockThresholdBlockLowAndAbove
	case thresholdMedium:
		return genai.HarmBlockThresholdBlockMediumAndAbove
	case thresholdHigh:
		return genai.HarmBlockThresholdBlockOnlyHigh
	case thresholdOff:
		return genai.HarmBlockThresholdOff
	default:
		return genai.HarmBlockThresholdUnspecified
	}
}

// SafetySetting is a custom type that wraps genai.SafetySetting
// and uses the custom HarmBlockThreshold for serialization.
type SafetySetting struct {
	Category  genai.HarmCategory `json:"category"`
	Threshold Threshold          `json:"threshold"`
}

// ApplicationData encapsulates application state and configuration.
// Note that the chat history is stored in plain text format.
type ApplicationData struct {
	SystemPrompts  map[string]gemini.SystemInstruction      `json:"system_prompts"`
	SafetySettings []*SafetySetting                         `json:"safety_settings"`
	History        map[string][]*gemini.SerializableContent `json:"history"`
}

// newDefaultApplicationData returns a new ApplicationData with default values.
func newDefaultApplicationData() *ApplicationData {
	defaultSafetySettings := []*SafetySetting{
		{Category: genai.HarmCategoryHarassment, Threshold: thresholdLow},
		{Category: genai.HarmCategoryHateSpeech, Threshold: thresholdLow},
		{Category: genai.HarmCategorySexuallyExplicit, Threshold: thresholdLow},
		{Category: genai.HarmCategoryDangerousContent, Threshold: thresholdLow},
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

// GenaiSafetySettings converts the application data safety settings to genai safety settings.
func (d *ApplicationData) GenaiSafetySettings() []*genai.SafetySetting {
	genaiSafetySettings := make([]*genai.SafetySetting, len(d.SafetySettings))
	for i, s := range d.SafetySettings {
		genaiSafetySettings[i] = &genai.SafetySetting{
			Category:  s.Category,
			Threshold: s.Threshold.toGenai(),
		}
	}

	return genaiSafetySettings
}
