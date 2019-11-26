package gin

const defaultMultipartMemory = 32 << 20 // 32MB

var (
	default404Body   = []byte(`404 page not found`)
	default405Body   = []byte(`405 method not allowed`)
	defaultAppEngine bool
)

type HandlerFunc func(*Context)
