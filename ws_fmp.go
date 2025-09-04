package goetna

import (
	"fmt"
	"net/http"
	"time"

	gjson "github.com/goccy/go-json"
	gws "github.com/gorilla/websocket"
	sch "github.com/long-js/goetna/schema"
)

// NewFmpWS creates the instance of FmpWS
func NewFmpWS(name string, fmpKey string, logger Logger, hdlConn ConnHandler, hdlDisconn DisconnHandler) *FmpWS {
	ws := FmpWS{
		WSClient:     NewWSClient(name, logger, hdlConn, hdlDisconn),
		fmpKey:       fmpKey,
		subsciptions: make(map[string]struct{}),
		QuotesChan:   make(chan sch.FmpQuote, 1000),
	}
	ws.SetConnectFunc(ws.connect)
	ws.SetTopicFunc(getFmpEvent)
	ws.SetMessageHandler(ws.onMessage)
	return &ws
}

type FmpWS struct {
	WSClient
	fmpKey       string
	subsciptions map[string]struct{}
	QuotesChan   chan sch.FmpQuote
}

// connect establishes a new WebSocket connection to the FMP API.
// It sets headers, dials the server and updates the connection status.
func (ws *FmpWS) connect() error {
	if (*ws).conn != nil {
		return fmt.Errorf("connection already exists")
	}
	header := make(http.Header)
	header["User-Agent"] = []string{"qant/2.0"}
	header["Accept-Encoding"] = []string{"gzip, deflate"}

	dialer := gws.Dialer{EnableCompression: true, HandshakeTimeout: 45 * time.Second}
	if conn, response, err := dialer.DialContext((*ws).ctx, DefaultConfig.WSUrlPubFMP, header); err != nil {
		return fmt.Errorf("failed to connect, status: %s, %+v", (*response).Status, err)
	} else {
		conn.SetPongHandler((*ws).onPong)
		if (*ws).hdlDisconnect != nil {
			conn.SetCloseHandler((*ws).hdlDisconnect)
		}
		(*ws).mu.Lock()
		(*ws).conn = conn
		(*ws).mu.Unlock()
		(*ws).connected.Store(true)
		(*ws).logger.Info("connected: %s [%s], close: %t", response.Header.Get("Server"),
			response.Header.Get("Date"), response.Close)
	}

	if err := (*ws).sendJson(&sch.FmpReq{Event: "login", Data: map[string]string{"apiKey": (*ws).fmpKey}}); err != nil {
		return err
	}
	return nil
}

// sendJson marshals a sch.Subscription struct into JSON and sends it as a binary WebSocket message.
func (ws *FmpWS) sendJson(message *sch.FmpReq) error {
	if buf, err := gjson.Marshal(*message); err != nil {
		return fmt.Errorf("can't marshal %+v, %+v", *message, err)
	} else {
		(*ws).reqChan <- buf
	}
	return nil
}

// Subscribe sends a subscription request for a specific topic and keys.
// It checks for an existing connection and prevents duplicate subscriptions.
func (ws *FmpWS) Subscribe(key string) error {
	(*ws).mu.Lock()
	if (*ws).conn == nil {
		(*ws).mu.Unlock()
		return fmt.Errorf("not connected")
	}
	(*ws).mu.Unlock()
	if _, exist := (*ws).subsciptions[key]; !exist {
		m := sch.FmpReq{Event: "subscribe", Data: map[string]string{"ticker": key}}
		if err := (*ws).sendJson(&m); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("already subscribed %s", key)
	}
	return nil
}

// Unsubscribe sends an unsubscription request for a specific topic and keys.
// It checks for an existing connection and the presence of the subscription before sending the unsubscribe command.
func (ws *FmpWS) Unsubscribe(key string) error {
	if (*ws).conn == nil {
		return fmt.Errorf("not connected")
	}
	if _, exist := (*ws).subsciptions[key]; !exist {
		return fmt.Errorf("subscription is absent: %s", key)
	} else {
		m := sch.FmpReq{Event: "unsubscribe", Data: map[string]string{"ticker": key}}
		if err := (*ws).sendJson(&m); err != nil {
			return err
		}
		delete((*ws).subsciptions, key)
	}
	return nil
}

// onMessage processes incoming WebSocket messages based on the topic.
// It decodes the JSON payload into the corresponding struct and sends it to the appropriate channel.
func (ws *FmpWS) onMessage(topic string, dec *gjson.Decoder) error {
	var (
		quote sch.FmpQuote
		resp  sch.FmpResponse
	)

	switch topic {
	case sch.WSTopicQuote:
		if err := dec.Decode(&quote); err != nil {
			return fmt.Errorf("decoding fault, %+v", err)
		}
		if quote.Last != 0. && quote.Type == "T" {
			(*ws).QuotesChan <- quote
		}
	case sch.WSTopicEvent:
		if err := dec.Decode(&resp); err != nil {
			return fmt.Errorf("FMP event decoding fault %+v", err)
		} else if resp.Event != sch.WSEvtHB && resp.Status != 200 {
			(*ws).logger.Error("FMP: %d %s", resp.Status, resp.Message)
			if resp.Event == sch.WSEvtLogin {
				(*ws).loggedIn.Store(false)
			}
			return nil
		}

		switch resp.Event {
		case sch.WSEvtHB:
			(*ws).logger.Debug("HB: %d", resp.Timestamp)
		case sch.WSEvtSub:
			if len(resp.Message) < 15 {
				return fmt.Errorf("FMP wrong response message %s", resp.Message)
			} else {
				key := resp.Message[14:]
				(*ws).subsciptions[key] = struct{}{}
				(*ws).logger.Info("Subscribed: %d %s", resp.Status, key)
			}
		case sch.WSEvtUnsub:
			if len(resp.Message) < 19 {
				return fmt.Errorf("FMP wrong response message %s", resp.Message)
			}

			key := resp.Message[18:]
			if _, exist := (*ws).subsciptions[key]; !exist {
				return fmt.Errorf("subscription doesn't exist: %s, %s", key, resp.Message)
			} else {
				delete((*ws).subsciptions, key)
				(*ws).logger.Info("Unsubscribed: %d %s", resp.Status, key)
			}
		case sch.WSEvtLogin:
			(*ws).loggedIn.Store(true)
			(*ws).logger.Info("Logged in: %d %s", resp.Status, resp.Message)
		}
	}
	return nil
}

func getFmpEvent(data []byte) (string, error) {
	if len(data) < 10 || data[0] != 123 {
		return "", fmt.Errorf("wrong FMP data: %s", data)
	}
	if string(data[:9]) == sch.FieldEvent {
		return sch.WSTopicEvent, nil
	}
	return sch.WSTopicQuote, nil
}
