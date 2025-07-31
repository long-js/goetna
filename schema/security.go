package schema

import "time"

// Security...
type Security struct {
	Symbol          string    `json:"Symbol"`
	Description     string    `json:"Description"`
	Exchange        string    `json:"Exchange"`
	Currency        string    `json:"Currency"`
	AddedDate       time.Time `json:"AddedDate"`
	ModifyDate      time.Time `json:"ModifyDate"`
	Type            string    `json:"Type"`
	Id              int32     `json:"Id"`
	TickSize        float64   `json:"TickSize"`
	ContractSize    float64   `json:"ContractSize"`
	Precision       uint8     `json:"Precision"`
	VolumePrecision uint8     `json:"VolumePrecision"`
	Enabled         bool      `json:"Enabled"`
	AllowTrade      bool      `json:"AllowTrade"`
	AllowMargin     bool      `json:"AllowMargin"`
	AllowShort      bool      `json:"AllowShort"`
}

type TradingSession string
