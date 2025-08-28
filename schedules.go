package goetna

import "github.com/long-js/goetna/schema"

// BaseSchedules contains the time of morning, regular and evening trading sessions
var BaseSchedules = map[string]schema.MarketSchedule{
	"NGS": {
		MonOpen:  14400, // 04:00
		MonClose: 34199, // 09:29:59
		RegOpen:  34200, // 09:30
		RegClose: 56700, // 15:45
		EvnOpen:  57600, // 16:00
		EvnClose: 72000, // 20:00
		GateOpen: 29400, // 08:10
		Delta:    -14400,
	},
}
