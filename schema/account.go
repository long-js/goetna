package schema

import "time"

type Account struct {
	Id               int    `json:"Id"`
	ClearingAccount  string `json:"ClearingAccount"`
	AccessType       string `json:"AccessType"`
	MarginType       string `json:"MarginType"`
	OwnerType        string `json:"OwnerType"`
	Enabled          bool   `json:"Enabled"`
	ClearingFirm     string `json:"ClearingFirm"`
	IsAverageAccount bool   `json:"IsAverageAccount"`
	Owners           []struct {
		UserId     int       `json:"UserId"`
		FirstName  string    `json:"FirstName"`
		MiddleName string    `json:"MiddleName"`
		LastName   string    `json:"LastName"`
		Login      string    `json:"Login"`
		Email      string    `json:"Email"`
		Role       int       `json:"Role"`
		AddedDate  time.Time `json:"AddedDate"`
		Salutation string    `json:"Salutation"`
		Suffix     string    `json:"Suffix"`
	} `json:"Owners"`
}

type TradingBalance struct {
	Cash    float64 `json:"cash"`    // The amount of funds available on the trading account.
	NetCash float64 `json:"netCash"` // The amount of funds available on the account minus the options margin requirement.

	// The amount of funds that can either be withdrawn or used to open new positions.
	// This value is used as the basis for calculating buying power.
	// Excess = Equity - TMMR - Pending Cash - Uncleared Cash, where TMMR - Total Maintenance Margin Requirement
	Excess float64 `json:"excess"`

	ChangeAbsolute          float64 `json:"changeAbsolute"`          // The difference between the account's value and the account's equity at the closing of the previous trading session.
	ChangePercent           float64 `json:"changePercent"`           // Identical to changeAbsolute but expressed in percentage terms.
	EquityTotal             float64 `json:"equityTotal"`             // The gross valuation of all equity on the trading account.
	PendingOrdersCount      float64 `json:"pendingOrdersCount"`      // The number of pending orders on the account.
	NetLiquidity            float64 `json:"netLiquidity"`            // The amount of funds will be available after all active positions are terminated.
	StockLongMarketValue    float64 `json:"stockLongMarketValue"`    // The gross market value of all long stock positions.
	StockShortMarketValue   float64 `json:"stockShortMarketValue"`   // The gross market value of all short stock positions.
	OptionLongMarketValue   float64 `json:"optionLongMarketValue"`   // The gross market value of all long option positions.
	OptionShortMarketValue  float64 `json:"optionShortMarketValue"`  // The gross market value of all short option positions.
	ForexLongMarketValue    float64 `json:"forexLongMarketValue"`    // The gross market value of all long forex positions.
	ForexShortMarketValue   float64 `json:"forexShortMarketValue"`   // The gross market value of all short forex positions.
	DayTrades               float64 `json:"dayTrades"`               // The number of day trades that have been executed during the last five trading sessions.
	StockBuyingPower        float64 `json:"stockBuyingPower"`        // The gross amount of stocks that can be purchased, adjusted for the available margin debt.
	OptionBuyingPower       float64 `json:"optionBuyingPower"`       // The gross amount of options that can be purchased, adjusted for the available margin debt.
	ForexBuyingPower        float64 `json:"forexBuyingPower"`        // The gross amount of forex positions that can be opened, adjusted for the available margin debt.
	PendingCash             float64 `json:"pendingCash"`             // The amount of funds reserved to complete pending transactions.
	MaintenanceMargin       float64 `json:"maintenanceMargin"`       // The minimum amount of equity that must be maintained in a margin account.
	OptionMaintenanceMargin float64 `json:"optionMaintenanceMargin"` // The minimum amount of equity that must be maintained for option securities.
	OpenPL                  float64 `json:"openPL"`                  // The amount of unrealized profit or loss for all positions.
	ClosePL                 float64 `json:"closePL"`                 // The amount of realized profit or loss during the current trading session.
	MarketValue             float64 `json:"marketValue"`             // The market value of all open long and short positions.
	TotalPL                 float64 `json:"totalPL"`
}

type TradingBalanceValue struct {
	Date  time.Time `json:"Date"`
	Value float64   `json:"Value"`
}
