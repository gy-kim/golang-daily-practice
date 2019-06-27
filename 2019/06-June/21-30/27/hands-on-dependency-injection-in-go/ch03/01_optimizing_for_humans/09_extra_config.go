package humans

func PetFetcher(search string, limit int, offset int, sortBy string, sortAscending bool) []Pet {
	return []Pet{}
}

func PetFetcherTypicalUsage() {
	_ = PetFetcher("Fido", 0, 0, "", true)
}
