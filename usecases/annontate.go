package usecases

import (
//. "github.com/sjhitchner/sourcegraph/domain"
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
	return html, nil
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
