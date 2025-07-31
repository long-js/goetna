package schema

import "time"

// Transfer representation
type Transfer struct {
	Id                    string    `json:"Id"`
	ExternalId            string    `json:"ExternalId"`
	Mechanism             string    `json:"Mechanism"`
	Status                string    `json:"Status"`
	Comment               string    `json:"Comment"`
	ClearingAccountNumber string    `json:"ClearingAccountNumber"`
	Amount                float64   `json:"Amount"`
	TotalAmount           float64   `json:"TotalAmount"`
	TransferDate          time.Time `json:"TransferDate"`
	CreatedAt             time.Time `json:"CreatedAt"`
	AccountId             uint32    `json:"AccountId"`
	IsDeposit             bool      `json:"IsDeposit"`
}

type RespTransfers struct {
	Result           []Transfer `json:"Result"`
	NextPageLink     string     `json:"NextPageLink"`
	PreviousPageLink string     `json:"PreviousPageLink"`
	TotalCount       uint16     `json:"TotalCount"`
}
