package workspace

// MetadataCollection is a collection of solutions to interactively choose from.
type MetadataCollection []*Metadata

// NewMetadataCollection loads up the solution metadata for each of the provided paths.
func NewMetadataCollection(paths []string) (MetadataCollection, error) {
	var collection []*Metadata

	for _, path := range paths {
		metadata, err := NewMetadata(path)
		if err != nil {
			return nil, err
		}
		collection = append(collection, metadata)
	}
	return collection, nil
}
