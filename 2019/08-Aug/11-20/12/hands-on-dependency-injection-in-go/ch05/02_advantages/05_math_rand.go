package advantages

type Rand struct {
	src Source
}

func (r *Rand) Int() int {
	//  code changed for brevity
	value := r.src.Int63()
	return int(value)
}

var globalRand = New(&lockedSource{})

func Int() int { return globalRand.Int() }

func New(src Source) *Rand {
	return &Rand{
		src: src,
	}
}

type lockedSource struct {
}

func (l *lockedSource) Int63() int64 {
	return 0
}

type Source interface {
	Int63() int64
}
