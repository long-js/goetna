package goetna

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	gjson "github.com/goccy/go-json"
	gws "github.com/gorilla/websocket"
	sch "github.com/long-js/goetna/schema"
)

type ConnState uint32

// NewEtnaWS creates the instance of EtnaWS, connects it to th e `baseUrl` and starts receive and ping goroutines.
func NewEtnaWS(name, url string, login, passwd []byte, userSessId, streamSessId sch.SessionId,
	logger Logger, hdlConn ConnHandler, hdlDisconn DisconnHandler) *EtnaWS {
	ws := EtnaWS{
		WSClient:      NewWSClient(name, logger, hdlConn, hdlDisconn),
		url:           url,
		login:         login,
		passwd:        passwd,
		userSessId:    userSessId,
		streamSessId:  streamSessId,
		subsciptions:  map[string][]sch.EtnaSubReq{},
		QuotesChan:    make(chan sch.EtnaQuote, 1000),
		BarsChan:      make(chan sch.Bar, 100),
		BalanceChan:   make(chan sch.TradingBalance, 20),
		PositionsChan: make(chan sch.Position, 20),
		OrdersChan:    make(chan sch.Order, 100),
	}
	ws.SetConnectFunc(ws.connect)
	ws.SetTopicFunc(getEtnaTopic)
	ws.SetMessageHandler(ws.onMessage)
	return &ws
}

type EtnaWS struct {
	WSClient
	url                      string
	streamSessId, userSessId sch.SessionId
	userId                   int32
	login, passwd            []byte
	subsciptions             map[string][]sch.EtnaSubReq

	QuotesChan    chan sch.EtnaQuote
	BarsChan      chan sch.Bar
	BalanceChan   chan sch.TradingBalance
	PositionsChan chan sch.Position
	OrdersChan    chan sch.Order
}

// createUrl generates the WebSocket connection URL with the necessary authentication parameters.
// It prioritizes using existing session credentials if available, otherwise it decodes and uses login and password.
func (ws *EtnaWS) createUrl() (string, error) {
	var (
		err  error
		v    url.Values
		decL = make([]byte, base64.StdEncoding.DecodedLen(len((*ws).login)))
		decP = make([]byte, base64.StdEncoding.DecodedLen(len((*ws).passwd)))
	)

	if (*ws).streamSessId != "" && (*ws).userSessId != "" && (*ws).userId != 0 {
		v = url.Values{
			"User":     {fmt.Sprintf("%d:%s", (*ws).userId, (*ws).userSessId)},
			"Password": {string((*ws).streamSessId)}, "HttpClientType": {"WebSocket"}}
	} else if _, err = base64.StdEncoding.Decode(decL, (*ws).login); err != nil {
		return "", fmt.Errorf("can't decode login %v", err)
	} else if _, err = base64.StdEncoding.Decode(decP, (*ws).passwd); err != nil {
		return "", fmt.Errorf("can't decode password %v", err)
	} else {
		v = url.Values{
			"User": {string(bytes.Trim(decL, "\x00"))}, "Password": {string(bytes.Trim(decP, "\x00"))},
			"HttpClientType": {"WebSocket"}}
	}
	return fmt.Sprintf("%s/CreateSession.txt?%s", (*ws).url, v.Encode()), nil
}

// connect establishes a new WebSocket connection to the Etna API.
// It handles URL creation, sets headers, dials the server, configures handlers for pong and disconnect events,
// and updates the connection status.
func (ws *EtnaWS) connect() error {
	var (
		err error
		uri string
	)
	if (*ws).conn != nil {
		return fmt.Errorf("connection already exists")
	} else if uri, err = (*ws).createUrl(); err != nil {
		return err
	}
	header := make(http.Header)
	header["User-Agent"] = []string{"qant/2.0"}
	header["Accept-Encoding"] = []string{"gzip, deflate"}

	dialer := gws.Dialer{EnableCompression: true, HandshakeTimeout: 45 * time.Second}
	if isTest, err := strconv.ParseBool(os.Getenv("TEST_ENV")); err == nil && isTest {
		dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	conn, response, err := dialer.DialContext((*ws).ctx, uri, header)
	if err != nil {
		return fmt.Errorf("failed to connect %s: status: %s, %+v", uri, (*response).Status, err)
	}
	conn.SetPongHandler((*ws).onPong)
	if (*ws).hdlDisconnect != nil {
		conn.SetCloseHandler((*ws).hdlDisconnect)
	}
	(*ws).conn = conn
	(*ws).connected.Store(true)
	(*ws).logger.Info("connected: %s [%s] %s, close: %t", (*ws).url, response.Header.Get("Server"),
		response.Header.Get("Date"), response.Close)
	return nil
}

// sendJson marshals a sch.EtnaSubReq struct into JSON and sends it as a binary WebSocket message.
func (ws *EtnaWS) sendJson(message *sch.EtnaSubReq) error {
	var (
		err error
		buf []byte
	)
	if buf, err = gjson.Marshal(*message); err != nil {
		return fmt.Errorf("can't marshal %+v, %+v", *message, err)
	}
	(*ws).reqChan <- buf
	return nil
}

// Subscribe sends a subscription request for a specific topic and keys.
// It checks for an existing connection and prevents duplicate subscriptions.
func (ws *EtnaWS) Subscribe(topic string, keys string) error {
	if (*ws).conn == nil {
		return fmt.Errorf("not connected")
	}
	var (
		exist bool
		subs  []sch.EtnaSubReq
	)
	if subs, exist = (*ws).subsciptions[topic]; !exist {
		subs = make([]sch.EtnaSubReq, 0, 1)
		(*ws).subsciptions[topic] = subs
	}

	exist = false
	for i := 0; i < len(subs); i++ {
		if subs[i].Keys == keys {
			exist = true
			break
		}
	}
	if !exist {
		sub := sch.EtnaSubReq{
			Cmd: "Subscribe.txt", SessionId: (*ws).userSessId, Keys: keys, Topic: topic, HttpClientType: "WebSocket"}
		if err := (*ws).sendJson(&sub); err != nil {
			return err
		}
	}
	return nil
}

// Unsubscribe sends an unsubscription request for a specific topic and keys.
// It checks for an existing connection and the presence of the subscription before sending the unsubscribe command.
func (ws *EtnaWS) Unsubscribe(topic string, keys string) error {
	if (*ws).conn == nil {
		return fmt.Errorf("not connected")
	}
	var (
		exist bool
		subs  []sch.EtnaSubReq
	)
	if subs, exist = (*ws).subsciptions[topic]; !exist {
		return fmt.Errorf("subscription type is absent: %s, %s", topic, keys)
	}
	for i := 0; i < len(subs); i++ {
		if subs[i].Keys == keys {
			sub := subs[i]
			sub.Cmd = "Unsubscribe.txt"
			if err := (*ws).sendJson(&sub); err != nil {
				return err
			}
			subs = append(subs[:i], subs[i+1:]...)
			break
		}
	}
	return nil
}

// onMessage processes incoming WebSocket messages based on the topic.
// It decodes the JSON payload into the corresponding struct (EtnaQuote, Order, Balance, or Position)
// and sends it to the appropriate channel.
func (ws *EtnaWS) onMessage(topic string, dec *gjson.Decoder) error {
	var (
		err      error
		quote    sch.EtnaQuote
		bar      sch.Bar
		balance  sch.TradingBalance
		position sch.Position
		order    sch.Order
		sub      sch.EtnaSubReq
	)
	switch topic {
	case sch.WSTopicQuote:
		sBuf := map[string]string{}
		if err = dec.Decode(&sBuf); err != nil {
			return fmt.Errorf("quote decoding fault %+v", err)
		} else if err = quote.Parse(sBuf); err != nil {
			return fmt.Errorf("quote decoding fault %+v", err)
		}
		(*ws).QuotesChan <- quote
	case sch.WSTopicCandle:
		sBuf := map[string]string{}
		if err = dec.Decode(&sBuf); err != nil {
			return fmt.Errorf("bar decoding fault %+v", err)
		} else if err = bar.Parse(sBuf); err != nil {
			return fmt.Errorf("bar decoding fault %+v", err)
		} else if bar.IsCompleted {
			(*ws).BarsChan <- bar
		}
	case sch.WSTopicOrder:
		sBuf := map[string]string{}
		if err = dec.Decode(&sBuf); err != nil {
			return fmt.Errorf("order decoding fault %+v", err)
		} else if err = order.Parse(sBuf); err != nil {
			return fmt.Errorf("order decoding fault %+v", err)
		}
		(*ws).OrdersChan <- order
	case sch.WSTopicBalance:
		if err = dec.Decode(&balance); err != nil {
			return fmt.Errorf("balance decoding fault %+v", err)
		} else if err = balance.Parse(); err != nil {
			return fmt.Errorf("balance parsing fault %+v", err)
		}
		(*ws).BalanceChan <- balance
	case sch.WSTopicPosition:
		sBuf := map[string]string{}
		if err = dec.Decode(&sBuf); err != nil {
			return fmt.Errorf("position decoding fault %+v", err)
		} else if err = position.Parse(sBuf); err != nil {
			return fmt.Errorf("position decoding fault %+v", err)
		}
		(*ws).PositionsChan <- position
	case sch.WSCmdPing:
		(*ws).reqChan <- sch.WSPongMsg
	case sch.WSCmdSub:
		if err = dec.Decode(&sub); err != nil {
			return fmt.Errorf("subscription decoding fault %+v", err)
		}
		if subs, exist := (*ws).subsciptions[sub.Topic]; exist {
			subs = append(subs, sub)
			(*ws).logger.Info("Subscribed %s: %s [%s]", sub.Topic, sub.Keys, sub.SessionId)
		}
	case sch.WSCmdUnsub:
		if err = dec.Decode(&sub); err != nil {
			return fmt.Errorf("unsubscription decoding fault %+v", err)
		}
		if _, exist := (*ws).subsciptions[sub.Topic]; exist {
			(*ws).logger.Info("TODO Unsubscribed %s: %s [%s]", sub.Topic, sub.Keys, sub.SessionId)
		}
	case sch.WSCmdCreate:
		msg := map[string]string{}
		if err = dec.Decode(&msg); err != nil {
			return fmt.Errorf("CreateSession decoding fault %+v", err)
		}
		(*ws).userSessId = sch.SessionId(msg["SessionId"])
		(*ws).loggedIn.Store(true)
		(*ws).hdlConnect((*ws).name)
		(*ws).logger.Info("Websocket session created: %s", msg["SessionId"])
	default:
		return fmt.Errorf("wrong message %s", topic)
	}
	return nil
}

// getTopic returns the message topic.
func getEtnaTopic(data []byte) (string, error) {
	var (
		err   error
		topic string
	)
	if topic, err = etnaMsgType(data, false); err != nil {
		if topic, err = etnaMsgType(data, true); err != nil {
			return "", fmt.Errorf("WS message has neither topic not cmd: %+v, %s", err, data)
		}
	}
	return topic, nil
}

// goPing is a goroutine that periodically sends WebSocket ping messages to the server
// to maintain the connection. It checks the last received message timestamp and initiates
// a disconnect and reconnect if the silent period exceeds a defined threshold.
func (ws *EtnaWS) goPing() {
	defer (*ws).logger.Info("pinger finished")

	done := (*ws).ctx.Done()
	ticker := time.NewTicker(time.Second).C
	for now := range ticker {
		if !(*ws).connected.Load() {
			break
		}
		select {
		case <-done:
			return
		default:
			lastTs := (*ws).lastMsgTs.Load()
			if now.Unix()-lastTs > WSMaxSilentPeriod {
				(*ws).mu.Lock()
				if (*ws).conn == nil {
					(*ws).mu.Unlock()
					continue
				}
				err := (*(*ws).conn).WriteControl(gws.PingMessage, []byte(strconv.FormatInt(now.Unix(), 10)),
					now.Add(500*time.Millisecond))
				(*ws).mu.Unlock()
				if err != nil {
					(*ws).connected.Store(false)
					(*ws).logger.Error("ping fault: %v", err)
				}

			}
		}
	}
}

// etnaMsgType extracts the message type from a raw WebSocket message byte slice.
func etnaMsgType(data []byte, isCmd bool) (string, error) {
	end := 15
	searchField := sch.FieldEntytyType
	if isCmd {
		end = 8
		searchField = sch.FieldCmd
	}
	if data[0] != 123 || string(data[1:end]) != searchField {
		return "", fmt.Errorf("wrong data")
	}

	res := ""
	for i := end + 1; i < len(data)-2; i++ {
		if data[i] == 34 {
			res = string(data[end+1 : i])
			break
		}
	}
	if res == "" {
		return res, fmt.Errorf("topic not found")
	}
	return res, nil
}
