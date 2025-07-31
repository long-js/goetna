package schema

import "time"

// Security...
type Security struct {
	Id              int       `json:"Id"`
	Symbol          string    `json:"Symbol"`
	Description     string    `json:"Description"`
	Exchange        string    `json:"Exchange"`
	Currency        string    `json:"Currency"`
	AddedDate       time.Time `json:"AddedDate"`
	ModifyDate      time.Time `json:"ModifyDate"`
	Type            string    `json:"Type"`
	Precision       int       `json:"Precision"`
	VolumePrecision int       `json:"VolumePrecision"`
	TickSize        float64   `json:"TickSize"`
	Enabled         bool      `json:"Enabled"`
	AllowTrade      bool      `json:"AllowTrade"`
	AllowMargin     bool      `json:"AllowMargin"`
	AllowShort      bool      `json:"AllowShort"`
}

type TradingSession string
