package schema

type EtnaResponse struct {
	State  string `json:"State,omitempty"`
	Step   string `json:"Step,omitempty"`
	Reason string `json:"Reason,omitempty"`
}

type FmpResponse struct {
	Event     string `json:"event"`
	Message   string `json:"message"`
	Status    int16  `json:"status"`
	Timestamp uint64 `json:"timestamp"`
}
