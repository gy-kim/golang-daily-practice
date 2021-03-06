package unit_tests

import "database/sql"

type PetSaver struct{}

// Save the supplied pet and return the ID
func (p PetSaver) Save(pet Pet) (int, error) {
	err := p.validate(pet)
	if err != nil {
		return 0, err
	}

	result, err := p.save(pet)
	if err != nil {
		return 0, err
	}

	return p.extractID(result)
}

func (p PetSaver) validate(pet Pet) error {
	return nil
}

func (p PetSaver) save(pet Pet) (sql.Result, error) {
	return nil, nil
}

func (p PetSaver) extractID(result sql.Result) (int, error) {
	return 0, nil
}
