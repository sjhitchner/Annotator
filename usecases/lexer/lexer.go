package lexer

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	itemError itemType = iota
	itemEOF
	itemName
	itemText
	itemHyperlink
	itemLeftHyper
	itemRightHyper

	eof = -1
)

type itemType int

type item struct {
	typ   itemType
	value string
}

func (t item) String() string {
	switch t.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return t.value
	}
	if len(t.value) > 50 {
		return fmt.Sprintf("%d:%.50q...", t.typ, t.value)
	}
	return fmt.Sprintf("%d:%q", t.typ, t.value)
}

type stateFunction func(*lexer) stateFunction

// Lexer
type lexer struct {
	input string
	start int
	pos   int
	width int
	state stateFunction
	items chan item
}

// No copying, just a slice
// Channel usage adds some overhead could use ring buffer
func NewLexer(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan item, 2),
		state: lexName,
	}
	return l
}

func (t *lexer) NextItem() item {
	for {
		select {
		case item := <-t.items:
			return item
		default:
			t.state = t.state(t)
		}
	}
	panic("should never get here")
}

func (t *lexer) run() {
	for state := lexName; state != nil; {
		state = state(t)
	}
	close(t.items)
}

func (t *lexer) emit(i itemType) {
	t.items <- item{i, t.input[t.start:t.pos]}
	t.start = t.pos
}

func (t *lexer) next() rune {
	r, w := utf8.DecodeRuneInString(t.input[t.pos:])
	t.width = w
	t.pos += t.width

	if int(t.pos) >= len(t.input) {
		t.width = 0
		return eof
	}
	return r
}

func (t *lexer) errorf(format string, args ...interface{}) stateFunction {
	t.items <- item{itemError, fmt.Sprintf(format, args...)}
	return nil
}

const (
	LEFT_HYPER  = "<a"
	RIGHT_HYPER = "</a>"
)

//State Functions
//
// States
//  Alphanumeric String (possible Name)
//  Non Alphanumeric String (ignored text)
//  Text between <a..a> (hyperlink)
//

// Parses out alphanumeric names
func lexName(l *lexer) stateFunction {
	for {
		if strings.HasPrefix(l.input[l.pos:], LEFT_HYPER) {
			if l.pos > l.start {
				l.emit(itemName)
			}
			return lexLeftHyperlink
		}

		if strings.HasPrefix(l.input[l.pos:], RIGHT_HYPER) {
			return l.errorf("Invalid HTML Snippet")
		}

		if !IsAlphaNumeric(l.input[l.pos]) {
			if l.pos > l.start {
				l.emit(itemName)
			}
			return lexText
		}

		r := l.next()
		if r == '\n' || r == '\r' {
			return l.errorf("Invalid HTML Snippet - has line break")
		} else if r == eof {
			break
		}
	}

	if l.pos > l.start {
		l.emit(itemName)
	}

	l.emit(itemEOF)

	return nil
}

func lexLeftHyperlink(l *lexer) stateFunction {
	l.pos += len(LEFT_HYPER)
	return lexInsideHyperlink
}

func lexRightHyperlink(l *lexer) stateFunction {
	l.pos += len(RIGHT_HYPER)
	l.emit(itemHyperlink)
	return lexName
}

// Parsers out hyperlinks
func lexInsideHyperlink(l *lexer) stateFunction {
	for {
		if strings.HasPrefix(l.input[l.pos:], LEFT_HYPER) {
			return l.errorf("Invalid HTML Snippet")
		}

		if strings.HasPrefix(l.input[l.pos:], RIGHT_HYPER) {
			return lexRightHyperlink
		}

		r := l.next()
		if r == '\n' || r == '\r' {
			return l.errorf("Invalid HTML Snippet - has line break")
		} else if r == eof {
			return l.errorf("Invalid HTML Snippet - unclosed hyperlink")
			break
		}
	}
	return nil
}

// Parsers out non-name text
func lexText(l *lexer) stateFunction {
	for {
		if strings.HasPrefix(l.input[l.pos:], LEFT_HYPER) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexLeftHyperlink
		}
		if IsAlphaNumeric(l.input[l.pos]) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexName
		}

		r := l.next()
		if r == '\n' || r == '\r' {
			return l.errorf("Invalid HTML Snippet - has line break")
		} else if r == eof {
			break
		}
	}

	if l.pos > l.start {
		l.emit(itemText)
	}

	l.emit(itemEOF)

	return nil
}

// Help Functions

func IsAlphaNumeric(i uint8) bool {
	r := rune(i)
	return 'A' <= r && r <= 'Z' || 'a' <= r && r <= 'z' || '0' <= r && r <= '9'
}
