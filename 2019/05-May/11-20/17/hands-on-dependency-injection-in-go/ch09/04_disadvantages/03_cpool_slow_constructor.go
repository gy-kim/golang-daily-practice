package disadvantages

func newConnectionPool() ConnectionPool {
	pool := &myConnectionPool{}

	// initailize the pool
	pool.init()

	// return a "ready to use pool"
	return pool
}
