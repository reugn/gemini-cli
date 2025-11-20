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

const (
	toolGoogleSearch = "GOOGLE_SEARCH"
	toolURLContext   = "URL_CONTEXT"
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
// and uses the custom Threshold for serialization.
type SafetySetting struct {
	Category  genai.HarmCategory `json:"category"`
	Threshold Threshold          `json:"threshold"`
}

// Tool represents a model tool configuration.
type Tool struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

// ApplicationData encapsulates application state and configuration.
// Note that the chat history is stored in plain text format.
type ApplicationData struct {
	SystemPrompts  map[string]gemini.SystemInstruction      `json:"system_prompts"`
	SafetySettings []SafetySetting                          `json:"safety_settings"`
	Tools          []Tool                                   `json:"tools"`
	History        map[string][]*gemini.SerializableContent `json:"history"`
}

// newDefaultApplicationData returns a new ApplicationData with default values.
func newDefaultApplicationData() *ApplicationData {
	defaultSafetySettings := []SafetySetting{
		{Category: genai.HarmCategoryHarassment, Threshold: thresholdLow},
		{Category: genai.HarmCategoryHateSpeech, Threshold: thresholdLow},
		{Category: genai.HarmCategorySexuallyExplicit, Threshold: thresholdLow},
		{Category: genai.HarmCategoryDangerousContent, Threshold: thresholdLow},
	}

	defaultTools := []Tool{
		{Name: toolGoogleSearch, Enabled: true},
		{Name: toolURLContext, Enabled: true},
	}

	return &ApplicationData{
		SystemPrompts:  make(map[string]gemini.SystemInstruction),
		SafetySettings: defaultSafetySettings,
		Tools:          defaultTools,
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

// GenaiTools builds a genai Tool slice using enabled entries.
func (d *ApplicationData) GenaiTools() []*genai.Tool {
	tools := make([]*genai.Tool, 0, len(d.Tools))
	for _, tool := range d.Tools {
		if !tool.Enabled {
			continue
		}

		var genaiTool *genai.Tool
		switch tool.Name {
		case toolGoogleSearch:
			genaiTool = &genai.Tool{GoogleSearch: &genai.GoogleSearch{}}
		case toolURLContext:
			genaiTool = &genai.Tool{URLContext: &genai.URLContext{}}
		default:
			continue // Skip unknown tools
		}

		tools = append(tools, genaiTool)
	}

	return tools
}

// GenaiContentConfig builds a genai GenerateContentConfig with the current
// safety settings and enabled tools.
func (d *ApplicationData) GenaiContentConfig() *genai.GenerateContentConfig {
	return &genai.GenerateContentConfig{
		SafetySettings: d.GenaiSafetySettings(),
		Tools:          d.GenaiTools(),
	}
}
