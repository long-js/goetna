package schema

type WSSessionType uint8
type SessionId string

type SessionResp struct {
	Id SessionId `json:"Id"`
}

type Streamer struct {
	Url       string    `json:"Url"`
	Type      string    `json:"Type"`
	SessionId SessionId `json:"SessionId"`
}

type Streamers struct {
	QuoteAddresses []Streamer `json:"QuoteAddresses"`
	DataAddresses  []Streamer `json:"DataAddresses"`
}

type Subscription struct {
	Cmd            string    `json:"Cmd"`
	SessionId      SessionId `json:"SessionId"`
	Keys           string    `json:"Keys"`
	Topic          string    `json:"EntityType"`
	HttpClientType string    `json:"HttpClientType"`
}
