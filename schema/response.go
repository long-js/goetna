package schema

type Response struct {
	State  string `json:"State,omitempty"`
	Step   string `json:"Step,omitempty"`
	Reason string `json:"Reason,omitempty"`
}
