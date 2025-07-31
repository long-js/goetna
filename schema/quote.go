package schema

type Quote struct {
	TsNano int64   `json:"t"`
	Ap     float64 `json:"ap"`
	As     float64 `json:"as"`
	Bp     float64 `json:"bp"`
	Bs     float64 `json:"bs"`
	Lp     float64 `json:"lp"`
	Ls     float64 `json:"ls"`
	Symbol string  `json:"Key"`
	Type   string  `json:"type"`
}

// const QuoteTimeLayout = "01/02/2006 15:04:05"
//
// type QuoteTime time.Time
//
// func (qt *QuoteTime) UnmarshalJSON(b []byte) error {
// 	b = bytes.Trim(b, `"`)
// 	if t, err := time.Parse(QuoteTimeLayout, string(b)); err != nil {
// 		*qt = QuoteTime{}
// 		return err
// 	} else {
// 		*qt = QuoteTime(t)
// 	}
// 	return nil
// }
