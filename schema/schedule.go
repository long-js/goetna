package schema

type MarketSchedule struct {
	RegOpen  uint32
	RegClose uint32
	MonOpen  uint32
	MonClose uint32
	EvnOpen  uint32
	EvnClose uint32
	Delta    int32 // UTC delta
}
