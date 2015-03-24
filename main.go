package main

import (
	"github.com/gorilla/mux"
	"github.com/sjhitchner/sourcegraph/interfaces/db"
	"github.com/sjhitchner/sourcegraph/interfaces/rest"
	uc "github.com/sjhitchner/sourcegraph/usecases"
	"log"
	"net/http"
)

func main() {

	router := mux.NewRouter()
	store := db.NewNamesRepository()
	interactor := uc.NewAnnotationInteractor(store)

	namesResource := rest.NewNamesResource(interactor)
	namesResource.Register(router.PathPrefix("/names").Subrouter())

	annotateResource := rest.NewAnnotateResource(interactor)
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
