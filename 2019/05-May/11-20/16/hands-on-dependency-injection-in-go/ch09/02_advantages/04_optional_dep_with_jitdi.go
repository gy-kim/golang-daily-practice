package advantages

func NewLoaderWithJIT(ds DataStore) *LoaderWithJIT {
	return &LoaderWithJIT{
		datastore: ds,
	}
}

type LoaderWithJIT struct {
	datastore DataStore

	OptionalCache Cache
}

func (l *LoaderWithJIT) Load(ID int) (*Animal, error) {
	// attemp to load from cache
	output := l.cache().Get(ID)
	if output != nil {
		return output, nil
	}

	// load from data store
	output, err := l.datastore.Load(ID)
	if err != nil {
		return nil, err
	}

	// cache the loaded value
	l.cache().Put(ID, output)

	// output the result
	return output, nil
}

func (l *LoaderWithJIT) cache() Cache {
	if l.OptionalCache == nil {
		l.OptionalCache = &noopCache{}
	}

	return l.OptionalCache
}

type noopCache struct {
}

func (n *noopCache) Get(ID int) *Animal {
	return nil
}

func (n *noopCache) Put(ID int, value *Animal) {

}
