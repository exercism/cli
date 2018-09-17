package workspace

// ExerciseMetadataCollection is a collection of exercise metadata to interactively choose from.
type ExerciseMetadataCollection []*ExerciseMetadata

// NewExerciseMetadataCollection loads up the exercise metadata for each of the provided paths.
func NewExerciseMetadataCollection(paths []string) (ExerciseMetadataCollection, error) {
	var metadataCollection []*ExerciseMetadata

	for _, path := range paths {
		metadata, err := NewExerciseMetadata(path)
		if err != nil {
			return []*ExerciseMetadata{}, err
		}
		metadataCollection = append(metadataCollection, metadata)
	}
	return metadataCollection, nil
}
