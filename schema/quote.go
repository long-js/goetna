package schema

import (
	"bytes"
	"strconv"
	"time"
)

type Quote struct {
	Time          QuoteTime
	Ask           float64
	Bid           float64
	Last          float64
	Size          float64
	AskSize       float64
	BidSize       float64
	PreviousClose float64
	SymbolId      string
	Type          string
}

func (q *Quote) Parse(values map[string]string) error {
	var err error

	for k, v := range values {
		switch k {
		case "Date":
			qt := QuoteTime{}
			if err = qt.UnmarshalJSON([]byte(v)); err != nil {
				return err
			}
			(*q).Time = qt
		case "Ask":
			(*q).Ask, err = strconv.ParseFloat(v, 64)
		case "Bid":
			(*q).Bid, err = strconv.ParseFloat(v, 64)
		case "Price":
			(*q).Last, err = strconv.ParseFloat(v, 64)
		case "Volume":
			(*q).Size, err = strconv.ParseFloat(v, 64)
		case "AskSize":
			(*q).AskSize, err = strconv.ParseFloat(v, 64)
		case "BidSize":
			(*q).BidSize, err = strconv.ParseFloat(v, 64)
		case "PreviousClose":
			(*q).PreviousClose, err = strconv.ParseFloat(v, 64)
		case "Key":
			(*q).SymbolId = v
		case "QuoteTypes":
			(*q).Type = v
		}
		if err != nil {
			return err
		}
	}
	return nil
}

const QuoteTimeLayout = "01/02/2006 15:04:05"

type QuoteTime time.Time

func (qt *QuoteTime) UnmarshalJSON(b []byte) error {
	b = bytes.Trim(b, `"`)
	if t, err := time.Parse(QuoteTimeLayout, string(b)); err != nil {
		*qt = QuoteTime{}
		return err
	} else {
		*qt = QuoteTime(t)
	}
	return nil
}

func (qt QuoteTime) String() string {
	return time.Time(qt).Format(time.DateTime)
}
