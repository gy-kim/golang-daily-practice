package disadvantages

type MyPresonLoader interface {
	Load(ID int) (*Person, error)
}
