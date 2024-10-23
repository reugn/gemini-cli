package cli

import (
	"slices"

	"github.com/manifoldco/promptui"
)

var (
	inputMode = []string{"single-line", "multi-line"}
)

// selectModel returns the selected generative model name.
func selectModel(current string, models []string) (string, error) {
	prompt := promptui.Select{
		Label:        "Select generative model",
		HideSelected: true,
		Items:        models,
		CursorPos:    slices.Index(models, current),
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

// selectInputMode returns true if multiline input is selected;
// otherwise, it returns false.
func selectInputMode(multiline bool) (bool, error) {
	var cursorPos int
	if multiline {
		cursorPos = 1
	}

	prompt := promptui.Select{
		Label:        "Select input mode",
		HideSelected: true,
		Items:        inputMode,
		CursorPos:    cursorPos,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	return result == inputMode[1], nil
}
