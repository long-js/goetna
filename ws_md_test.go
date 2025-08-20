//go:build test

package goetna

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/long-js/goetna/schema"
)

func createWS(private, isNvb bool) *EtnaWS {
	var (
		err      error
		sessType schema.WSSessionType
		stream   schema.Streamer
	)
	l, p := loadCreds()
	if private {
		sessType = schema.WSSessData
		stream.Url = DefaultConfig.WSUrlPriv
		stream.SessionId, err = (*rest).RecoverStreamerSession(ctx, sessType)
	} else {
		sessType = schema.WSSessQuote
		if isNvb {
			stream.Url = DefaultConfig.WSUrlPubNvb
		} else {
			stream.Url = DefaultConfig.WSUrlPub
			stream.SessionId, err = (*rest).RecoverStreamerSession(ctx, sessType)
		}
	}
	if err != nil || stream.SessionId == "" {
		if resp, err := (*rest).GetStreamers(ctx, true); err != nil {
			panic(err)
		} else if private {
			stream = resp.DataAddresses[0]
		} else {
			stream = resp.QuoteAddresses[0]
		}
		stream.Url, _ = strings.CutSuffix(stream.Url, ":443")
	}
	return NewEtnaWS("TestWS", stream.Url, l, p, stream.SessionId, ColouredLogger("WSData"), onConnect, onDisconnect)
}

func onConnect(name string) {
	fmt.Printf("Connected callback %s", name)
}

func onDisconnect(code int, text string) error {
	fmt.Printf("Disconnected callback: %d %s\n", code, text)
	return nil
}

func TestWsStart(t *testing.T) {
	(*t).Skip()
	ws := createWS(false, false)

	if err := ws.Start(); err != nil {
		(*t).Error(err)
	} else {
		time.Sleep(100 * time.Second)
	}
}

func TestWsQuotes(t *testing.T) {
	// (*t).Skip()
	ws := createWS(false, false)

	if err := (*ws).Start(); err != nil {
		(*t).Error(err)
	} else if err = (*ws).Subscribe("Candle", "AAPL|NGS|USD:1m;AAPL|NGS|USD:15m"); err != nil {
		// } else if err = (*ws).Subscribe(schema.WSTopicQuote, "3803"); err != nil {
		(*t).Error(err)
	}

	for cnt := 0; cnt < 100; cnt++ {
		select {
		case q := <-(*ws).QuotesChan:
			(*t).Logf("QUOT: %s %f %f\n", q.Time, q.Last, q.Size)
		case bar := <-(*ws).BarsChan:
			(*t).Logf("BAR: %+v\n", bar)
		}
	}
}
