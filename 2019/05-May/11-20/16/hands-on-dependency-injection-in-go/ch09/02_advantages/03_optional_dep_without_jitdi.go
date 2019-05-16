package advantages

func NewLoaderWithoutJIT(ds DataStore) *LoaderWithoutJIT {
	return &LoaderWithoutJIT{
		datastore: ds,
	}
}

type LoaderWithoutJIT struct {
	datastore DataStore

	OptionalCache Cache
}

type Cache interface {
	Get(ID int) *Animal
	Put(ID int, value *Animal)
}

type DataStore interface {
	Load(ID int) (*Animal, error)
	Save(ID int, value *Animal) error
}

type Animal struct {
}
