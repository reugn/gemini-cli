package handler

import (
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/styles"
)

// RendererOptions represents configuration options for the terminal renderer.
type RendererOptions struct {
	StylePath string
	WordWrap  int
}

func (o RendererOptions) newTermRenderer() (*glamour.TermRenderer, error) {
	var styleOption glamour.TermRendererOption
	switch {
	case o.StylePath == styles.AutoStyle && os.Getenv("GLAMOUR_STYLE") != "":
		styleOption = glamour.WithEnvironmentConfig()
	default:
		styleOption = glamour.WithStylePath(o.StylePath)
	}

	return glamour.NewTermRenderer(
		styleOption,
		glamour.WithWordWrap(o.WordWrap),
	)
}
