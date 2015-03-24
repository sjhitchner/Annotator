package usecases

import (
	"fmt"
	. "github.com/sjhitchner/annotator/domain"
	"github.com/sjhitchner/annotator/interfaces/db"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

func TestNames(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	repo := db.NewNamesRepository()
	interactor := NewAnnotationInteractor(repo)

	for _, name := range readNames() {
		err := interactor.UpdateURLForName(name, buildURL(name))
		if err != nil {
			t.Fatalf("Failed adding name [%s] to repo [%v]", name, err)
		}
	}

	for _, name := range readNames() {
		url, err := interactor.GetURLForName(name)
		if err != nil {
			t.Fatalf("Failed finding name [%s] in repo [%v]", name, err)
		}

		expected := buildURL(name)
		if url != expected {
			t.Fatalf("expected [%s] got url [%s]", expected, url)
		}
	}

	input := `Aaden Anton Brea Cicero Delinda Ellis Gavin Idabelle Jodi Kinte Lonnie Matie Nora Rettie Shelvie Toni Yvette Aaden Anton Brea Cicero Delinda Ellis Gavin Idabelle Jodi Kinte Lonnie Matie Nora Rettie Shelvie Toni Yvette`
	expected := `<a href="http://aaden.com">Aaden</a> <a href="http://anton.com">Anton</a> <a href="http://brea.com">Brea</a> <a href="http://cicero.com">Cicero</a> <a href="http://delinda.com">Delinda</a> <a href="http://ellis.com">Ellis</a> <a href="http://gavin.com">Gavin</a> <a href="http://idabelle.com">Idabelle</a> <a href="http://jodi.com">Jodi</a> <a href="http://kinte.com">Kinte</a> <a href="http://lonnie.com">Lonnie</a> <a href="http://matie.com">Matie</a> <a href="http://nora.com">Nora</a> <a href="http://rettie.com">Rettie</a> <a href="http://shelvie.com">Shelvie</a> <a href="http://toni.com">Toni</a> <a href="http://yvette.com">Yvette</a> <a href="http://aaden.com">Aaden</a> <a href="http://anton.com">Anton</a> <a href="http://brea.com">Brea</a> <a href="http://cicero.com">Cicero</a> <a href="http://delinda.com">Delinda</a> <a href="http://ellis.com">Ellis</a> <a href="http://gavin.com">Gavin</a> <a href="http://idabelle.com">Idabelle</a> <a href="http://jodi.com">Jodi</a> <a href="http://kinte.com">Kinte</a> <a href="http://lonnie.com">Lonnie</a> <a href="http://matie.com">Matie</a> <a href="http://nora.com">Nora</a> <a href="http://rettie.com">Rettie</a> <a href="http://shelvie.com">Shelvie</a> <a href="http://toni.com">Toni</a> <a href="http://yvette.com">Yvette</a>`
	html, err := interactor.AnnotateHTML(input)
	if err != nil {
		t.Fatalf("failed annotating html [%s] [%s]", input, err)
	}
	if expected != html {
		t.Fatalf("expected\n%s\ngot\n%s\n", expected, html)
	}

	if err := interactor.DeleteAllNames(); err != nil {
		t.Fatalf("Failed deleting all from repo [%v]", err)
	}

}

func readNames() []Name {
	content, err := ioutil.ReadFile("test_names.txt")
	if err != nil {
		panic(err)

	}
	arr := strings.Split(string(content), "\n")

	names := make([]Name, 0, len(arr))
	for _, n := range arr {
		name := Name(n)
		if err := name.Validate(); err != nil {
			continue
		}
		names = append(names, name)
	}
	return names
}

func buildURL(name Name) URL {
	return URL(fmt.Sprintf("http://%s.com", strings.ToLower(string(name))))
}

/*
f, err := os.Open("names.txt")
    if err != nil {
        panic(err)
    }



    reader := bufio.NewReader(f)
    b4, err := r4.Peek(5)
    check(err)
    fmt.Printf("5 bytes: %s\n", string(b4))
Close the file when you’re done (usually this would be scheduled immediately after Opening with defer).
    f.Close()


	interactor := setupInteractor(
		"alex", "http://alex.com",
	)

	annotateTest(
		t,
		interactor,
		`my name is alex`,
		`my name is <a href="http://alex.com">alex</a>`,
	)
}
*/
