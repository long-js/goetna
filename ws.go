package goetna

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	gjson "github.com/goccy/go-json"
	gws "github.com/gorilla/websocket"
	sch "github.com/long-js/goetna/schema"
)

const (
	ConnectTimeout    = 20 // the amount of seconds to wait for connection and authentication
	WSReconnInterval  = 12 // base period in seconds for the reconnect period calculation
	WSMaxSilentPeriod = 30 // maximum period of silence, seconds
)

type ConnHandler func(name string)
type DisconnHandler func(code int, text string) error
type MessageHandler func(topic string, dec *gjson.Decoder) error

func NewWSClient(name string, logger Logger, hdlConn ConnHandler, hdlDisconn DisconnHandler) WSClient {
	ctx, ctxCancel := context.WithCancel(context.Background())
	return WSClient{
		name:          name,
		logger:        logger,
		ctx:           ctx,
		ctxCancel:     ctxCancel,
		mu:            sync.Mutex{},
		wg:            sync.WaitGroup{},
		hdlConnect:    hdlConn,
		hdlDisconnect: hdlDisconn,
		reqChan:       make(chan []byte, 100),
	}
}

type WSClient struct {
	name                string
	logger              Logger
	ctx                 context.Context
	ctxCancel           func()
	mu                  sync.Mutex
	wg                  sync.WaitGroup
	conn                *gws.Conn
	connected, loggedIn atomic.Bool
	connectFn           func() error
	topicGetterFn       func(data []byte) (string, error)
	hdlConnect          ConnHandler
	hdlDisconnect       DisconnHandler
	hdlMessage          MessageHandler
	lastMsgTs           atomic.Int64
	reqChan             chan []byte
}

// IsOperational returns the current connection status of the WebSocket client.
//
// Returns:
//
//	true if the client is currently connected, false otherwise.
func (ws *WSClient) IsOperational() bool {
	return (*ws).connected.Load() && (*ws).loggedIn.Load()
}

func (ws *WSClient) SetConnectFunc(f func() error) {
	(*ws).connectFn = f
}

func (ws *WSClient) SetTopicFunc(f func(data []byte) (string, error)) {
	(*ws).topicGetterFn = f
}

func (ws *WSClient) SetMessageHandler(h MessageHandler) {
	(*ws).hdlMessage = h
}

// Start initiates the WebSocket connection and starts background processes for receiving messages.
func (ws *WSClient) Start() error {
	if (*ws).connectFn == nil {
		return fmt.Errorf("connect function is absent")
	} else if err := (*ws).connectFn(); err != nil {
		return err
	}

	go (*ws).goReceiver()
	go (*ws).goSender()
	// go (*ws).goPing()

	var i int
	for ; !(*ws).IsOperational() && i < ConnectTimeout; i++ {
		time.Sleep(1 * time.Second)
	}
	if i == ConnectTimeout && !(*ws).IsOperational() {
		return fmt.Errorf("connection timeout")
	}
	return nil
}

func (ws *WSClient) Stop() {
	(*ws).ctxCancel()
}

func (ws *WSClient) reconnect() {
	var err error

	for i := float64(0); i < 3<<8; i++ {
		period := math.Abs(WSReconnInterval*10*math.Sin(i/(2*WSReconnInterval))+i) + WSReconnInterval
		time.Sleep(time.Duration(period) * time.Second)
		(*ws).conn = nil
		if err = (*ws).Start(); err == nil {
			// TODO check for automatic resubsciption after session's been restored.
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
		(*ws).logger.Error("reconnection fault %+v", err)
	}
	if err != nil {
		(*ws).logger.Error("giving up with reconnection: %+v", err)
	}
}

func (ws *WSClient) disconnect() error {
	(*ws).connected.Store(false)
	return nil
}

// onPong handles incoming WebSocket pong messages.
func (ws *WSClient) onPong(data string) error {
	(*ws).logger.Debug("<-- PONG %s", data)
	(*ws).lastMsgTs.Store(time.Now().Unix())
	return nil
}

// goReceiver is a goroutine that continuously reads WebSocket messages, extracts the topic,
// and dispatches them to the onMessage handler. It also handles disconnections, potential panics,
// and initiates reconnection attempts.
func (ws *WSClient) goReceiver() {
	defer (*ws).logger.Info("receiver finished: %s", (*ws).name)
	defer func() {
		if errMsg := recover(); errMsg != nil {
			(*ws).logger.Error("receiver got panic: %+v\n%s", errMsg, debug.Stack())
			go (*ws).goReceiver()
		} else {
			select {
			case <-(*ws).ctx.Done():
			default:
				if err := (*ws).disconnect(); err != nil {
					(*ws).logger.Error("receiver disconnect fault: %s %+v", (*ws).name, errMsg)
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

		if (*ws).topicGetterFn != nil {
			if topic, err = (*ws).topicGetterFn(sockBuf); err != nil {
				(*ws).logger.Error("can't get topic: %+v", err)
				continue
			}
		}

		// if topic != sch.WSTopicCandle && topic != sch.WSTopicQuote && topic != sch.WSCmdPing {
		if topic != sch.WSCmdPing {
			(*ws).logger.Debug("<-- %s", sockBuf)
		}
		if _, err = buffer.Write(sockBuf); err != nil {
			(*ws).logger.Error("can't write to buffer %+v", err)
			continue
		} else if err = (*ws).hdlMessage(topic, dec); err != nil {
			(*ws).logger.Error("message processing fault: %s, %+v", topic, err)
		}
		buffer.Reset()
	}
}

// goSender is a goroutine that continuously reads WebSocket messages, extracts the topic,
func (ws *WSClient) goSender() {
	defer (*ws).logger.Info("sender finished")
	(*ws).wg.Add(1)
	defer (*ws).wg.Done()

	var (
		err     error
		req     []byte
		done    = (*ws).ctx.Done()
		cmdPong = "{\"Cmd\":\"Pong"
	)

loop:
	for connected := (*ws).connected.Load(); connected; connected = (*ws).connected.Load() {
		select {
		case <-done:
			break loop
		case req = <-(*ws).reqChan:
			if err = (*(*ws).conn).WriteMessage(gws.TextMessage, req); err != nil {
				(*ws).logger.Error("can't send message %+v, %+v", req, err)
				continue
			}
			if string(req[:12]) != cmdPong {
				(*ws).logger.Debug("--> %s", req)
			}
		}
	}
}
