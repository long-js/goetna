package schema

import "strconv"

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

// BarSize...
type BarSize struct {
	Name    string
	Seconds uint32
}

type ReqBars struct {
	Security   ReqBarsSecurity `json:"Security"`
	Settings   ReqBarsSettings `json:"SecurityHistorySettings"`
	Indicators []string        `json:"IndicatorsHistorySettings"`
}

type ReqBarsSecurity struct {
	Symbol string `json:"Symbol"`
}

type ReqBarsSettings struct {
	StartDate     int64  `json:"StartDate,omitempty"`
	EndDate       int64  `json:"EndDate,omitempty"`
	Count         int16  `json:"CandlesCount,omitempty"`
	Interval      int8   `json:"Interval"`
	IncludeNonRTH bool   `json:"IncludeNonMarketData,omitempty"`
	Period        string `json:"Period"`
}

type RespBars struct {
	Bars       [][]Bar `json:"SecurityHistory"`
	Indicators [][]struct {
		Date   int64
		Values []float64
	} `json:"IndicatorsHistory"`
}
