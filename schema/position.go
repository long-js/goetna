package schema

import (
	"strconv"
	"time"
)

// Position...
type Position struct {
	Quantity           int64   `json:"Quantity"`
	RealizedProfitLoss float64 `json:"RealizedProfitLoss"`
	StopLossPrice      float64 `json:"StopLossPrice"`
	TakeProfitPrice    float64 `json:"TakeProfitPrice"`
	CostBasis          float64 `json:"CostBasis"`
	AverageOpenPrice   float64 `json:"AverageOpenPrice"`
	MinContractSize    float64 `json:"ContractSize"`
	AccountId          uint32  `json:"AccountId"`
	Id                 uint32  `json:"Id"`
	SecurityId         uint32  `json:"SecurityId"`
	Symbol             string  `json:"Symbol"`
	Exchange           string
	SecurityCurrency   string    `json:"SecurityCurrency"`
	SecurityType       string    `json:"SecurityType"`
	CreateDate         time.Time `json:"CreateDate"`
	ModifyDate         time.Time `json:"ModifyDate"`
}

func (p *Position) Parse(values map[string]string) error {
	var err error

	for k, v := range values {
		switch k {
		case "Id":
			var pid64 uint64
			if pid64, err = strconv.ParseUint(v, 10, 32); err == nil {
				(*p).Id = uint32(pid64)
			}
		case "AccountId":
			var pid64 uint64
			if pid64, err = strconv.ParseUint(v, 10, 32); err == nil {
				(*p).AccountId = uint32(pid64)
			}
		case "SecurityId":
			var pid64 uint64
			if pid64, err = strconv.ParseUint(v, 10, 32); err == nil {
				(*p).SecurityId = uint32(pid64)
			}
		case "ContractSize":
			(*p).MinContractSize, err = strconv.ParseFloat(v, 64)
		case "Quantity":
			(*p).Quantity, err = strconv.ParseInt(v, 10, 64)
		case "RealizedProfitLoss":
			(*p).RealizedProfitLoss, err = strconv.ParseFloat(v, 64)
		case "CostBasis":
			(*p).CostBasis, err = strconv.ParseFloat(v, 64)
		case "AverageOpenPrice":
			(*p).AverageOpenPrice, err = strconv.ParseFloat(v, 64)
		case "Symbol":
			(*p).Symbol = v
		case "Exchange":
			(*p).Exchange = v
		case "Currency":
			(*p).SecurityCurrency = v
		case "SecurityType":
			(*p).SecurityType = v
		case "CreateDate":
			var ms int64
			if ms, err = strconv.ParseInt(v, 10, 64); err == nil {
				(*p).CreateDate = time.UnixMilli(ms)
			}
		case "ModifyDate":
			var ms int64
			if ms, err = strconv.ParseInt(v, 10, 64); err == nil {
				(*p).ModifyDate = time.UnixMilli(ms)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

type RespPositions struct {
	Result           []Position `json:"Result"`
	NextPageLink     string     `json:"NextPageLink"`
	PreviousPageLink string     `json:"PreviousPageLink"`
	TotalCount       uint8      `json:"TotalCount"`
}
