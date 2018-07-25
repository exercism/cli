package config

import (
	"fmt"
	"regexp"
)

var (
	defaultBaseURL = "https://api.exercism.io/v1"
)

// InferSiteURL guesses what the website URL is.
// The basis for the guess is which API we're submitting to.
func InferSiteURL(apiURL string) string {
	if apiURL == "" {
		apiURL = defaultBaseURL
	}
	if apiURL == "https://api.exercism.io/v1" {
		return "https://exercism.io"
	}
	re := regexp.MustCompile("^(https?://[^/]*).*")
	return re.ReplaceAllString(apiURL, "$1")
}

// SettingsURL provides a link to where the user can find their API token.
func SettingsURL(apiURL string) string {
	return fmt.Sprintf("%s%s", InferSiteURL(apiURL), "/my/settings")
}
