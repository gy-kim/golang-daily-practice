package defining_dependency_injection

import (
	"encoding/json"
	"errors"
)

// Saver persists the supplied bytes
type Saver interface {
	Save(data []byte) error
}

// SavePerson will validate and persist the supplied person
func SavePerson(person *Person, saver Saver) error {
	// validate the inputs
	err := person.validate()

	if err != nil {
		return err
	}

	bytes, err := person.encode()
	if err != nil {
		return err
	}

	return saver.Save(bytes)
}

// Person data object
type Person struct {
	Name  string
	Phone string
}

// validate the person object
func (p *Person) validate() error {
	if p.Name == "" {
		return errors.New("name missing")
	}
	if p.Phone == "" {
		return errors.New("phone missing")
	}
	return nil
}

func (p *Person) encode() ([]byte, error) {
	return json.Marshal(p)
}
