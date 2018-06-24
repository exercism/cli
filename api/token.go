package api

// ValidateToken calls the API to determine whether the token is valid.
func ValidateToken() error {
	client, err := NewClient()
	if err != nil {
		return err
	}
	return client.checkAuthorization()
}
