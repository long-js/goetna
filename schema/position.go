package schema

import "time"

// Position...
type Position struct {
	Id                   int       `json:"Id"`
	AccountId            int       `json:"AccountId"`
	SecurityId           int       `json:"SecurityId"`
	Symbol               string    `json:"Symbol"`
	Name                 string    `json:"Name"`
	SecurityCurrency     string    `json:"SecurityCurrency"`
	SecurityType         string    `json:"SecurityType"`
	ContractSize         int       `json:"ContractSize"`
	CostBasis            float64   `json:"CostBasis"`
	DailyCostBasis       float64   `json:"DailyCostBasis"`
	CreateDate           time.Time `json:"CreateDate"`
	ModifyDate           time.Time `json:"ModifyDate"`
	Quantity             int       `json:"Quantity"`
	RealizedProfitLoss   int       `json:"RealizedProfitLoss"`
	AverageOpenPrice     float64   `json:"AverageOpenPrice"`
	AverageClosePrice    int       `json:"AverageClosePrice"`
	StopLossPrice        int       `json:"StopLossPrice"`
	TakeProfitPrice      int       `json:"TakeProfitPrice"`
	DailyCloseProfitLoss int       `json:"DailyCloseProfitLoss"`
	ExcessChanges        int       `json:"ExcessChanges"`
	DayQuantity          int       `json:"DayQuantity"`
	MarketValueEOD       float64   `json:"MarketValueEOD"`
}

type RespPositions struct {
	Result           []Position `json:"Result"`
	NextPageLink     string     `json:"NextPageLink"`
	PreviousPageLink string     `json:"PreviousPageLink"`
	TotalCount       uint8      `json:"TotalCount"`
}
