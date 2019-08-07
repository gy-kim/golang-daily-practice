package unit_tests

type PersonLoaderStub struct {
	Person *Person
	Error  error
}

func (p *PersonLoaderStub) Load(ID int) (*Person, error) {
	return p.Person, p.Error
}
