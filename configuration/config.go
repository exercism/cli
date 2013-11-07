package configuration

type Config struct {
	GithubUsername    string `json:"githubUsername"`
	ApiKey            string `json:"apiKey"`
	ExercismDirectory string `json:"exercismDirectory"`
	Hostname          string `json:"hostname"`
}
