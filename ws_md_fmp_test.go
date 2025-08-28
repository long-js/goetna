//go:build test

package goetna

import (
	"testing"
	"time"
)

func createFmpWS() *FmpWS {
	var fmpKey string

	if resp, err := (*rest).GetStreamers(ctx, true); err != nil {
		panic(err)
	} else {
		fmpKey = resp.FMPKey
	}

	return NewFmpWS("TestFmpWS", fmpKey, ColouredLogger("WSFmp"), onConnect, onDisconnect)
}

func TestFmpWsStart(t *testing.T) {
	(*t).Skip()
	ws := createFmpWS()

	if err := ws.Start(); err != nil {
		(*t).Error(err)
	} else {
		time.Sleep(100 * time.Second)
	}
}

func TestFmpWsQuotes(t *testing.T) {
	// (*t).Skip()
	ws := createFmpWS()

	if err := (*ws).Start(); err != nil {
		(*t).Error(err)
	} else if err = (*ws).Subscribe("aapl"); err != nil {
		(*t).Error(err)
	} else if err = (*ws).Subscribe("nvda"); err != nil {
		(*t).Error(err)
	}

	for cnt := 0; cnt < 100; cnt++ {
		q := <-(*ws).QuotesChan
		(*t).Logf("QUOT: %s %f %f\n", time.Unix(q.NTs/1e9, q.NTs%1e9), q.Last, q.Size)
	}
}

func TestEtnaWsReconnect(t *testing.T) {
	(*t).Skip()
	// TODO
}
