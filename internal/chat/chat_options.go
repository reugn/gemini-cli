package chat

import "github.com/reugn/gemini-cli/internal/handler"

// Opts represents the Chat configuration options.
type Opts struct {
	GenerativeModel string
	Multiline       bool
	LineTerminator  string
	StylePath       string
	WordWrap        int
}

func (o *Opts) rendererOptions() handler.RendererOptions {
	return handler.RendererOptions{
		StylePath: o.StylePath,
		WordWrap:  o.WordWrap,
	}
}
