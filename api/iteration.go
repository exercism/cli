package api

// Iteration represents a version of a particular exercise.
// This gets submitted to the API.
type Iteration struct {
	Key  string `json:"key"`
	Code string `json:"code"`
	Path string `json:"path"`
	Dir  string `json:"dir"`
}
