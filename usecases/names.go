package usecases

import (
	. "github.com/sjhitchner/sourcegraph/domain"
	"log"
)

type nameInteractorImpl struct {
	repo NameRepository
}

func newNameInteractor(repo NameRepository) NameInteractor {
	return &nameInteractorImpl{
		repo,
	}
}

func (t nameInteractorImpl) UpdateURLForName(name Name, url URL) error {
	log.Printf("UpdateURLForName [%s] [%s]\n", name, url)

	if err := name.Validate(); err != nil {
		return err
	}

	if err := url.Validate(); err != nil {
		return err
	}

	return t.repo.Put(name, url)
}

func (t nameInteractorImpl) GetURLForName(name Name) (URL, error) {
	log.Printf("GetURLForName [%s]\n", name)

	if err := name.Validate(); err != nil {
		return "", err
	}

	url, err := t.repo.Get(name)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (t nameInteractorImpl) DeleteAllNames() error {
	log.Printf("DeleteAllNames\n")

	return t.repo.DeleteAll()
}
