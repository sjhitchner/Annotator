package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	. "github.com/sjhitchner/annotator/domain"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	HEADER_CONTENT_TYPE = "Content-Type"

	CONTENT_TYPE_JSON = "application/json"
	CONTENT_TYPE_TEXT = "text/plain"
	CONTENT_TYPE_HTML = "text/html"
)

type Resource interface {
	Register(router *mux.Router)
}

type NamesResource interface {
	Resource
	UpdateURLForName(resp http.ResponseWriter, req *http.Request)
	RetrieveName(resp http.ResponseWriter, req *http.Request)
	RemoveAllNames(resp http.ResponseWriter, req *http.Request)
}

type AnnotateResource interface {
	Resource
	AnnotateHTML(resp http.ResponseWriter, req *http.Request)
}

func ReadPayload(request *http.Request, pointer interface{}) error {
	contentType := request.Header.Get(HEADER_CONTENT_TYPE)

	if strings.Contains(contentType, CONTENT_TYPE_JSON) {
		buffer, err := ioutil.ReadAll(request.Body)
		if err != nil {
			return err
		}
		decoder := json.NewDecoder(bytes.NewReader(buffer))
		decoder.UseNumber()
		return decoder.Decode(pointer)
	}
	return fmt.Errorf("Unable to unmarshal JSON payload")
}

func OK(response http.ResponseWriter, payload interface{}) {
	/*
		if LOG_RESPONSE {
			b, err := json.MarshalIndent(payload, "", "    ")
			if err != nil {
				nerr.NewInternalError(err, "marshalling error")
			}
			log.Println("RESPONSE\n", string(b))
		}
	*/

	var buf []byte
	switch val := payload.(type) {
	case string:
		buf = []byte(val)
	default:
		var err error
		buf, err = json.Marshal(payload)
		if err != nil {
			ERROR(response, err)
			return
		}
	}

	response.WriteHeader(http.StatusOK)
	response.Write(buf)
	response.Write([]byte("\n"))
	return
}

func ERROR(response http.ResponseWriter, err error) {

	if strings.Contains(err.Error(), NOT_FOUND_ERROR) {
		response.WriteHeader(http.StatusNotFound)
		response.Write([]byte(err.Error()))
		return
	}

	if strings.Contains(err.Error(), MALFORMED_ERROR) {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte(err.Error()))
		return
	}

	response.WriteHeader(http.StatusInternalServerError)
	response.Write([]byte(err.Error()))
	return
}
