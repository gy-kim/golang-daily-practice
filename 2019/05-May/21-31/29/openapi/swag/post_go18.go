package swag

import "net/url"

func pathUnescape(path string) (string, error) {
	return url.PathUnescape(path)
}
