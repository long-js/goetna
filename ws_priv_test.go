//go:build test

package goetna

import (
	"testing"

	"github.com/long-js/goetna/schema"
)

func TestWsSubscription(t *testing.T) {
	// (*t).Skip()
	ws := createWS(true)

	if err := (*ws).Start(); err != nil {
		(*t).Error(err)
	}
	if err := (*ws).Subscribe(schema.WSTopicOrder, "292"); err != nil {
		(*t).Error(err)
	} else if err = (*ws).Subscribe(schema.WSTopicBalance, "292"); err != nil {
		(*t).Error(err)
	} else if err = (*ws).Subscribe(schema.WSTopicPosition, "292"); err != nil {
		(*t).Error(err)
	}

	for cnt := 0; cnt < 20; cnt++ {
		select {
		case d := <-(*ws).BalanceChan:
			(*t).Logf("BAL: %+v\n", d)
		case d := <-(*ws).PositionsChan:
			(*t).Logf("POS: %+v\n", d)
		case d := <-(*ws).OrdersChan:
			(*t).Logf("ORDER: %+v\n", d)
		}

	}
}
