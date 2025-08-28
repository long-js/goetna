//go:build test

package goetna

import (
	"testing"

	"github.com/long-js/goetna/schema"
)

func TestWsSubscription(t *testing.T) {
	// (*t).Skip()
	ws := createEtnaWS(true)

	if err := (*ws).Start(); err != nil {
		(*t).Error(err)
	}
	accId := "421" // 292
	if err := (*ws).Subscribe(schema.WSTopicOrder, accId); err != nil {
		(*t).Error(err)
	} else if err = (*ws).Subscribe(schema.WSTopicBalance, accId); err != nil {
		(*t).Error(err)
	} else if err = (*ws).Subscribe(schema.WSTopicPosition, accId); err != nil {
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
