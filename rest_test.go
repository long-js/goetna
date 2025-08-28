//go:build test

package goetna

import (
	"context"
	"encoding/base64"
	"os"
	"strings"
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
	r, err := NewEtnaREST(os.Getenv("ETNA_KEY"), os.Getenv("ETNA_NRTH_TOKEN"), l, p, isPrivate, ColouredLogger("REST"))
	if err != nil {
		panic(err)
	}
	return r, c
}

/*
	Public
*/

func TestGetBars(t *testing.T) {
	// (*t).Skip()
	var (
		err  error
		bars []sch.BarHist
	)
	today := time.Now().In(time.UTC).Format("2006-01-02")
	today_tm := time.Now().In(time.UTC).Format("2006-01-02 15:04")
	// yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")
	tests := TestTable{
		// "AAPL_1h": {
		// 	arg: sch.ReqBars{
		// 		Ticker: "AAPL", ExchangeId: 3, Options: sch.ReqBarsOptions{
		// 			StartDate: yesterday, EndDate: today_tm, Tf: "1h"}},
		// 	expect: map[string]interface{}{"count": 16}},
		// "AAPL_1h eq_date": {
		// 	arg: sch.ReqBars{
		// 		Ticker: "AAPL", ExchangeId: 3, Options: sch.ReqBarsOptions{
		// 			StartDate: today, EndDate: today, Tf: "1h"}},
		// 	expect: map[string]interface{}{"count": 16}},
		// "AAPL_15m": {
		// 	arg: sch.ReqBars{
		// 		Ticker: "AAPL", ExchangeId: 3, Options: sch.ReqBarsOptions{
		// 			StartDate: yesterday, EndDate: today_tm, Tf: "15m"}},
		// 	expect: map[string]interface{}{"count": 64}},
		"AAPL_1m": {
			arg: sch.ReqBars{
				Ticker: "AAPL", ExchangeId: 3, Options: sch.ReqBarsOptions{
					StartDate: today, EndDate: today_tm, Tf: "1m"}},
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

func TestGetAvailableExchanges(t *testing.T) {
	(*t).Skip()
	if resp, err := rest.GetAvailableExchanges(ctx); err != nil {
		(*t).Error(err)
	} else {
		(*t).Logf("EXCHs: %+v", resp)
	}
}

func TestGetStreamers(t *testing.T) {
	(*t).Skip()
	var (
		err  error
		resp sch.Streamers
	)

	if resp, err = rest.GetStreamers(ctx, true); err != nil {
		(*t).Error(err)
	} else if len(resp.QuoteAddresses) == 0 || resp.QuoteAddresses[0].Url != DefaultConfig.WSUrlPubFMP {
		(*t).Errorf("wrong streamers: %+v", resp.QuoteAddresses)
	} else {
		(*t).Logf("STREAMERS: %+v\n", resp)
	}
	if resp, err = rest.GetStreamers(ctx, false); err != nil {
		(*t).Error(err)
	} else if len(resp.QuoteAddresses) == 0 || strings.TrimSuffix(resp.QuoteAddresses[0].Url, ":443") != DefaultConfig.WSUrlPub {
		(*t).Errorf("wrong streamers: %+v", resp.QuoteAddresses)
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

/*
	Private
*/

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

func TestGetUserAccounts(t *testing.T) {
	(*t).Skip()
	var (
		err  error
		accs []sch.Account
	)
	if accs, err = rest.GetUserAccounts(ctx); err != nil {
		(*t).Error(err)
	} else {
		(*t).Logf("ACCOUNS: %+v\n", accs)
	}
}

func TestGetBalance(t *testing.T) {
	// (*t).Skip()
	var (
		err error
		bal sch.TradingBalance
	)
	accId := uint32(421) // 292
	if bal, err = rest.GetBalance(ctx, accId); err != nil {
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
	accId := uint32(421) // 292
	fmt := "2006-01-02T15:04:05.999999Z"
	fromTs := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC).Format(fmt)
	tillTs := time.Now().Format(fmt)
	if bals, err = rest.GetBalanceHistory(ctx, accId, fromTs, tillTs); err != nil {
		(*t).Error(err)
	} else {
		(*t).Logf("BALS: %+v\n", bals)
	}
}

func TestGetTransfers(t *testing.T) {
	(*t).Skip()
	accId := uint32(421) // 292
	if resp, err := rest.GetTransfers(ctx, accId); err != nil {
		(*t).Error(err)
	} else {
		(*t).Logf("TRANSES: %+v", resp)
	}
}

func TestGetPositions(t *testing.T) {
	(*t).Skip()
	var (
		err   error
		poses []sch.Position
		accId = uint32(421) // 292
	)
	if poses, err = rest.GetPositions(ctx, accId); err != nil {
		(*t).Error(err)
	} else {
		(*t).Logf("POSES: %+v\n", poses)
	}
}

func TestGetOrders(t *testing.T) {
	(*t).Skip()
	var (
		err   error
		ords  []sch.Order
		accId = uint32(421) // 292
	)
	if ords, err = rest.GetOrders(ctx, accId, true); err != nil {
		(*t).Error(err)
	} else if ords, err = rest.GetOrders(ctx, accId, false); err != nil {
		(*t).Error(err)
	} else {
		for _, o := range ords {
			(*t).Logf("ORDERS: %+v\n", o)
		}
	}
}

func TestGetOrder(t *testing.T) {
	(*t).Skip()
	var (
		err   error
		ord   sch.Order
		accId = uint32(421) // 292
	)
	if ord, err = rest.GetOrder(ctx, accId, 1357); err != nil {
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
		accId = uint32(421) // 292
	)
	tests := TestTable{
		"Limit Buy": {
			arg: sch.OrderParams{
				Quantity: 1., Price: 310., Symbol: "TSLA", ClientId: "LB001", Type: sch.OrderLimit,
				Side: sch.SideBuy, Comment: "TSLA (/!@%^*&#_';$ LB001"},
			expect: map[string]interface{}{
				"symbol": "TSLA", "qty": 1., "side": "Buy", "cid": "LB001", "comment": "TSLA (/!@%^*&#_';$ LB001"}},
		"Limit Sell": {
			arg: sch.OrderParams{
				Quantity: 1., Price: 340., Symbol: "TSLA", ClientId: "LS001", Type: sch.OrderLimit,
				Side: sch.SideSell},
			expect: map[string]interface{}{"symbol": "TSLA", "qty": 1., "side": "Sell", "cid": "LS001", "comment": ""}},
	}
	active := make([]uint64, 0, len(tests))
	for name, tc := range tests {
		(*t).Run(name, func(t *testing.T) {
			params := tc.arg.(sch.OrderParams)
			if ord, err = rest.PlaceOrder(ctx, accId, &params); err != nil {
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
		if err = rest.CancelOrder(ctx, accId, oid); err != nil {
			(*t).Error(err)
		}
	}
}

func TestCancelOrder(t *testing.T) {
	(*t).Skip()
	var (
		err   error
		accId = uint32(421) // 292
	)
	if err = rest.CancelOrder(ctx, accId, 1662); err != nil {
		(*t).Error(err)
	}
}
