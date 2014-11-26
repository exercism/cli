package api

// Submission is an iteration that has been submitted to the API.
type Submission struct {
	URL      string 					 `json:"url"`
	TrackID  string 					 `json:"track_id"`
	Language string 					 `json:"language"`
	Slug     string 					 `json:"slug"`
	Name     string 					 `json:"name"`
	UserName string 					 `json:"username"`
	Problem	 map[string]string `json:"problem"`
	Code		 map[string]string `json:"solution"`
}
