package disadvantages

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const downstreamServer = "http://www.example.com"

// GetchRates rates from downstream sevice
type FetchRates struct{}

func (f *FetchRates) Fetch() ([]Rate, error) {
	// build the URL from which to fetch the rates
	url := downstreamServer + "/rates"

	// build request
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	//  fetch rates
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// read the content of the response
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// convert JSON bytes to Go structs
	out := &downStreamResponse{}
	err = json.Unmarshal(data, out)
	if err != nil {
		return nil, err
	}

	return out.Rates, nil
}

type downStreamResponse struct {
	Rates []Rate `json:"rates"`
}

type Rate struct {
	Code  string
	Value float64
}
