package api

// Submission is an iteration that has been submitted to the API.
type Submission struct {
	URL           string            `json:"url"`
	TrackID       string            `json:"track_id"`
	Language      string            `json:"language"`
	Slug          string            `json:"slug"`
	Name          string            `json:"name"`
	Username      string            `json:"username"`
	ProblemFiles  map[string]string `json:"problem_files"`
	SolutionFiles map[string]string `json:"solution_files"`
	Iteration     int               `json:"iteration"`
}
