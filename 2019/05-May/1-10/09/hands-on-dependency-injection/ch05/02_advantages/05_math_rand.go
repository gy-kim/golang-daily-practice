package advantages

// A Rand is a source of random numbers.
type Rand struct {
	src Source
}

// Int returns a non-negarive pseudo-random int.
func (r *Rand) Int() int {
	// code chaged for brevity
	value := r.src.Int63()
	return int(value)
}

/*
 * Top-level conveniene functions
 */

var globalRand = New(&lockedSource{})

// Int returns a non-nerative pseudo-random int from the default Source.
func Int() int { return globalRand.Int() }

// New returns a new Rand that uses random values from src
// to generate other random values.
func New(src Source) *Rand {
	// code changed for brevity
	return &Rand{
		src: src,
	}
}

type lockedSource struct {
}

func (l *lockedSource) Int63() int64 {
	return 0
}

// A Source represents a source of uniformly-distributetd
// psuedo-randomo int64 values in the range [0, 1<<64]
type Source interface {
	Int63() int64
}
