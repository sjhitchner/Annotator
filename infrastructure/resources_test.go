package infrastructure

import (
	"github.com/emicklei/go-restful"
	. "github.com/sjhitchner/sourcegraph/domain"
	. "github.com/sjhitchner/sourcegraph/usecases"
	"net/http/httptest"
	"testing"
)

func NewTestServer(container *restful.Container) *httptest.Server {
	ts := httptest.NewServer(container)
	return ts
}

type MockNameRepo struct {
}

func (t MockNameRepo) Get(name Name) (URL, error) {

}

func (t MockNameRepo) Put(name Name, url URL) error {

}

func (t MockNameRepo) DeleteAll() error {

}

type MockNamesInteractor struct {
}

func (t MockNamesInteractor) UpdateURLForName(name Name, url URL) error {
	return nil
}

func (t MockNamesInteractor) GetURLForName(name Name) (URL, error) {
	return "", nil
}

func (t MockNamesInteractor) DeleteAllNames() error {
	return nil
}

func TestGetNameForURL(t *testing.T) {
	container := restful.NewContainer()
	resource := NewNamesResource(MockNamesInteractor{})
	resource.Register(container)
	ts := NewTestServer(container)
}
