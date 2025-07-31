package goetna

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	gjson "github.com/goccy/go-json"
	gws "github.com/gorilla/websocket"
	sch "github.com/khokhlomin/goetna/schema"
)

type ConnState uint32
type ConnHandler func()
type DisconnHandler func(code int, text string) error

// NewEtnaWS creates the instance of EtnaWS, connects it to th e `url` and starts receive and ping goroutines.
func NewEtnaWS(url, login, passwd string, sessionId sch.SessionId, logger Logger, hdlConn ConnHandler,
	hdlDisconn DisconnHandler) *EtnaWS {
	ws := EtnaWS{
		mu: sync.Mutex{}, url: url, login: login, passwd: passwd, sessionId: sessionId, logger: logger,
		hdlConnect: hdlConn, hdlDisconnect: hdlDisconn, subsciptions: map[string][]sch.Subscription{}}
	ws.ctx, ws.ctxCancel = context.WithCancel(context.Background())
	ws.reqChan = make(chan []byte)
	ws.msgChan = make(chan []byte)
	return &ws
}

type EtnaWS struct {
	url                      string
	sessionId, userSessionId sch.SessionId
	userId                   int32
	login, passwd            string

	logger           Logger
	ctx              context.Context
	ctxCancel        func()
	mu               sync.Mutex
	wg               sync.WaitGroup
	conn             *gws.Conn
	connected        atomic.Bool
	hdlConnect       ConnHandler
	hdlDisconnect    DisconnHandler
	lastMsgTs        atomic.Int64
	reqChan, msgChan chan []byte
	subsciptions     map[string][]sch.Subscription

	QuotesChan    chan sch.Quote
	BalanceChan   chan sch.TradingBalance
	PositionsChan chan sch.Position
	OrdersChan    chan sch.Order
}

// IsConnected returns the current connection status of the WebSocket client.
//
// Returns:
//
//	true if the client is currently connected, false otherwise.
func (ws *EtnaWS) IsConnected() bool {
	return (*ws).connected.Load()
}

// Start initiates the WebSocket connection and starts the background processes for receiving messages and sending pings.
func (ws *EtnaWS) Start() error {
	if err := (*ws).connect(); err != nil {
		return err
	}
	go (*ws).goReceive()
	go (*ws).goPing()
	return nil
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
	header := http.Header{"User-Agent": []string{"qant/2.0"}}
	dialer := gws.Dialer{EnableCompression: true, HandshakeTimeout: 45 * time.Second}
	conn, response, err := dialer.DialContext((*ws).ctx, uri, header)
	if err != nil {
		return fmt.Errorf("failed to connect: %+v, status: %s", err, (*response).Status)
	}
	conn.SetPongHandler((*ws).onPong)
	if (*ws).hdlDisconnect != nil {
		conn.SetCloseHandler((*ws).hdlDisconnect)
	}
	(*ws).conn = conn
	(*ws).connected.Store(true)
	(*ws).logger.Info("connected: %+v, close: %t %s", response.Header, response.Close, uri)
	return nil
}

func (ws *EtnaWS) disconnect() error {
	(*ws).connected.Store(false)
	return nil
}

func (ws *EtnaWS) reconnect() {
	var err error

	for i := float64(0); i < 3<<8; i++ {
		time.Sleep(time.Duration(math.Abs(sch.WSReconnInterval*10*math.Sin(i/(2*sch.WSReconnInterval))+i)+
			sch.WSReconnInterval) * time.Second)
		if err = (*ws).Start(); err == nil {
			// TODO check for automatic resubsciption after session's been resored.
			// resubscribe
			// for topic, subs := range (*ws).subsciptions {
			// 	for sIdx := 0; sIdx < len(subs); sIdx++ {
			// 		if err = (*ws).sendJson(&subs[sIdx]); err != nil {
			// 			(*ws).logger.Error("can't resubscribe: %s %+v, %+v", topic, subs[sIdx], err)
			// 		} else {
			// 			(*ws).logger.Debug("resubscribed %s %+v", topic, subs[sIdx])
			// 		}
			// 	}
			// }
			break
		}
		(*ws).logger.Error("reconnect fault %+v", err)
	}
	if err != nil {
		(*ws).logger.Error("giving up with reconnect: %+v", err)
	}
}

// createUrl generates the WebSocket connection URL with the necessary authentication parameters.
// It prioritizes using existing session credentials if available, otherwise it decodes and uses login and password.
func (ws *EtnaWS) createUrl() (string, error) {
	var v url.Values

	if (*ws).sessionId != "" && (*ws).userSessionId != "" && (*ws).userId != 0 {
		v = url.Values{
			"User":     {fmt.Sprintf("%d:%s", (*ws).userId, (*ws).userSessionId)},
			"Password": {string((*ws).sessionId)}, "HttpClientType": {"WebSocket"}}
	} else if login, err := base64.StdEncoding.DecodeString((*ws).login); err != nil {
		return "", fmt.Errorf("can't decode creds %v", err)
	} else if pwd, err := base64.StdEncoding.DecodeString((*ws).passwd); err != nil {
		return "", fmt.Errorf("can't decode creds %v", err)
	} else {
		v = url.Values{"User": {string(login)}, "Password": {string(pwd)}, "HttpClientType": {"WebSocket"}}
	}
	return fmt.Sprintf("%s/CreateSession.txt?%s", (*ws).url, v.Encode()), nil
}

// sendJson marshals a sch.Subscription struct into JSON and sends it as a binary WebSocket message.
func (ws *EtnaWS) sendJson(message *sch.Subscription) error {
	var (
		err error
		buf []byte
	)
	if buf, err = gjson.Marshal(*message); err != nil {
		return fmt.Errorf("can't marshal %+v, %+v", *message, err)
	} else if err = (*(*ws).conn).WriteMessage(gws.BinaryMessage, buf); err != nil {
		return fmt.Errorf("can't send message %+v, %+v", *message, err)
	}
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
		subs  []sch.Subscription
	)
	if subs, exist = (*ws).subsciptions[topic]; !exist {
		subs = make([]sch.Subscription, 0, 1)
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
		sub := sch.Subscription{
			Cmd: "Subscribe.txt", SessionId: (*ws).sessionId, Keys: keys, Topic: topic,
			HttpClientType: "WebSocket"}
		if err := (*ws).sendJson(&sub); err != nil {
			return err
		} else {
			subs = append(subs, sub)
		}
		(*ws).logger.Info("Subscribed %s: %s", topic, keys)
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
		subs  []sch.Subscription
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
// It decodes the JSON payload into the corresponding struct (Quote, Order, Balance, or Position)
// and sends it to the appropriate channel.
func (ws *EtnaWS) onMessage(topic string, dec *gjson.Decoder) error {
	var (
		err      error
		quote    sch.Quote
		balance  sch.TradingBalance
		position sch.Position
		order    sch.Order
	)
	switch topic {
	case sch.WSTopicQuote:
		if err = dec.Decode(&quote); err != nil {
			return fmt.Errorf("quote decoding fault %+v", err)
		}
		(*ws).QuotesChan <- quote
	case sch.WSTopicOrder:
		if err = dec.Decode(&order); err != nil {
			return fmt.Errorf("order decoding fault %+v", err)
		}
		(*ws).OrdersChan <- order
	case sch.WSTopicBalance:
		if err = dec.Decode(&balance); err != nil {
			return fmt.Errorf("balance decoding fault %+v", err)
		}
		(*ws).BalanceChan <- balance
	case sch.WSTopicPosition:
		if err = dec.Decode(&position); err != nil {
			return fmt.Errorf("position decoding fault %+v", err)
		}
		(*ws).PositionsChan <- position
	default:
		return fmt.Errorf("wrong message %s", topic)
	}
	return nil
}

// onPong handles incoming WebSocket pong messages.
func (ws *EtnaWS) onPong(data string) error {
	(*ws).logger.Debug("<-- PONG %s", data)
	(*ws).lastMsgTs.Store(time.Now().Unix())
	return nil
}

// goReceive is a goroutine that continuously reads WebSocket messages, extracts the topic,
// and dispatches them to the onMessage handler. It also handles disconnections, potential panics,
// and initiates reconnection attempts.
func (ws *EtnaWS) goReceive() {
	defer (*ws).logger.Info("receiver finished")
	defer func() {
		if errMsg := recover(); errMsg != nil {
			(*ws).logger.Error("receiver got panic: %+v\n%s", errMsg, debug.Stack())
			go (*ws).goReceive()
		} else {
			select {
			case <-(*ws).ctx.Done():
			default:
				if err := (*ws).disconnect(); err != nil {
					(*ws).logger.Error("receiver disconnect fault: %+v", errMsg)
				}
				(*ws).reconnect()
			}
		}
	}()
	(*ws).wg.Add(1)
	defer (*ws).wg.Done()

	var (
		err     error
		sockBuf []byte
		topic   string
		dec     *gjson.Decoder
		buf     = make([]byte, 0, 1024)
	)
	buffer := bytes.NewBuffer(buf)
	dec = gjson.NewDecoder(buffer)
	for connected := (*ws).connected.Load(); connected; connected = (*ws).connected.Load() {
		if _, sockBuf, err = (*(*ws).conn).ReadMessage(); err != nil {
			(*ws).lastMsgTs.Store(time.Now().Unix())
			(*ws).connected.Store(false)
			(*ws).logger.Error("reading message fault: %v", err)
			continue
		}
		(*ws).lastMsgTs.Store(time.Now().Unix())

		if topic, err = getTopic(sockBuf); err != nil {
			(*ws).logger.Error("WS topic: %+v, %s", err, sockBuf)
			continue
		} else if _, err = buffer.Write(sockBuf); err != nil {
			(*ws).logger.Error("can't write to buffer %+v", err)
			continue
		}
		if err = (*ws).onMessage(topic, dec); err != nil {
			(*ws).logger.Error("message processing fault: %s, %+v", topic, err)
		}
		buffer.Reset()
	}
}

// goPing is a goroutine that periodically sends WebSocket ping messages to the server
// to maintain the connection. It checks the last received message timestamp and initiates
// a disconnect and reconnect if the silent period exceeds a defined threshold.
func (ws *EtnaWS) goPing() {
	defer (*ws).logger.Info("pinger finished")

	ticker := time.NewTicker(time.Second).C
	for now := range ticker {
		if !(*ws).connected.Load() {
			break
		}
		select {
		case <-(*ws).ctx.Done():
			return
		default:
			lastTs := (*ws).lastMsgTs.Load()
			if now.Unix()-lastTs > sch.WSMaxSilentPeriod {
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

// getTopic extracts the topic string from a raw WebSocket message byte slice.
// It checks for a specific prefix and searches for the topic.
func getTopic(data []byte) (string, error) {
	if data[0] != 123 || !bytes.Equal(data[1:14], sch.EntytyType) {
		return "", fmt.Errorf("wrong data")
	}
	res := ""
	for i := 15; i < len(data)-2; i++ {
		if data[i] == 34 {
			res = string(data[15:i])
			break
		}
	}
	if res == "" {
		return res, fmt.Errorf("topic not found")
	}
	return res, nil
}
