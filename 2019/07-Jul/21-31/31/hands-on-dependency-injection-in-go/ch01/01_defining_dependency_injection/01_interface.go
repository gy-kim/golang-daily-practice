package defining_dependency_injection

import (
	"encoding/json"
	"errors"
)

// Saver persists the supplied bytes
type Saver interface {
	Save(data []byte) error
}

// SaverPerson will validate and persists the supplied person
func SavePerson(person *Person, saver Saver) error {
	// validate the inputs
	err := person.validate()

	if err != nil {
		return err
	}

	// encode person to bytes
	bytes, err := person.encode()
	if err != nil {
		return err
	}

	// save the person and return the result
	return saver.Save(bytes)
}

// Person data object
type Person struct {
	Name  string
	Phone string
}

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
