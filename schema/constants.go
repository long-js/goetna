package schema

const (
	ShowModeShow    ShowMode = "Show"
	ShowModeNotShow          = "DoNotShow"
	ShowModeIfError          = "ShowIfAnError"
)
const (
	OrderMarket    OrderType = "Market"
	OrderLimit               = "Limit"
	OrderStop                = "Stop"
	OrderStopLimit           = "Stop Limit"
)
const (
	SideBuy        OrderSide = "Buy"
	SideSell                 = "Sell"
	SideSellShort            = "SellShort"
	SideBuyToCover           = "BuyToCover"
)
const (
	StatusNew = OrderStatus(iota)
	StatusPartial
	StatusFilled
	StatusDoneDay
	StatusCanceled
	StatusReplaced
	StatusPendingCancel
	StatusStopped
	StatusRejected
	StatusSuspended
	StatusPending
	StatusCalculated
	StatusExpired
	StatusAcceptedBidding
	StatusPendingReplace
	StatusError
)
const (
	// Day. The order automatically expires at the end of the regular trading session.
	TimeInForceDay TimeInForce = "Day"

	// GoodTillCancel. The order persists indefinitely until it is executed or manually cancelled.
	TimeInForceGTC TimeInForce = "GoodTillCancel"

	// GoodTillDate. The order will be active until the date specified in the ExpireDate attribute.
	TimeInForceGTD TimeInForce = "GoodTillDate"
)
const (
	SessReg     TradingSession = "REG"     // regular sessions
	SessPre     TradingSession = "PRE"     // pre-market session
	SessPost    TradingSession = "POST"    // post-market session
	SessRegPost TradingSession = "REGPOST" // regular and post-market sessions
	SessAll     TradingSession = "ALL"     // pre-market, regular and post-market sessions
)

var VALID_TFS = map[string]BarSize{
	"1m":  {"1", 60},
	"2m":  {"2", 120},
	"3m":  {"3", 180},
	"5m":  {"5", 300},
	"10m": {"10", 600},
	"15m": {"15", 900},
	"30m": {"30", 1800},
	"1h":  {"60", 3600},
	"2h":  {"120", 7200},
	"1D":  {"1440", 86400},
}

/*
	WebSocket constants
*/

const (
	WSReconnInterval  = 12 // base period in seconds for the reconnect period calculation
	WSMaxSilentPeriod = 30 // maximum period of silence, seconds
)
const (
	WSSessData  WSSessionType = 0 // data session
	WSSessQuote WSSessionType = 1 // quote session
)
const (
	WSTopicQuote    = "Quote"
	WSTopicCandle   = "Candle"
	WSTopicBalance  = "AccountBalance"
	WSTopicPosition = "Position"
	WSTopicOrder    = "Order"
	WSCmdCreate     = "CreateSession.txt"
	WSCmdSub        = "Subscribe.txt"
	WSCmdUnsub      = "Unsubscribe.txt"
	WSCmdPing       = "Ping"
)

const (
	FieldEntytyType = "\"EntityType\": "
	FieldCmd        = "\"Cmd\": "
)

var WSPongMsg = []byte("{\"Cmd\":\"Pong\",\"StatusCode\":\"Ok\"}")
