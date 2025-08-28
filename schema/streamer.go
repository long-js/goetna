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
	FMPKey         string
}

type EtnaSubReq struct {
	Cmd            string    `json:"Cmd"`
	SessionId      SessionId `json:"SessionId"`
	Keys           string    `json:"Keys"`
	Topic          string    `json:"EntityType"`
	HttpClientType string    `json:"HttpClientType"`
}

type FmpStreamers struct {
	Success bool `json:"success"`
	Data    map[string]struct {
		Streamers Streamers         `json:"streamers"`
		Creds     map[string]string `json:"credentials"`
	} `json:"data"`
}

type FmpReq struct {
	Event string            `json:"event"`
	Data  map[string]string `json:"data"`
}
