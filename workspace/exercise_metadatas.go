package workspace

// ExerciseMetadatas is a collection of exercise metadata to interactively choose from.
type ExerciseMetadatas []*ExerciseMetadata

// NewExerciseMetadatas loads up the exercise metadata for each of the provided paths.
func NewExerciseMetadatas(paths []string) (ExerciseMetadatas, error) {
	var metadatas []*ExerciseMetadata

	for _, path := range paths {
		metadata, err := NewExerciseMetadata(path)
		if err != nil {
			return []*ExerciseMetadata{}, err
		}
		metadatas = append(metadatas, metadata)
	}
	return metadatas, nil
}
