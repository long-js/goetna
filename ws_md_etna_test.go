//go:build test

package goetna

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/long-js/goetna/schema"
)

func createEtnaWS(private bool) *EtnaWS {
	var (
		sessType     schema.WSSessionType
		stream       schema.Streamer
		streamSessId schema.SessionId
	)
	l, p := loadCreds()

	if resp, err := (*rest).GetStreamers(ctx, false); err != nil {
		panic(err)
	} else if private {
		stream = resp.DataAddresses[0]
	} else {
		stream = resp.QuoteAddresses[0]
	}
	stream.Url, _ = strings.CutSuffix(stream.Url, ":443")

	if private {
		sessType = schema.WSSessData
	} else {
		sessType = schema.WSSessQuote
	}
	streamSessId, _ = (*rest).RecoverStreamerSession(ctx, sessType)
	return NewEtnaWS("TestEtnaWS", stream.Url, l, p, stream.SessionId, streamSessId,
		ColouredLogger("WSData"), onConnect, onDisconnect)
}

func onConnect(name string) {
	fmt.Printf("Connected callback %s\n", name)
}

func onDisconnect(code int, text string) error {
	fmt.Printf("Disconnected callback: %d %s\n", code, text)
	return nil
}

func TestEtnaWsStart(t *testing.T) {
	(*t).Skip()
	ws := createEtnaWS(false)

	if err := ws.Start(); err != nil {
		(*t).Error(err)
	} else {
		time.Sleep(100 * time.Second)
	}
}

func TestEtnaWsQuotes(t *testing.T) {
	// (*t).Skip()
	ws := createEtnaWS(false)

	if err := (*ws).Start(); err != nil {
		(*t).Error(err)
	} else if err = (*ws).Subscribe(schema.WSTopicQuote, "230226"); err != nil { // 3803 AAPL demo, 230226 AAPL prod
		(*t).Error(err)
		// } else if err = (*ws).Subscribe(schema.WSTopicCandle, "AAPL|NGS|USD:1m"); err != nil { // 3803 AAPL demo, 230226 AAPL prod
		// 	(*t).Error(err)
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
