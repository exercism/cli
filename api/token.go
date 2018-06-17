package api

import (
	"fmt"
	"regexp"
)

// ValidateToken checks if token is a valid UUID,
// and optionally verifies that the remote API accepts it.
func ValidateToken(token string, skipAuthCheck bool) error {
	tokenIsUUID, err := regexp.MatchString("^[[:alnum:]]{8}-([[:alnum:]]{4}-){3}[[:alnum:]]{12}$", token)
	if err != nil {
		return err
	}

	if !tokenIsUUID {
		return fmt.Errorf("the token \"%s\" doesn't look like a valid token", token)
	}

	if skipAuthCheck {
		return nil
	}

	client, err := NewClient()
	if err != nil {
		return err
	}
	client.UserConfig.Token = token

	return client.checkAuthorization()
}

// checkAuthorization calls the API to check if
// the client is authorized.
func (client *Client) checkAuthorization() error {
	url := client.APIConfig.URL("validate")
	req, err := client.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	_, err = client.Do(req, nil)

	return err
}
