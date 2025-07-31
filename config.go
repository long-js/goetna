package goetna

import "time"

type EtnaConfig struct {
	RestTimeout       time.Duration
	RestUrl           string
	WSPingTimeout     time.Duration
	WSMaxSilentPeriod int64 // Maximum period of silence, seconds
}

func DefaultConfig() *EtnaConfig {
	return &EtnaConfig{
		RestTimeout:       12000000000,
		RestUrl:           "https://pub-api-nvb-demo-prod.etnasoft.us/api/",
		WSPingTimeout:     30,
		WSMaxSilentPeriod: 7000000000,
	}
}

var defaultConfig = DefaultConfig()
