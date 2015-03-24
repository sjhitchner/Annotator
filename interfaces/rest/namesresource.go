package rest

import (
	"github.com/gorilla/mux"
	. "github.com/sjhitchner/annotator/domain"
	uc "github.com/sjhitchner/annotator/usecases"
	"net/http"
)

type NameUrlRequest struct {
	URL URL `json:"url"`
}

type NameUrlResponse struct {
	Name Name `json:"name"`
	URL  URL  `json:"url"`
}

type namesResourceImpl struct {
	interactor uc.NamesInteractor
}

func NewNamesResource(interactor uc.NamesInteractor) NamesResource {
	return namesResourceImpl{
		interactor,
	}
}

func (t namesResourceImpl) Register(router *mux.Router) {
	router.Methods("GET").
		Path("/{name:[A-Za-z0-9]+}").
		HandlerFunc(t.RetrieveName)

	router.Methods("PUT").
		//Path("/{name:[A-Za-z0-9]+}").
		Path("/{name}").
		Headers(HEADER_CONTENT_TYPE, CONTENT_TYPE_JSON).
		HandlerFunc(t.UpdateURLForName)

	router.Methods("DELETE").
		HandlerFunc(t.RemoveAllNames)
}

func getRequestName(request *http.Request) (Name, error) {
	params := mux.Vars(request)
	name := Name(params["name"])
	if err := name.Validate(); err != nil {
		return "", err
	}
	return name, nil
}

func (t namesResourceImpl) UpdateURLForName(response http.ResponseWriter, request *http.Request) {
	name, err := getRequestName(request)
	if err != nil {
		ERROR(response, err)
		return
	}

	var nameRequest NameUrlRequest
	if err := ReadPayload(request, &nameRequest); err != nil {
		ERROR(response, err)
		return
	}

	if err := t.interactor.UpdateURLForName(name, nameRequest.URL); err != nil {
		ERROR(response, err)
		return
	}

	OK(response, "")
	return
}

func (t namesResourceImpl) RetrieveName(response http.ResponseWriter, request *http.Request) {
	name, err := getRequestName(request)
	if err != nil {
		ERROR(response, err)
		return
	}

	url, err := t.interactor.GetURLForName(name)
	if err != nil {
		ERROR(response, err)
		return
	}

	result := NameUrlResponse{
		Name: name,
		URL:  url,
	}

	OK(response, result)
	return
}

func (t namesResourceImpl) RemoveAllNames(response http.ResponseWriter, request *http.Request) {
	if err := t.interactor.DeleteAllNames(); err != nil {
		ERROR(response, err)
		return
	}

	OK(response, "")
	return
}
