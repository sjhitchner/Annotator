package usecases

import (
	. "github.com/sjhitchner/sourcegraph/domain"
)

type NameRepository interface {
	Get(name Name) (URL, error)
	Put(name Name, url URL) error
	DeleteAll() error
}

type NameInteractor interface {
	UpdateURLForName(name Name, url URL) error
	GetURLForName(name Name) (URL, error)
	DeleteAllNames() error
}

type AnnotateInteractor interface {
	AnnotateHTML(html string) (string, error)
}

type AnnotationInteractor interface {
	NameInteractor
	AnnotateInteractor
}

type annotationInteractor struct {
	//nameInteractorImpl
	//annotateInteractorImpl
	NameInteractor
	AnnotateInteractor
}

func NewAnnotationInteractor(repo NameRepository) AnnotationInteractor {
	return &annotationInteractor{
		newNameInteractor(repo),
		newAnnotateInteractor(repo),
	}
}
