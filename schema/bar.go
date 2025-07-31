package schema

import (
	"bytes"
	"strconv"
	"time"
)

// Bar...
type Bar struct {
	Open        float64 `json:"Open"`
	High        float64 `json:"High"`
	Low         float64 `json:"Low"`
	Close       float64 `json:"Close"`
	Volume      float64 `json:"Volume"`
	Time        uint32  `json:"Time"`
	IsCompleted bool
	IsMarket    bool `json:"IsMarket"` // Indicates if the bar is positioned during the RTH (true).
}

func (b Bar) Parse(values map[string]string) error {
	var err error

	for k, v := range values {
		switch k {
		case "Open":
			b.Open, err = strconv.ParseFloat(v, 64)
		case "High":
			b.High, err = strconv.ParseFloat(v, 64)
		case "Low":
			b.Low, err = strconv.ParseFloat(v, 64)
		case "Close":
			b.Close, err = strconv.ParseFloat(v, 64)
		case "Volume":
			b.Volume, err = strconv.ParseFloat(v, 64)
		case "Time":
			if t, err := strconv.ParseUint(v, 10, 32); err != nil {
				return err
			} else {
				b.Time = uint32(t)
			}
		case "IsCompleted":
			b.IsCompleted, err = strconv.ParseBool(v)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

type BarHist struct {
	Open   float64     `json:"open"`
	High   float64     `json:"high"`
	Low    float64     `json:"low"`
	Close  float64     `json:"close"`
	Volume float64     `json:"volume"`
	Time   BarHistTime `json:"date"`
}

const BarHistTimeLayout = "2006-01-02 15:04:05"

type BarHistTime time.Time

func (bt *BarHistTime) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)
	if t, err := time.Parse(BarHistTimeLayout, string(b)); err != nil {
		*bt = BarHistTime{}
		return err
	} else {
		*bt = BarHistTime(t)
	}
	return nil
}

func (bt BarHistTime) String() string {
	return time.Time(bt).Format(time.DateTime)
}

// BarSize...
type BarSize struct {
	Name    string
	Seconds uint32
}

type ReqBars struct {
	Ticker     string         `schema:"ticker"`
	ExchangeId uint8          `schema:"exchange_id"`
	Options    ReqBarsOptions `schema:"options"`
}

type ReqBarsOptions struct {
	StartDate string `schema:"options[start_date]"`
	EndDate   string `schema:"options[end_date]"`
	Tf        string `schema:"options[timeframe]"`
	Extended  uint8  `schema:"options[extended]"`
}

type RespBars struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    []BarHist `json:"data"`
}
