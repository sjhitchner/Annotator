package infrastructure

import (
	"github.com/emicklei/go-restful"
	. "github.com/sjhitchner/sourcegraph/domain"
	"net/http"
	"strings"
)

type Resource interface {
	Register(container *restful.Container)
}

type NamesResource interface {
	Resource
	UpdateURLForName(request *restful.Request, response *restful.Response)
	RetrieveName(request *restful.Request, response *restful.Response)
	RemoveAllNames(request *restful.Request, response *restful.Response)
}

type AnnotateResource interface {
	Resource
	AnnotateHTML(request *restful.Request, response *restful.Response)
}

func OK(response *restful.Response, payload interface{}) {
	/*
		if LOG_RESPONSE {
			b, err := json.MarshalIndent(payload, "", "    ")
			if err != nil {
				nerr.NewInternalError(err, "marshalling error")
			}
			log.Println("RESPONSE\n", string(b))
		}
	*/

	response.WriteHeader(http.StatusOK)
	response.WriteAsJson(payload)
	return
}

func ERROR(response *restful.Response, err error) {
	if strings.Contains(err.Error(), NOT_FOUND_ERROR) {
		response.WriteErrorString(http.StatusNotFound, err.Error())
		return
	}
	if strings.Contains(err.Error(), MALFORMED_ERROR) {
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}

	response.WriteErrorString(http.StatusInternalServerError, err.Error())
	return
}
