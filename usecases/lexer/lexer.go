package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	ItemError ItemType = iota
	ItemEOF
	ItemName
	ItemText
	ItemLeftAngle
	ItemRightAngle
	ItemHyperlink
	ItemLeftHyper
	ItemRightHyper

	eof = -1
)

type ItemType int

type Item struct {
	Type  ItemType
	Value string
}

func (t Item) String() string {
	switch t.Type {
	case ItemEOF:
		return "EOF"
	case ItemError:
		return t.Value
	}
	if len(t.Value) > 50 {
		return fmt.Sprintf("%d:%.50q...", t.Type, t.Value)
	}
	return fmt.Sprintf("%d:%q", t.Type, t.Value)
}

type stateFunction func(*lexer) stateFunction

// Lexer
type lexer struct {
	input string
	start int
	pos   int
	width int
	state stateFunction
	items chan Item
}

// No copying, just a slice
// Channel usage adds some overhead could use ring buffer
func NewLexer(input string) *lexer {
	l := &lexer{
		input: input,
		items: make(chan Item, 2),
		state: lexName,
	}
	return l
}

func (t *lexer) NextItem() Item {
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

func (t *lexer) emit(i ItemType) {
	t.items <- Item{i, t.input[t.start:t.pos]}
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
	t.items <- Item{ItemError, fmt.Sprintf(format, args...)}
	return nil
}

const (
	LEFT_ANGLE  = "<"
	RIGHT_ANGLE = ">"
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
				l.emit(ItemName)
			}
			return lexLeftHyperlink
		}

		if strings.HasPrefix(l.input[l.pos:], RIGHT_HYPER) {
			return l.errorf("Invalid HTML Snippet")
		}

		if !IsAlphaNumeric(l.input[l.pos]) {
			if l.pos > l.start {
				l.emit(ItemName)
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
		l.emit(ItemName)
	}

	l.emit(ItemEOF)

	return nil
}

func lexLeftHyperlink(l *lexer) stateFunction {
	l.pos += len(LEFT_HYPER)
	return lexInsideHyperlink
}

func lexRightHyperlink(l *lexer) stateFunction {
	l.pos += len(RIGHT_HYPER)
	l.emit(ItemHyperlink)
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
				l.emit(ItemText)
			}
			return lexLeftHyperlink
		}
		if IsAlphaNumeric(l.input[l.pos]) {
			if l.pos > l.start {
				l.emit(ItemText)
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
		l.emit(ItemText)
	}

	l.emit(ItemEOF)

	return nil
}

// Help Functions

func IsAlphaNumeric(i uint8) bool {
	r := rune(i)
	return 'A' <= r && r <= 'Z' || 'a' <= r && r <= 'z' || '0' <= r && r <= '9'
}
