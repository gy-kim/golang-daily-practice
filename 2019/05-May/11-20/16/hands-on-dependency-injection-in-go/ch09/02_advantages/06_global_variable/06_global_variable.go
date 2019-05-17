package global_variable

// Glboal singleton of connections to our data store
var storage UserStorage

type Savor struct {
}

func (s *Savor) Do(in *User) error {
	err := s.validate(in)
	if err != nil {
		return nil
	}

	return storage.Save(in)
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
