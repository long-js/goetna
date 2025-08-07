//go:build test

package goetna

import (
	"context"
	"encoding/base64"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	sch "github.com/long-js/goetna/schema"
)

type TestTable map[string]struct {
	arg    interface{}
	expect map[string]interface{}
}

var rest, ctx = createREST(false)

func loadCreds() ([]byte, []byte) {
	login := []byte(os.Getenv("ETNA_LOGIN"))
	passwd := []byte(os.Getenv("ETNA_PASSWD"))
	bLogin := make([]byte, base64.StdEncoding.EncodedLen(len(login)))
	bPwd := make([]byte, base64.StdEncoding.EncodedLen(len(passwd)))
	base64.StdEncoding.Encode(bLogin, login)
	base64.StdEncoding.Encode(bPwd, passwd)
	return bLogin, bPwd
}

func createREST(isPrivate bool) (*EtnaREST, context.Context) {
	l, p := loadCreds()

	c := context.Background()
	r := NewEtnaREST(os.Getenv("ETNA_KEY"), os.Getenv("ETNA_HIST_TOKEN"), isPrivate, ColouredLogger("REST"))
	if err := r.Authenticate(c, l, p); err != nil {
		panic(err)
	}
	return r, c
}

func TestGetBars(t *testing.T) {
	(*t).Skip()
	var (
		err  error
		bars []sch.BarHist
	)
	tests := TestTable{
		"AAPL_1h": {
			arg: sch.ReqBars{
				Ticker: "AAPL", ExchangeId: 3, Options: sch.ReqBarsOptions{
					StartDate: "2025-07-24", EndDate: "2025-07-24 14:00", Tf: "1h"}},
			expect: map[string]interface{}{"count": 32}},
		"AAPL_1h eq_date": {
			arg: sch.ReqBars{
				Ticker: "AAPL", ExchangeId: 3, Options: sch.ReqBarsOptions{
					StartDate: "2025-07-24", EndDate: "2025-07-24", Tf: "1h"}},
			expect: map[string]interface{}{"count": 16}},
		"AAPL_15m": {
			arg: sch.ReqBars{
				Ticker: "AAPL", ExchangeId: 3, Options: sch.ReqBarsOptions{
					StartDate: "2025-07-24", EndDate: "2025-07-24", Tf: "15m"}},
			expect: map[string]interface{}{"count": 128}},
		"AAPL_1m": {
			arg: sch.ReqBars{
				Ticker: "AAPL", ExchangeId: 3, Options: sch.ReqBarsOptions{
					StartDate: "2025-07-24", EndDate: "2025-07-24", Tf: "1m"}},
			expect: map[string]interface{}{"count": 951}},
	}
	for name, tc := range tests {
		(*t).Run(name, func(t *testing.T) {
			params := tc.arg.(sch.ReqBars)
			cnt := tc.expect["count"].(int)
			if bars, err = rest.GetBars(ctx, &params); err != nil {
				(*t).Error(err)
			} else if len(bars) != cnt {
				(*t).Errorf("wrong response: #%d", len(bars))
			}
		})
	}
}

func TestGetSecurity(t *testing.T) {
	(*t).Skip()
	var (
		err error
		sec sch.Security
	)
	tests := TestTable{
		"TSLA": {
			arg:    "TSLA",
			expect: map[string]interface{}{"symbol": "TSLA", "tickSize": .01}},
		"NVDA": {
			arg:    "NVDA",
			expect: map[string]interface{}{"symbol": "NVDA", "tickSize": .01}},
		"AAPL": {
			arg:    "AAPL",
			expect: map[string]interface{}{"symbol": "AAPL", "tickSize": .01}},
	}
	for name, tc := range tests {
		(*t).Run(name, func(t *testing.T) {
			symbol := tc.arg.(string)
			symb := tc.expect["symbol"].(string)
			tick := tc.expect["tickSize"].(float64)
			if sec, err = rest.GetSecurity(ctx, symbol); err != nil {
				(*t).Error(err)
			} else if sec.Symbol != symb || sec.TickSize != tick {
				(*t).Errorf("wrong response: %+v", sec)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	(*t).Skip()
	if resp, err := rest.GetUser(ctx); err != nil {
		(*t).Error(err)
	} else if resp.UserId == 0 || resp.Login == "" {
		(*t).Errorf("wrong user info: %+v", resp)
	}
}

func TestGetUserSettings(t *testing.T) {
	(*t).Skip()
	if resp, err := rest.GetUserSettings(ctx); err != nil {
		(*t).Error(err)
	} else {
		(*t).Logf("user trading settings: %+v", resp)
	}
}

func TestGetAvailableExchanges(t *testing.T) {
	(*t).Skip()
	if resp, err := rest.GetAvailableExchanges(ctx); err != nil {
		(*t).Error(err)
	} else {
		(*t).Logf("EXCHs: %+v", resp)
	}
}

func TestGetUserAccounts(t *testing.T) {
	// (*t).Skip()
	var (
		err  error
		accs []sch.Account
	)
	if accs, err = rest.GetUserAccounts(ctx); err != nil {
		(*t).Error(err)
	} else {
		(*t).Logf("ACCOUNS: %+v\n", accs)
		// [
		// 	{"Id":292,"ClearingAccount":"292","AccessType":"Owner","MarginType":"Cash",
		// 		"OwnerType":"IndividualCustomer",
		// 		"Currency":"USD"}]
	}
}

func TestGetBalance(t *testing.T) {
	// (*t).Skip()
	var (
		err error
		bal sch.TradingBalance
	)
	if bal, err = rest.GetBalance(ctx, 292); err != nil {
		(*t).Error(err)
	} else {
		(*t).Logf("BAL: %+v\n", bal)
	}
}

func TestGetBalanceHistory(t *testing.T) {
	(*t).Skip()
	var (
		err  error
		bals []sch.BalanceHistoryValue
	)
	fmt := "2006-01-02T15:04:05.999999Z"
	fromTs := time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC).Format(fmt)
	tillTs := time.Now().Format(fmt)
	if bals, err = rest.GetBalanceHistory(ctx, 292, fromTs, tillTs); err != nil {
		(*t).Error(err)
	} else {
		(*t).Logf("BALS: %+v\n", bals)
	}
}

func TestGetTransfers(t *testing.T) {
	(*t).Skip()
	if resp, err := rest.GetTransfers(ctx, 292); err != nil {
		(*t).Error(err)
	} else {
		(*t).Logf("TRANSES: %+v", resp)
	}
}

func TestGetPositions(t *testing.T) {
	// (*t).Skip()
	var (
		err   error
		poses []sch.Position
	)
	if poses, err = rest.GetPositions(ctx, 292); err != nil {
		(*t).Error(err)
	} else {
		(*t).Logf("POSES: %+v\n", poses)
	}
}

func TestGetOrders(t *testing.T) {
	// (*t).Skip()
	var (
		err  error
		ords []sch.Order
	)
	if ords, err = rest.GetOrders(ctx, 292, true); err != nil {
		(*t).Error(err)
	} else if ords, err = rest.GetOrders(ctx, 292, false); err != nil {
		(*t).Error(err)
	} else {
		(*t).Logf("ORDERS: %+v\n", ords)
	}
}

func TestGetOrder(t *testing.T) {
	(*t).Skip()
	var (
		err error
		ord sch.Order
	)
	if ord, err = rest.GetOrder(ctx, 292, 1357); err != nil {
		(*t).Error(err)
	} else {
		(*t).Logf("ORDER: %+v\n", ord)
	}
}

func TestPlaceOrder(t *testing.T) {
	(*t).Skip()
	var (
		err   error
		ord   sch.Order
		AccId = uint32(292)
	)
	tests := TestTable{
		"Limit Buy": {
			arg: sch.OrderParams{
				Quantity: 1., Price: 240., Symbol: "TSLA", ClientId: "LB001", Type: sch.OrderLimit,
				Side: sch.SideBuy, Comment: "TSLA (/!@%^*&#_';$ LB001"},
			expect: map[string]interface{}{
				"symbol": "TSLA", "qty": 1., "side": "Buy", "cid": "LB001", "comment": "TSLA (/!@%^*&#_';$ LB001"}},
		"Limit Sell": {
			arg: sch.OrderParams{
				Quantity: 1., Price: 265., Symbol: "TSLA", ClientId: "LS001", Type: sch.OrderLimit,
				Side: sch.SideSell},
			expect: map[string]interface{}{"symbol": "TSLA", "qty": 1., "side": "Sell", "cid": "LS001", "comment": ""}},
	}
	active := make([]uint64, 0, len(tests))
	for name, tc := range tests {
		(*t).Run(name, func(t *testing.T) {
			params := tc.arg.(sch.OrderParams)
			if ord, err = rest.PlaceOrder(ctx, AccId, &params); err != nil {
				(*t).Error(err)
			} else {
				active = append(active, ord.Id)
			}
			d1 := cmp.Diff(tc.expect["symbol"].(string), ord.Symbol)
			d2 := cmp.Diff(tc.expect["qty"].(float64), ord.Quantity)
			d3 := cmp.Diff(tc.expect["comment"].(string), ord.Comment)
			if d1 != "" || d2 != "" || d3 != "" {
				(*t).Errorf("%s %s %s\n", d1, d2, d3)
			}
		})
	}
	time.Sleep(2 * time.Second)
	for _, oid := range active {
		if err = rest.CancelOrder(ctx, AccId, oid); err != nil {
			(*t).Error(err)
		}
	}
}

func TestGetStreamers(t *testing.T) {
	// (*t).Skip()
	var (
		err  error
		resp sch.Streamers
	)
	if resp, err = rest.GetStreamers(ctx); err != nil {
		(*t).Error(err)
	} else {
		(*t).Logf("STREAMERS: %+v\n", resp)
	}
}

func TestRecoverSession(t *testing.T) {
	(*t).Skip()
	var (
		err          error
		respD, respQ sch.SessionId
	)
	if respD, err = rest.RecoverStreamerSession(ctx, sch.WSSessData); err != nil {
		(*t).Error(err)
	} else if respQ, err = rest.RecoverStreamerSession(ctx, sch.WSSessQuote); err != nil {
		(*t).Error(err)
	} else {
		(*t).Logf("RECOVERED:\ndata: %s\nquote: %s\n", respD, respQ)
	}
}
