//go:build test

package goetna

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/khokhlomin/goetna/schema"
)

func createWS(private bool) *EtnaWS {
	var (
		err      error
		sessType schema.WSSessionType
		stream   schema.Streamer
	)
	l, p := loadCreds()
	if private {
		sessType = schema.WSSessData
		stream.Url = defaultConfig.WSUrlPriv
	} else {
		sessType = schema.WSSessQuote
		stream.Url = defaultConfig.WSUrlPub
	}
	if stream.SessionId, err = (*rest).RecoverStreamerSession(ctx, sessType); err != nil || stream.SessionId == "" {
		if resp, err := (*rest).GetStreamers(ctx); err != nil {
			panic(err)
		} else if private {
			stream = resp.DataAddresses[0]
		} else {
			stream = resp.QuoteAddresses[0]
		}
		stream.Url, _ = strings.CutSuffix(stream.Url, ":443")
	}
	return NewEtnaWS(stream.Url, l, p, stream.SessionId, ColouredLogger("WSData"), onConnect, onDisconnect)
}

func onConnect() {
	fmt.Println("Connected callback")
}

func onDisconnect(code int, text string) error {
	fmt.Printf("Disconnected callback: %d %s\n", code, text)
	return nil
}

func TestWsStart(t *testing.T) {
	(*t).Skip()
	ws := createWS(false)

	if err := ws.Start(); err != nil {
		(*t).Error(err)
	} else {
		time.Sleep(100 * time.Second)
	}
}

func TestWsQuotes(t *testing.T) {
	// (*t).Skip()
	ws := createWS(false)

	if err := (*ws).Start(); err != nil {
		(*t).Error(err)
		// } else if err = (*ws).Subscribe("Candle", "AAPL|NGS|USD:1m"); err != nil {
	} else if err = (*ws).Subscribe(schema.WSTopicQuote, "3803"); err != nil {
		(*t).Error(err)
	}

	for cnt := 0; cnt < 100; cnt++ {
		q := <-(*ws).QuotesChan
		(*t).Logf("QUOT: %s %f %f\n", q.Time, q.Last, q.Size)

	}
}
