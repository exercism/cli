package api

import (
	"fmt"
	"regexp"
)

// ValidateToken checks if token is a valid UUID,
// and verifies that the remote API accepts it, unless opted out.
func ValidateToken(token string, noAuthCheck bool) error {
	tokenIsUUID, err := regexp.MatchString("^[[:alnum:]]{8}-([[:alnum:]]{4}-){3}[[:alnum:]]{12}$", token)
	if err != nil {
		return err
	}

	if !tokenIsUUID {
		return fmt.Errorf("the token \"%s\" doesn't look like a valid token", token)
	}

	if noAuthCheck {
		return nil
	}

	client, err := NewClient()
	if err != nil {
		return err
	}
	client.UserConfig.Token = token

	return client.checkAuthorization()
}
