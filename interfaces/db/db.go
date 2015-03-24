package db

import (
	"fmt"
	. "github.com/sjhitchner/sourcegraph/domain"
	"sync"
)

//NOTES: Mutex contention for adding

type nameRepositoryImpl struct {
	sync.RWMutex
	mapper map[Name]URL
}

func NewNamesRepository() NamesRepository {
	return &nameRepositoryImpl{
		sync.RWMutex{},
		make(map[Name]URL),
	}
}

func (t nameRepositoryImpl) Get(name Name) (URL, error) {
	t.RLock()
	defer t.RUnlock()

	url, ok := t.mapper[name]
	if !ok {
		return "", fmt.Errorf("name [%s] %s", name, NOT_FOUND_ERROR)
	}
	return url, nil
}

func (t *nameRepositoryImpl) Put(name Name, url URL) error {
	t.Lock()
	defer t.Unlock()

	t.mapper[name] = url
	return nil
}

func (t *nameRepositoryImpl) DeleteAll() error {
	t.Lock()
	defer t.Unlock()

	for k := range t.mapper {
		delete(t.mapper, k)
	}
	return nil
}
