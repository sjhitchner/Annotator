package main

import (
	"github.com/sjhitchner/sourcegraph/infrastructure"
	//"github.com/sjhitchner/sourcegraph/interfaces"
	"github.com/gorilla/mux"
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

	router := mux.NewRouter()
	store := infrastructure.NewNameRepository()
	interactor := uc.NewAnnotationInteractor(store)

	namesResource := infrastructure.NewNamesResource(interactor)
	namesResource.Register(router.PathPrefix("/names").Subrouter())

	annotateResource := infrastructure.NewAnnotateResource(interactor)
	annotateResource.Register(router.PathPrefix("/annotate").Subrouter())

	router.Path("/ping").
		Methods("GET").
		HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("pong"))
	})

	http.Handle("/", router)

	log.Printf("Started listening on localhost:3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}
