package schema

// Bar...
type Bar struct {
	Open     float64 `json:"Open"`
	High     float64 `json:"High"`
	Low      float64 `json:"Low"`
	Close    float64 `json:"Close"`
	Volume   float64 `json:"Volume"`
	Time     uint32  `json:"Time"`
	IsMarket bool    `json:"IsMarket"` // Indicates if the bar is positioned during the RTH (true).
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
