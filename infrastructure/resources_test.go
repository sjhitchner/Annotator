package infrastructure

import (
	"github.com/gorilla/mux"
	. "github.com/sjhitchner/sourcegraph/domain"
	//. "github.com/sjhitchner/sourcegraph/usecases"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func NewTestServer(router *mux.Router) *httptest.Server {
	ts := httptest.NewServer(container)
	return ts
}

func NewNamesServer() *httptest.Server {
	router := mux.NewRouter()
	resource := NewNamesResource(NewMockNamesInteractor())
	resource.Register(router)
	return NewTestServer(router)
}

/*
type MockNameRepo struct {
}

func (t MockNameRepo) Get(name Name) (URL, error) {

}

func (t MockNameRepo) Put(name Name, url URL) error {

}

func (t MockNameRepo) DeleteAll() error {

}
*/

type MockNamesInteractor struct {
	m map[Name]URL
}

func NewMockNamesInteractor() MockNamesInteractor {
	return MockNamesInteractor{
		make(map[Name]URL),
	}
}

func (t MockNamesInteractor) UpdateURLForName(name Name, url URL) error {
	t.m[name] = url
	return nil
}

func (t MockNamesInteractor) GetURLForName(name Name) (URL, error) {
	return t.m[name], nil
}

func (t MockNamesInteractor) DeleteAllNames() error {
	for k, _ := range t.m {
		delete(t.m, k)
	}
	return nil
}

func TestNamesResource(t *testing.T) {
	ts := NewNamesServer()
	defer ts.Close()

	name := Name("steve")
	url := URL("stephenhitchner.com")

	status1, text := CallUpdateURLForName(ts.URL, name, url)
	if status1 != http.StatusOK {
		t.Fatalf("Call failed [%d] [%s]", status1, text)
	}

	status2, urlResp := CallGetURLForName(ts.URL, name)
	if status2 != http.StatusOK {
		t.Fatalf("Call failed [%d] [%s]", status2, text)
	}

	if url != urlResp {
		t.Fatalf("[%s] expected got [%s]", url, urlResp)
	}
	/*
		res, err := http.Get(ts.URL)
		if err != nil {
			t.Fatal(err)
		}

		greeting, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			t.Fatal(err)
		}

		fmt.Printf("%s", greeting)
	*/
}

func CallUpdateURLForName(serverURL string, name Name, url URL) (int, string) {
	nr := NameUrlRequest{
		URL: url,
	}

	b, _ := json.Marshal(nr)

	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%s/names/%s", serverURL, name),
		bytes.NewBuffer(b),
	)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, string(body)
}

func CallGetURLForName(serverURL string, name Name) (int, URL) {
	resp, err := http.Get(fmt.Sprintf("%s/names/%s", serverURL, name))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	result := NameUrlResponse{}
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}

	return resp.StatusCode, result.URL
}
