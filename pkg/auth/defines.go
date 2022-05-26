package auth

type Authentication struct {
	AccessToken string `json:"access_token"`
	Success     bool   `json:"success"`
	Error       string `json:"error,omitempty"`
}
