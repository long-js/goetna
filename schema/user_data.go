package schema

type UserInfo struct {
	UserId     int32  `json:"UserId"`               // the internal ID of the user in ETNA Trader
	FirstName  string `json:"FirstName,omitempty"`  // the first name of the user
	MiddleName string `json:"MiddleName,omitempty"` // the middle name of the user
	LastName   string `json:"LastName,omitempty"`   // the last name of the user
	Login      string `json:"Login"`                // the user's login in ETNA Trader
	Email      string `json:"Email"`                // the email address of the user in ETNA Trader
	AddedDate  string `json:"AddedDate"`            // the date on which this user account was added to ETNA Trader
}

type UserTradingSettings struct {
	Instruments map[string]SecTypeSettings
	// The step by which the number of securities must be increased when placing an order.
	QuantityStepIncrementMultiplier int32
	// The step by which the limit or stop price must be increased when placing an order.
	PriceStepIncrementMultiplier int32
	// Whether the order verification view should be displayed when placing an order.
	SkipVerifyOrder ShowMode
	// Whether the order verification view should be displayed when cancelling an order.
	SkipVerifyCancelOrder ShowMode
	// Whether the order verification view should be displayed when closing an existing position.
	SkipVerifyClosingPosition ShowMode
	// Whether the order verification view should be displayed when replacing an order.
	SkipVerifyOrderReplace ShowMode
	// Whether the order status view should be displayed after an order has been placed.
	SkipPlaceOrderStatus ShowMode
	// Whether the order status view should be displayed after an order has been cancelled.
	SkipCancelOrderStatus ShowMode
	// Whether the order status view should be displayed after a positions has been closed.
	SkipClosingPositionStatus ShowMode
	// Whether the order status view should be displayed after an order has been replaced.
	SkipOrderReplaceStatus ShowMode
	// The maximum number of securities that can be traded in a single stock order.
	MaxStocksQuantity int64
	// The maximum number of securities that can be traded in a single option order.
	MaxOptionsQuantity int64
}

type SecTypeSettings struct {
	OrderType           // The default order type that is set whenever a security of the specified type is traded.
	Quantity     int32  // The default number of securities of an order for the specified security type.
	DurationType string // The default duration of an order for the specified security type.
	ExchangeType string // The default execution venue for the specified security type.
	AON          bool   // Indicates whether orders should be All-Or-None by default.
}

type ReqUserRegister struct {
	Credentials struct {
		Login    string `json:"Login"`
		Email    string `json:"Email"`
		Password string `json:"Password"`
	} `json:"Credentials"`
	Name struct {
		FirstName string `json:"FirstName"`
		LastName  string `json:"LastName"`
	} `json:"Name"`
}
