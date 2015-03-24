package rest

import (
	"github.com/gorilla/mux"
	uc "github.com/sjhitchner/annotator/usecases"
	"io/ioutil"
	"log"
	"net/http"
)

type AnnotateRequest struct {
	Text string
}

type annotateResourceImpl struct {
	interactor uc.AnnotateInteractor
}

func NewAnnotateResource(interactor uc.AnnotateInteractor) AnnotateResource {
	return annotateResourceImpl{
		interactor,
	}
}

func (t annotateResourceImpl) Register(router *mux.Router) {
	router.Methods("POST").
		Headers(HEADER_CONTENT_TYPE, CONTENT_TYPE_TEXT).
		HandlerFunc(t.AnnotateHTML)
	router.Methods("POST").
		Headers(HEADER_CONTENT_TYPE, CONTENT_TYPE_HTML).
		HandlerFunc(t.AnnotateHTML)
}

func (t annotateResourceImpl) AnnotateHTML(resp http.ResponseWriter, req *http.Request) {
	log.Println("annotate")

	buffer, err := ioutil.ReadAll(req.Body)
	if err != nil {
		ERROR(resp, err)
		return
	}

	html, err := t.interactor.AnnotateHTML(string(buffer))
	if err != nil {
		ERROR(resp, err)
		return
	}

	OK(resp, html)
	return
}
