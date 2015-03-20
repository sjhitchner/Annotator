package infrastructure

import (
	"github.com/emicklei/go-restful"
	. "github.com/sjhitchner/sourcegraph/domain"
	uc "github.com/sjhitchner/sourcegraph/usecases"
)

type NameUrlRequest struct {
	URL URL `json:"url"`
}

type NameUrlResponse struct {
	Name Name `json:"name"`
	URL  URL  `json:"url"`
}

type namesResourceImpl struct {
	interactor uc.NameInteractor
}

func NewNamesResource(interactor uc.NameInteractor) NamesResource {
	return namesResourceImpl{
		interactor,
	}
}

func (t namesResourceImpl) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/names").
		Doc("Manage Names").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.PUT("/{name}").To(t.UpdateURLForName).
		Doc("update a url for name").
		Operation("updateURLForName").
		Param(ws.PathParameter("name", "name of person").DataType("string")).
		Returns(400, "malformed request", nil).
		Reads(NameUrlRequest{}))

	ws.Route(ws.GET("/{name}").To(t.RetrieveName).
		Doc("Retrieve a name").
		Operation("retrieveName").
		Param(ws.PathParameter("name", "name of person").DataType("string")).
		Returns(400, "invalid name", nil).
		Returns(404, "name does not exist", nil).
		Writes(NameUrlResponse{}))

	ws.Route(ws.DELETE("").To(t.RemoveAllNames).
		Doc("delete all names").
		Operation("removeAllNames").
		Param(ws.PathParameter("name", "name of person").DataType("string")))

	container.Add(ws)
}

func getRequestName(req *restful.Request) (Name, error) {
	name := Name(req.PathParameter("name"))
	if err := name.Validate(); err != nil {
		return "", err
	}
	return name, nil
}

func (t namesResourceImpl) UpdateURLForName(request *restful.Request, response *restful.Response) {
	name, err := getRequestName(request)
	if err != nil {
		ERROR(response, err)
		return
	}

	var url URL
	if err := request.ReadEntity(&url); err != nil {
		ERROR(response, err)
		return
	}

	if err := t.interactor.UpdateURLForName(name, url); err != nil {
		ERROR(response, err)
		return
	}

	OK(response, nil)
	return
}

func (t namesResourceImpl) RetrieveName(request *restful.Request, response *restful.Response) {
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

func (t namesResourceImpl) RemoveAllNames(request *restful.Request, response *restful.Response) {
	if err := t.interactor.DeleteAllNames(); err != nil {
		ERROR(response, err)
		return
	}

	OK(response, nil)
	return
}
