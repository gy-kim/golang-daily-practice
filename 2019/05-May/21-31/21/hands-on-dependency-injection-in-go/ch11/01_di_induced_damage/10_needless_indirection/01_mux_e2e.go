package needless_indirection

import "net/http"

func buildRouter(mux *http.ServeMux) {
	mux.Handle("/get", &getEndpint{})
	mux.Handle("/list", &listEndpoint{})
	mux.Handle("/save", &saveEndpoint{})
}

type getEndpint struct{}

func (*getEndpint) ServeHTTP(resp http.ResponseWriter, _ *http.Request) {
	_, _ = resp.Write([]byte(`Hi from Get!`))
}

type listEndpoint struct{}

func (*listEndpoint) ServeHTTP(resp http.ResponseWriter, _ *http.Request) {
	_, _ = resp.Write([]byte(`Hi from List!`))
}

type saveEndpoint struct{}

func (*saveEndpoint) ServeHTTP(resp http.ResponseWriter, _ *http.Request) {
	_, _ = resp.Write([]byte(`Hi from Save!`))
}
