package main

import (
	"github.com/emicklei/go-restful"
	"github.com/sjhitchner/sourcegraph/infrastructure"
	//"github.com/sjhitchner/sourcegraph/interfaces"
	uc "github.com/sjhitchner/sourcegraph/usecases"
	"log"
	"net/http"
)

/*
PUT /names/[name]
{"url": "[url goes here]"}

GET /names/[name]
{"name": "[name goes here]", "url":"[url goes here]"}

DELETE /names

POST /annotate
<html>

</html>

*/

func main() {
	wsContainer := restful.NewContainer()

	store := infrastructure.NewNameRepository()
	interactor := uc.NewAnnotationInteractor(store)

	namesResource := infrastructure.NewNamesResource(interactor)
	namesResource.Register(wsContainer)

	annotateResource := infrastructure.NewAnnotateResource(interactor)
	annotateResource.Register(wsContainer)

	log.Printf("Started listening on localhost:3001")
	server := &http.Server{Addr: ":3001", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}
