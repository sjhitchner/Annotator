package usecases

import (
	. "github.com/sjhitchner/annotator/domain"
)

type NamesInteractor interface {
	UpdateURLForName(name Name, url URL) error
	GetURLForName(name Name) (URL, error)
	DeleteAllNames() error
}

type AnnotateInteractor interface {
	AnnotateHTML(html string) (string, error)
}

type AnnotationInteractor interface {
	NamesInteractor
	AnnotateInteractor
}

type annotationInteractor struct {
	NamesInteractor
	AnnotateInteractor
}

func NewAnnotationInteractor(repo NamesRepository) AnnotationInteractor {
	return &annotationInteractor{
		newNamesInteractor(repo),
		newAnnotateInteractor(repo),
	}
}
