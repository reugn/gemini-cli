package cli

import (
	"bufio"
	"fmt"
	"time"
)

const (
	moveCursorBackward  = "\033[%dD"
	clearLineFromCursor = "\033[K"
	progressRune        = rune('.') // 46
)

type spinner struct {
	length   int
	interval time.Duration
	writer   *bufio.Writer
	signal   chan struct{}
}

func newSpinner(length int, interval time.Duration, writer *bufio.Writer) *spinner {
	return &spinner{
		length:   length,
		interval: interval,
		writer:   writer,
		signal:   make(chan struct{}),
	}
}

//nolint:errcheck
func (s *spinner) start() {
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		var n int
		for {
			select {
			case <-s.signal:
				if n > 0 {
					s.clear(n)
				}
				s.signal <- struct{}{}
				return
			case <-ticker.C:
				if n < s.length {
					s.writer.WriteRune(progressRune)
					s.writer.Flush()
					n++
				} else {
					s.clear(n)
					n = 0
				}
			}
		}
	}()
}

func (s *spinner) clear(n int) {
	s.writer.WriteString(fmt.Sprintf(moveCursorBackward, n))
	s.writer.WriteString(clearLineFromCursor)
	s.writer.Flush()
}

func (s *spinner) stop() {
	s.signal <- struct{}{}
	<-s.signal
}
