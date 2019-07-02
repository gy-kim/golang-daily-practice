package disadvantages

import "errors"

func NewClient(service DepService) Client {
	return &clientImpl{
		service: service,
	}
}

type Client interface {
	DosomethingUseful() (bool, error)
}

type clientImpl struct {
	service DepService
}

type DepService interface {
	DoSomethingElse()
}

func (c *clientImpl) DosomethingUseful() (bool, error) {
	// this function does something useful
	return false, errors.New("not implemented")
}
