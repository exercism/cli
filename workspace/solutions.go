package workspace

// Solutions is a collection of solutions to interactively choose from.
type Solutions []*Solution

// NewSolutions loads up the solution metadata for each of the provided paths.
func NewSolutions(paths []string) (Solutions, error) {
	var solutions []*Solution

	for _, path := range paths {
		solution, err := NewSolution(path)
		if err != nil {
			return []*Solution{}, err
		}
		solutions = append(solutions, solution)
	}
	return solutions, nil
}
