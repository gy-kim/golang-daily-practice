package global_variable_jit

// Global sigleton of connections to our data store
var storage UserStorage

type Savor struct {
	storage UserStorage
}

func (s *Savor) Do(in *User) error {
	err := s.validate(in)
	if err != nil {
		return err
	}

	return s.getStorage().Save(in)
}

// Just-in-time DI
func (s *Savor) getStorage() UserStorage {
	if s.storage == nil {
		s.storage = storage
	}

	return s.storage
}

func (s *Savor) validate(in *User) error {
	return nil
}

type UserStorage interface {
	Save(in *User) error
}

type User struct {
	Name     string
	Password string
}
