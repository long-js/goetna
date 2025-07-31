package goetna

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"

	gjson "github.com/goccy/go-json"
	gschema "github.com/gorilla/schema"
	sch "github.com/khokhlomin/goetna/schema"
)

func NewEtnaREST(apiKey string, isPrivate bool) *EtnaREST {
	rest := EtnaREST{httpClient: &http.Client{Timeout: 12000000000}, enc: gschema.NewEncoder()}
	if isPrivate {
		rest.url = defaultConfig.RestUrlPriv
	} else {
		rest.url = defaultConfig.RestUrlPub
	}
	header := make(http.Header)
	header["User-Agent"] = []string{"qant-backend/2.0"}
	header["Content-Type"] = []string{"application/json"}
	header["Accept"] = []string{"application/json"}
	header["Connection"] = []string{"keep-alive"}
	header["Et-App-Key"] = []string{apiKey}
	rest.restHeader = header
	return &rest
}

type EtnaREST struct {
	httpClient *http.Client
	restHeader http.Header
	enc        *gschema.Encoder
	url        string
}

func (api *EtnaREST) callAPI(ctx context.Context, method, endpoint string, query url.Values,
	data, result interface{}) error {
	var (
		err    error
		bData  []byte
		req    *http.Request
		resp   *http.Response
		buffer *bytes.Buffer
		uri    string
	)
	// query
	if query != nil {
		uri = fmt.Sprintf("%s%s?%s", (*api).url, endpoint, query.Encode())
	} else {
		uri = fmt.Sprintf("%s%s", (*api).url, endpoint)
	}
	// body
	if data != nil {
		if bData, err = gjson.Marshal(data); err != nil {
			return err
		}
		buffer = bytes.NewBuffer(bData)
	} else {
		buffer = &bytes.Buffer{}
	}
	if req, err = http.NewRequestWithContext(ctx, method, uri, buffer); err != nil {
		return fmt.Errorf("request creation fault: %w", err)
	} else if len(bData) > 0 {
		(*api).restHeader["Content-Length"] = []string{fmt.Sprintf("%d", len(bData))}
	}
	(*req).Header = (*api).restHeader

	fmt.Printf("--> %s %s %s\n", method, endpoint, bData)
	resp, err = (*api).httpClient.Do(req)
	defer func() {
		if resp != nil {
			if err = resp.Body.Close(); err != nil {
				fmt.Printf("can't close response body %+v\n", err)
			}
		}
	}()

	if err != nil {
		return fmt.Errorf("request fault: %v", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		if err = (*api).readBody(resp, result); err != nil {
			return err
		}
	case http.StatusNoContent:
		return nil
	case http.StatusUnauthorized:
		return fmt.Errorf("%s", resp.Status)
	case http.StatusBadRequest:
		if err = (*api).readBody(resp, result); err != nil {
			return err
		}
		return fmt.Errorf("400, BAD_REQUEST: %s", result)
	default:
		if resp.StatusCode >= 500 {
			return fmt.Errorf("%d, server error: %s", resp.StatusCode, resp.Status)
		}
		return fmt.Errorf("invalid request: %s", resp.Status)
	}
	return nil
}

// readBody reads the response body from the provided http.Response, attempts to unmarshal it
// into the given result interface and returns an error if reading or unmarshaling fails.
func (api *EtnaREST) readBody(resp *http.Response, result interface{}) error {
	var (
		err error
		buf []byte
	)
	if buf, err = io.ReadAll(resp.Body); err != nil || len(buf) == 1 {
		return fmt.Errorf("error reading v2 response body: %w", err)
	}
	fmt.Printf("REST: %s\n", buf)
	if err = gjson.Unmarshal(buf, result); err != nil {
		if res, ok := result.(*sch.Response); ok {
			return fmt.Errorf("API error: %v", res)
		}
		return fmt.Errorf("can't unmarshal: %w", err)
	}
	return nil
}

// authenticate performs the authentication process against the API.
// It sets the Username and Password headers and calls the "token" API endpoint.
// If the authentication fails (either due to an API error or the SFA state not being "Succeeded"),
// it returns an error.
func (api *EtnaREST) Authenticate(ctx context.Context, login, passwd []byte) error {
	var (
		err      error
		sfa      sch.SFA
		decLogin = make([]byte, base64.StdEncoding.DecodedLen(len(login)))
		decPwd   = make([]byte, base64.StdEncoding.DecodedLen(len(passwd)))
	)
	if _, err = base64.StdEncoding.Decode(decLogin, login); err != nil {
		return fmt.Errorf("wrong username %+v", err)
	} else if _, err = base64.StdEncoding.Decode(decPwd, passwd); err != nil {
		return fmt.Errorf("wrong password %+v", err)
	}
	(*api).restHeader["Username"] = []string{string(bytes.Trim(decLogin, "\x00"))}
	(*api).restHeader["Password"] = []string{string(bytes.Trim(decPwd, "\x00"))}
	err = (*api).callAPI(ctx, http.MethodPost, "token", nil, nil, &sfa)
	if err != nil || sfa.State != "Succeeded" {
		return fmt.Errorf("authentication failed: %s %s, %s, %+v", sfa.State, sfa.Step, sfa.Reason, err)
	}
	(*api).restHeader.Del("Username")
	(*api).restHeader.Del("Password")
	(*api).restHeader["Authorization"] = []string{"Bearer " + sfa.Token}
	return nil
}

// GetStreamers retrieves a list of streamers from the API.
func (api *EtnaREST) GetStreamers(ctx context.Context) (sch.Streamers, error) {
	var resp sch.Streamers

	err := (*api).callAPI(ctx, http.MethodGet, "v1.0/streamers", nil, nil, &resp)
	if err != nil {
		return resp, fmt.Errorf("getStreamers failed: %+v", err)
	}
	return resp, nil
}

// RecoverStreamerSession attempts to recover a streamer session of the specified type.
func (api *EtnaREST) RecoverStreamerSession(ctx context.Context, sessType sch.WSSessionType) (sch.SessionId, error) {
	var resp sch.SessionResp

	qry := url.Values{"sessionType": []string{fmt.Sprintf("%d", sessType)}}
	err := (*api).callAPI(ctx, http.MethodPut, "v1.0/streamers/session/recover", qry, nil, &resp)
	if err != nil {
		return resp.Id, fmt.Errorf("recoverStreamerSession failed: %+v", err)
	}
	return resp.Id, nil
}

/*
 * Users
 */

func (api *EtnaREST) RegisterUser(ctx context.Context, params *sch.ReqUserRegister) (sch.UserInfo, error) {
	var resp sch.UserInfo
	return resp, nil
}

// GetUser retrieves the authenticated user's information.
func (api *EtnaREST) GetUser(ctx context.Context) (sch.UserInfo, error) {
	var resp sch.UserInfo
	if err := (*api).callAPI(ctx, http.MethodGet, "v1.0/users/@me/info", nil, nil, &resp); err != nil {
		return resp, fmt.Errorf("getUser failed: %+v", err)
	}
	return resp, nil
}

// GetUserSettings retrieves the authenticated user's trading settings.
func (api *EtnaREST) GetUserSettings(ctx context.Context) (sch.UserTradingSettings, error) {
	var resp sch.UserTradingSettings
	if err := (*api).callAPI(ctx, http.MethodGet, "v1.0/users/@me/settings/trading", nil, nil, &resp); err != nil {
		return resp, fmt.Errorf("getUserSettings failed: %+v", err)
	}
	return resp, nil
}

func (api *EtnaREST) ModifyUser(ctx context.Context) error {
	return nil
}

func (api *EtnaREST) UpdateUserPasswd(ctx context.Context) error {
	return nil
}

/*
 * Accounts, balances, positions
 */

// GetAvailableExchanges retrieves the list of user's available exchanges.
func (api *EtnaREST) GetAvailableExchanges(ctx context.Context) ([]string, error) {
	var resp = make([]string, 0, 5)
	err := (*api).callAPI(ctx, http.MethodGet, "v1.0/users/@me/exchanges", nil, nil, &resp)
	if err != nil {
		return resp, fmt.Errorf("getAvailableExchanges failed: %+v", err)
	}
	return resp, nil
}

// GetUserAccounts retrieves a slice of user accounts for the authenticated user.
func (api *EtnaREST) GetUserAccounts(ctx context.Context) ([]sch.Account, error) {
	var resp []sch.Account
	if err := (*api).callAPI(ctx, http.MethodGet, fmt.Sprintf("v1.0/users/@me/accounts"), nil, nil, &resp); err != nil {
		return nil, fmt.Errorf("getAllAccounts failed: %+v", err)
	}
	return resp, nil
}

// GetBalance retrieves a trading balancedata for the authenticated user.
func (api *EtnaREST) GetBalance(ctx context.Context, accId uint32) (sch.TradingBalance, error) {
	var resp sch.TradingBalance

	err := (*api).callAPI(ctx, http.MethodGet, fmt.Sprintf("v1.0/accounts/%d/info", accId), nil, nil, &resp)
	if err != nil {
		return resp, fmt.Errorf("getBalance failed: %+v", err)
	}
	return resp, nil
}

// GetBalanceHistory retrieves a slice of sch.TradingBalanceValue for the specified account of authenticated user.
func (api *EtnaREST) GetBalanceHistory(ctx context.Context, accId uint32,
	fromTs, tillTs string) ([]sch.BalanceHistoryValue, error) {
	var resp []sch.BalanceHistoryValue

	qry := url.Values{"startDate": {fromTs}, "endDate": {tillTs}, "step": {"1"}}
	err := (*api).callAPI(ctx, http.MethodGet, fmt.Sprintf("v1.0/accounts/%d/history", accId), qry, nil, &resp)
	if err != nil {
		return resp, fmt.Errorf("getBalanceHistory failed: %+v", err)
	}
	return resp, nil
}

// GetPositions retrieves a slice of positions for the authenticated user.
func (api *EtnaREST) GetPositions(ctx context.Context, accId uint32) ([]sch.Position, error) {
	var resp sch.RespPositions

	qry := url.Values{"pageNumber": {"0"}, "pageSize": {"99"}, "sortField": {"Symbol"}, "desc": {"false"}}
	err := (*api).callAPI(ctx, http.MethodGet, fmt.Sprintf("v1.0/accounts/%d/positions", accId), qry, nil, &resp)
	if err != nil {
		return resp.Result, fmt.Errorf("getPositions failed: %+v", err)
	}
	return resp.Result, nil
}

/*
 * Transfers
 */

func (api *EtnaREST) GetTransfers(ctx context.Context, accId uint32) ([]sch.Transfer, error) {
	var resp sch.RespTransfers

	qry := url.Values{"pageNumber": {"0"}, "pageSize": {"99"}, "sortBy": {"TransferDate"}, "isDesc": {"false"}}
	err := (*api).callAPI(ctx, http.MethodGet, fmt.Sprintf("v1.0/accounts/%d/transfers", accId), qry, nil, &resp)
	if err != nil {
		return resp.Result, fmt.Errorf("getTransfers failed: %+v", err)
	}
	return resp.Result, nil
}

/*
 * Orders, trades
 */

// GetOrders retrieves a list of orders for a specific account.
// Supports filtering for active orders and returns them sorted by creation date.
func (api *EtnaREST) GetOrders(ctx context.Context, accId uint32, active bool) ([]sch.Order, error) {
	var resp sch.RespOrders

	qry := url.Values{"pageNumber": {"0"}, "pageSize": {"99"}, "sortField": {"CreateDate"}, "desc": {"false"}}
	if active {
		qry.Set("filter", "Status in (0,1,10)")
	}
	err := (*api).callAPI(ctx, http.MethodGet, fmt.Sprintf("v1.0/accounts/%d/orders", accId), qry, nil, &resp)
	if err != nil {
		return resp.Result, fmt.Errorf("getOrders failed: %+v", err)
	}
	return resp.Result, nil
}

// GetOrder retrieves details for a specific order within an account.
func (api *EtnaREST) GetOrder(ctx context.Context, accId uint32, orderId uint64) (sch.Order, error) {
	var resp sch.Order

	err := (*api).callAPI(ctx, http.MethodGet, fmt.Sprintf("v1.0/accounts/%d/orders/%d", accId, orderId), nil, nil, &resp)
	if err != nil {
		return resp, fmt.Errorf("getOrder failed: %+v", err)
	}
	return resp, nil
}

// PlaceOrder submits a new order for a specific account.
// It automatically sets default values for TimeInforce and ExtendedHours if not provided.
func (api *EtnaREST) PlaceOrder(ctx context.Context, accId uint32, params *sch.OrderParams) (sch.Order, error) {
	var resp sch.Order

	if params.TimeInforce == "" {
		params.TimeInforce = sch.TimeInForceGTC
	}
	if params.ExtendedHours == "" {
		params.ExtendedHours = sch.SessAll
	}

	err := (*api).callAPI(ctx, http.MethodPost, fmt.Sprintf("v1.0/accounts/%d/orders", accId), nil, params, &resp)
	if err != nil {
		return resp, fmt.Errorf("placeOrder failed: %+v", err)
	}
	return resp, nil
}

// ReplaceOrder modifies an existing order for a specific account.
func (api *EtnaREST) ReplaceOrder(ctx context.Context, accId uint32, orderId int64, params *sch.OrderParams) (sch.Order,
	error) {
	var resp sch.Order

	err := (*api).callAPI(ctx, http.MethodPut, fmt.Sprintf("v1.0/accounts/%d/orders/%d", accId, orderId), nil, params, &resp)
	if err != nil {
		return resp, fmt.Errorf("replaceOrder failed: %+v", err)
	}
	return resp, nil
}

// CancelOrder cancels an existing order for a specific account.
func (api *EtnaREST) CancelOrder(ctx context.Context, accId uint32, orderId int64) error {
	err := (*api).callAPI(ctx, http.MethodDelete, fmt.Sprintf("v1.0/accounts/%d/orders/%d", accId, orderId), nil, nil, nil)
	if err != nil {
		return fmt.Errorf("cancelOrder failed: %+v", err)
	}
	return nil
}

/*
 * Market data (securities, bars, etc)
 */

// GetSecurity retrieves details for a specific security by its symbol.
func (api *EtnaREST) GetSecurity(ctx context.Context, symbol string) (sch.Security, error) {
	var resp sch.Security
	err := (*api).callAPI(ctx, http.MethodGet, fmt.Sprintf("v1.0/equities/%s", symbol), nil, nil, &resp)
	if err != nil {
		return resp, fmt.Errorf("getSecurity failed: %+v", err)
	}
	if _, exist := NasdaqMICs[resp.Exchange]; exist {
		resp.Exchange = NASDAQ
	}
	return resp, nil
}

// GetBars retrieves historical bar data for a security based on the provided parameters.
// It validates the timeframe and makes a PUT request to the history API.
func (api *EtnaREST) GetBars(ctx context.Context, params *sch.ReqBars) ([]sch.Bar, error) {
	if _, exist := sch.VALID_TFS[params.Settings.Period]; !exist {
		return nil, fmt.Errorf("wrong timeframe: %s", params.Settings.Period)
	}
	var resp sch.RespBars
	params.Settings.Interval = 0   // bars
	params.Indicators = []string{} // no indicators
	if err := (*api).callAPI(ctx, http.MethodPut, "v1.0/history/symbols", nil, params, &resp); err != nil {
		return nil, fmt.Errorf("getBars failed: %+v", err)
	} else if len(resp.Bars) == 0 {
		return nil, fmt.Errorf("bars are absent")
	}
	return resp.Bars[0], nil
}
