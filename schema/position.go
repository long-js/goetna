package schema

import "time"

// Position...
type Position struct {
	AccountId            uint32    `json:"AccountId"`
	Id                   uint32    `json:"Id"`
	SecurityId           uint32    `json:"SecurityId"`
	MinContractSize      uint32    `json:"ContractSize"`
	Quantity             int64     `json:"Quantity"`
	RealizedProfitLoss   float64   `json:"RealizedProfitLoss"`
	AverageClosePrice    float64   `json:"AverageClosePrice"`
	StopLossPrice        float64   `json:"StopLossPrice"`
	TakeProfitPrice      float64   `json:"TakeProfitPrice"`
	DailyCloseProfitLoss float64   `json:"DailyCloseProfitLoss"`
	ExcessChanges        float64   `json:"ExcessChanges"`
	CostBasis            float64   `json:"CostBasis"`
	DailyCostBasis       float64   `json:"DailyCostBasis"`
	AverageOpenPrice     float64   `json:"AverageOpenPrice"`
	MarketValueEOD       float64   `json:"MarketValueEOD"`
	DayQuantity          int       `json:"DayQuantity"`
	Symbol               string    `json:"Symbol"`
	Name                 string    `json:"Name"`
	SecurityCurrency     string    `json:"SecurityCurrency"`
	SecurityType         string    `json:"SecurityType"`
	CreateDate           time.Time `json:"CreateDate"`
	ModifyDate           time.Time `json:"ModifyDate"`
}

type RespPositions struct {
	Result           []Position `json:"Result"`
	NextPageLink     string     `json:"NextPageLink"`
	PreviousPageLink string     `json:"PreviousPageLink"`
	TotalCount       uint8      `json:"TotalCount"`
}
