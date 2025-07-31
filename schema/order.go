package schema

import "time"

type Order struct {
	Id                      int64          `json:"Id"`
	Quantity                float64        `json:"Quantity"`
	Price                   float64        `json:"Price"`
	StopPrice               float64        `json:"StopPrice"` // When this price is reached, the order will automatically be converted into a market order.
	ExecutedQuantity        float64        `json:"ExecutedQuantity"`
	LastPrice               float64        `json:"LastPrice"`
	LastQuantity            float64        `json:"LastQuantity"`
	LeavesQuantity          float64        `json:"LeavesQuantity"`
	AveragePrice            float64        `json:"AveragePrice"`
	BrokerServiceCommission float64        `json:"BrokerServiceCommission"`
	AccountId               int64          `json:"AccountId"`
	UserId                  int64          `json:"UserId"`
	RequestId               int64          `json:"RequestId"`
	ParentRequestId         int64          `json:"ParentRequestId"`
	StateId                 int64          `json:"StateId"`
	ParentId                int64          `json:"ParentId"`
	Date                    time.Time      `json:"Date"`
	TransactionDate         time.Time      `json:"TransactionDate"`
	ExpireDate              time.Time      `json:"ExpireDate"` // The expiration of the order.
	Symbol                  string         `json:"Symbol"`     // The ticker symbol of the underlying security in the new order.
	Currency                string         `json:"Currency"`
	ClientId                string         `json:"ClientId"` // The order ID on the client's side.
	Side                    OrderSide      `json:"Side"`
	Status                  string         `json:"Status"`
	ExecutionStatus         string         `json:"ExecutionStatus"`
	Type                    OrderType      `json:"Type"`
	RequestStatus           string         `json:"RequestStatus"`
	Target                  string         `json:"Target"`
	Comment                 string         `json:"Comment"`
	Description             string         `json:"Description,omitempty"`
	TimeInforce             TimeInForce    `json:"TimeInForce"`
	ClearingAccount         string         `json:"ClearingAccount"`
	ExecInst                string         `json:"ExecInst"` // Indicates if the order should be filled either entirely in one transaction or not at all. Possible value: 'AllOrNone'.
	Exchange                string         `json:"Exchange"` // The exchange on which the order should be executed.
	ExecutionVenue          string         `json:"ExecutionVenue"`
	InitialType             string         `json:"InitialType"`
	ExtendedHours           TradingSession `json:"ExtendedHours"` // If the order should be placed during the extended hours. (PRE, REG, REGPOST)
	ExecBrocker             string         `json:"ExecBrocker"`
	TransType               string         `json:"TransType"`
	ExecId                  string         `json:"ExecId"`
	QuantityQualifier       string         `json:"QuantityQualifier"`
	ExecutionInstructions   struct {
		Solicited          string `json:"Solicited"`
		PositionMarginType string `json:"PositionMarginType"`
		IP                 string `json:"IP"`
	} `json:"ExecutionInstructions"`
	IsExternal bool `json:"IsExternal"`
}

type RespOrders struct {
	Result           []Order `json:"Result"`
	NextPageLink     string  `json:"NextPageLink"`
	PreviousPageLink string  `json:"PreviousPageLink"`
	TotalCount       uint8   `json:"TotalCount"`
}

type OrderParams struct {
	ParentId  int64   `json:"ParentId,omitempty"`
	Quantity  float64 `json:"Quantity"` // The number of shares in the order.
	Price     float64 `json:"Price,omitempty"`
	StopPrice float64 `json:"StopPrice,omitempty"` // The order'll be converted automatically into a market order, when this price is reached.
	// Expiration            time.Time       `json:"ExpireDate,omitempty"` // The order'll be cancelled automatically after the specified time.
	Symbol                string            `json:"Symbol"`
	ClientId              string            `json:"ClientId,omitempty"` // The order ID on the client's side.
	BotName               string            `json:"-"`
	Type                  OrderType         `json:"Type"` // Possible values are: Market, Limit, Stop, Stop Limit.
	Side                  OrderSide         `json:"Side"` // Possible values are: Buy, Sell, SellShort, BuyToCover.
	Comment               string            `json:"Comment,omitempty"`
	ExecInst              string            `json:"ExecInst,omitempty"` // The possible value: AllOrNone.
	TimeInforce           TimeInForce       `json:"TimeInforce"`        // The period in which the order will be active.
	Exchange              string            `json:"Exchange,omitempty"` // The exchange, on which the order should be executed.
	ExtendedHours         TradingSession    `json:"ExtendedHours"`      // If the order should be placed during the extended hours (pre-market, post-market).
	ExecutionInstructions *ExecInstructions `json:"ExecutionInstructions,omitempty"`
	ValidationsToBypass   uint8             `json:"ValidationsToBypass,omitempty"`
}

type ExecInstructions struct {
	PerTradeCommission    string `json:"PerTradeCommission"`    // Specified in dollars ($1 per trade).
	PerContractCommission string `json:"PerContractCommission"` // Specified in cents (1 cent per contract).
}

type TimeInForce string
type OrderStatus uint8
type OrderSide string
type OrderType string
type ShowMode string
