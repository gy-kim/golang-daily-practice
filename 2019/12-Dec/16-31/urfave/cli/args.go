package cli

type Args interface {
	// Get returns the nth argument, or else a blank string
	Get(n int) string

	First() string

	Tail() []string

	Len() int

	Present() bool

	Slice() []string
}

type args []string

func (a *args) Get(n int) string {
	if len(*a) > n {
		return (*a)[n]
	}
	return ""
}

func (a *args) First() string {
	return a.Get(0)
}

func (a *args) Tail() []string {
	if a.Len() >= 2 {
		tail := []string((*a)[1:])
		ret := make([]string, len(tail))
		copy(ret, tail)
		return ret
	}
	return []string{}
}

func (a *args) Len() int {
	return len(*a)
}

func (a *args) Present() bool {
	return a.Len() != 0
}

func (a *args) Slice() []string {
	ret := make([]string, len(*a))
	copy(ret, *a)
	return ret
}
