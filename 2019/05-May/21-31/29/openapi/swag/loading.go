package swag

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// LoadHTTPTimeout the default timeout for load requests
var LoadHTTPTimeout = 30 * time.Second

// LoadFromFileOrHTTP loads the bytes from a file or a remote http server based on the path passed in
func LoadFromFileOrHTTP(path string) ([]byte, error) {
	return LoadStrategy(path, ioutil.ReadFile, loadHTTPBytes(LoadHTTPTimeout))(path)
}

// LoadFromFileOrHTTPWithTimeout loads the bytes from a file or a remote http server based on the path passed in
// timeout arg allows for per request overriding of the request timeout
func LoadFromFileOrHTTPWithTimeout(path string, timeout time.Duration) ([]byte, error) {
	return LoadStrategy(path, ioutil.ReadFile, loadHTTPBytes(timeout))(path)
}

// LoadStrategy returns a loader function for a given path or uri
func LoadStrategy(path string, local, remote func(string) ([]byte, error)) func(string) ([]byte, error) {
	if strings.HasPrefix(path, "http") {
		return remote
	}
	return func(pth string) ([]byte, error) {
		upth, err := pathUnescape(pth)
		if err != nil {
			return nil, err
		}
		return local(filepath.FromSlash(upth))
	}
}

func loadHTTPBytes(timeout time.Duration) func(path string) ([]byte, error) {
	return func(path string) ([]byte, error) {
		client := &http.Client{Timeout: timeout}
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			return nil, err
		}
		resp, err := client.Do(req)
		defer func() {
			if resp != nil {
				if e := resp.Body.Close(); e != nil {
					log.Println(e)
				}
			}
		}()
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("could not access document at %q [%s]", path, resp.Status)
		}
		return ioutil.ReadAll(resp.Body)
	}
}
