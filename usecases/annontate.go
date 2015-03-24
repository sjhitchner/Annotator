package usecases

import (
	"bytes"
	"fmt"
	. "github.com/sjhitchner/sourcegraph/domain"
	"github.com/sjhitchner/sourcegraph/usecases/lexer"
	"strings"
)

type annotateInteractorImpl struct {
	repo NamesRepository
}

func newAnnotateInteractor(repo NamesRepository) AnnotateInteractor {
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
