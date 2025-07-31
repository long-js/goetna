//go:build test

package goetna

import (
	"fmt"
	"testing"
	"time"
)

var wsQ, wsD = createWSs()

func createWSs() (*EtnaWS, *EtnaWS) {
	if resp, err := (*rest).GetStreamers(ctx); err != nil {
		panic(err)
	} else {
		da := resp.DataAddresses[0]
		qa := resp.QuoteAddresses[0]
		_wsD := NewEtnaWS(da.Url, "", "", da.SessionId, ColouredLogger("WSData"), onConnect, onDisconnect)
		_wsQ := NewEtnaWS(qa.Url, "", "", qa.SessionId, ColouredLogger("WSQuote"), onConnect, onDisconnect)
		return _wsQ, _wsD
	}
	return nil, nil
}

func onConnect() {
	fmt.Println("Connected")
}

func onDisconnect(code int, text string) error {
	fmt.Printf("Disconnected: %d %s\n", code, text)
	return nil
}

func TestStart(t *testing.T) {
	// (*t).Skip()
	var err error

	if err = wsD.Start(); err != nil {
		(*t).Error(err)
	} else {
		time.Sleep(5 * time.Second)
	}
}
