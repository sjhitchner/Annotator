package db

import (
	"fmt"
	. "github.com/sjhitchner/annotator/domain"
	"sync"
)

type mapBasedNameRepositoryImpl struct {
	sync.RWMutex
	mapper map[Name]URL
}

func NewNamesRepository() NamesRepository {
	return &mapBasedNameRepositoryImpl{
		sync.RWMutex{},
		make(map[Name]URL),
	}
}

func (t mapBasedNameRepositoryImpl) Get(name Name) (URL, error) {
	t.RLock()
	defer t.RUnlock()

	url, ok := t.mapper[name]
	if !ok {
		return "", fmt.Errorf("name [%s] %s", name, NOT_FOUND_ERROR)
	}
	return url, nil
}

func (t *mapBasedNameRepositoryImpl) Put(name Name, url URL) error {
	t.Lock()
	defer t.Unlock()

	t.mapper[name] = url
	return nil
}

func (t *mapBasedNameRepositoryImpl) DeleteAll() error {
	t.Lock()
	defer t.Unlock()

	for k := range t.mapper {
		delete(t.mapper, k)
	}
	return nil
}
