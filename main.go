package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/sjhitchner/annotator/interfaces/db"
	"github.com/sjhitchner/annotator/interfaces/rest"
	uc "github.com/sjhitchner/annotator/usecases"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	verbose := flag.Bool("verbose", false, "verbose logging")
	flag.Parse()

	// Create new Gorilla mux router
	router := mux.NewRouter()

	// Create a names repository and Annotation Interactor
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

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	log.Fatal(http.ListenAndServe(":3001", nil))
}
