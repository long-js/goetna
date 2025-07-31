package schema

type SFA struct {
	// The state of the request. 'Succeeded' value means that the request has been successfully made.
	State string `json:"State"`

	// WARNING: the authorization token lifetime is 24 hours.
	// The token that must be provided in all subsequent API requests as the authentication bearer token.
	// Token format: The value of this header must have the following format: Bearer BQ898r9fefi (Bearer + 1 space + the token).
	Token string `json:"Token,omitempty"`

	Step   string `json:"Step,omitempty"`
	Reason string `json:"Reason,omitempty"`
	Error  string `json:"error,omitempty"`
}
