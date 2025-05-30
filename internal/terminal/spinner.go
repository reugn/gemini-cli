package terminal

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

const (
	moveCursorBackward  = "\033[%dD"
	clearLineFromCursor = "\033[K"
	progressRune        = '.'
)

// Spinner is a visual indicator of progress displayed in the terminal as a
// scrolling dot animation.
type Spinner struct {
	writer   *bufio.Writer
	interval time.Duration
	signal   chan struct{}

	maxLength int
	length    int
}

// NewSpinner returns a new Spinner.
func NewSpinner(w io.Writer, interval time.Duration, length int) *Spinner {
	return &Spinner{
		writer:    bufio.NewWriter(w),
		interval:  interval,
		signal:    make(chan struct{}),
		maxLength: length,
	}
}

//nolint:errcheck
func (s *Spinner) Start() {
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		s.length = 0
		for {
			select {
			case <-s.signal:
				if s.length > 0 {
					s.Clear()
				}
				s.signal <- struct{}{}
				return
			case <-ticker.C:
				if s.length < s.maxLength {
					s.writer.WriteRune(progressRune)
					s.writer.Flush()
					s.length++
				} else {
					s.Clear()
					s.length = 0
				}
			}
		}
	}()
}

//nolint:errcheck,staticcheck
func (s *Spinner) Clear() {
	s.writer.WriteString(fmt.Sprintf(moveCursorBackward, s.length))
	s.writer.WriteString(clearLineFromCursor)
	s.writer.Flush()
}

func (s *Spinner) Stop() {
	s.signal <- struct{}{}
	<-s.signal
}
