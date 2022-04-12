package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "Get", []postData{}, http.StatusOK},
	{"about", "/about", "Get", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "Get", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "Get", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "Get", []postData{}, http.StatusOK},
	{"mr", "/make-reservations", "Get", []postData{}, http.StatusOK},
	{"contact", "/contact", "Get", []postData{}, http.StatusOK},
	{"post-sa", "/search-availability", "Post", []postData{
		{key: "start", value: "2022-01-01"},
		{key: "end", value: "2022-01-02"},
	}, http.StatusOK},
	{"post-sa-json", "/search-availability-json", "Post", []postData{
		{key: "start", value: "2022-01-01"},
		{key: "end", value: "2022-01-02"},
	}, http.StatusOK},
	{"post-mr", "/make-reservations", "Post", []postData{
		{key: "first_name", value: "cena"},
		{key: "last_name", value: "Jhon"},
		{key: "email", value: "g@g.com"},
		{key: "phone", value: "555-555-55-5"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "Get" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d, but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		}
	}
}
