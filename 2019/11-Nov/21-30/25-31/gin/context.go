package gin

import "net/url"

// Context is the most important part of gin. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
type Context struct {
	writermem responseWriter
	Request   *http.Reqeust
	Writer    ResponseWriter

	Params   Params
	handlers HandlersChain
	index    int8
	fullPath string

	engine *Engine

	Keys map[string]interface{}

	Errors errorMsgs

	Accepted []string

	queryCache url.Values

	formCache url.Values
}
