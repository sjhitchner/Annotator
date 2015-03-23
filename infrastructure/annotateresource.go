package infrastructure

import (
	"github.com/emicklei/go-restful"
	uc "github.com/sjhitchner/sourcegraph/usecases"
	"log"
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

func (t annotateResourceImpl) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/annotate").
		Doc("Manage Annotations")
	//("plain/html", "").
	//Produces("plain/html")

	ws.Route(ws.POST("").To(t.AnnotateHTML).
		Doc("update a url for name").
		Operation("updateURLForName").
		Param(ws.PathParameter("name", "name of person").DataType("string")).
		Returns(400, "malformed request", nil).
		Reads(AnnotateRequest{}))

	container.Add(ws)
}

func (t annotateResourceImpl) AnnotateHTML(request *restful.Request, response *restful.Response) {
	log.Println("annotate")
}
