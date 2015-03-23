package usecases

import (
	"bytes"
	"fmt"
	. "github.com/sjhitchner/sourcegraph/domain"
	"github.com/sjhitchner/sourcegraph/usecases/lexer"
	"strings"
)

type annotateInteractorImpl struct {
	repo NameRepository
}

func newAnnotateInteractor(repo NameRepository) AnnotateInteractor {
	return &annotateInteractorImpl{
		repo,
	}
}

func (t annotateInteractorImpl) AnnotateHTML(html string) (string, error) {
	var b bytes.Buffer

	lex := lexer.NewLexer(html)

	token := lex.NextItem()
	for ; token.Type != lexer.ItemEOF; token = lex.NextItem() {
		switch token.Type {
		case lexer.ItemName:
			name := Name(token.Value)
			url, err := t.repo.Get(name)
			if err == nil {
				b.WriteString(`<a href="`)
				b.WriteString(string(url))
				b.WriteString(`">`)
			}

			b.WriteString(string(name))

			if err == nil {
				b.WriteString(`</a>`)
			}
		case lexer.ItemError:
			return "", fmt.Errorf(token.Value)
		default:
			b.WriteString(token.Value)
		}

	}
	return strings.TrimSpace(b.String()), nil
}

/*
func main() {
	str := "steve h    alex.com <a href=\"qwerty\">qwerty</a> qwerty"
	if str2, err := Annotate(str); err == nil {
		fmt.Println("RESULT:", str2)
	} else {
		fmt.Println("ERROR:", err)
	}
}

func Annotate(str string) (string, error) {
	var b bytes.Buffer

	lexer := NewLexer(str)
	for token := lexer.NextItem(); token.typ != itemEOF; {
		fmt.Println(token)

		switch token.typ {
		case itemName:
			b.WriteString("<a href=\"\">")
			b.WriteString(token.value)
			b.WriteString("</a>")
		case itemError:
			return "", fmt.Errorf(token.value)
		default:
			b.WriteString(token.value)
		}
	}
	return b.String(), nil
}
*/
