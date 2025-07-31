package goetna

import "time"

type EtnaConfig struct {
	RestUrlPub, RestUrlHist, RestUrlPriv string
	WSUrlPub, WSUrlPriv                  string
	RestTimeout                          time.Duration
	WSPingTimeout                        time.Duration
	WSMaxSilentPeriod                    int64 // Maximum period of silence, seconds
}

func defaultConfig() *EtnaConfig {
	return &EtnaConfig{
		RestTimeout: 12000000000,
		RestUrlPub:  "https://pub-api-nvb-demo-prod.etnasoft.us/api/",
		RestUrlHist: "https://back-dev2.nvbrokerage.com/api/",
		// RestUrlPub:        "https://pub-api-nvb-live-prod.etnasoft.us/api/",
		RestUrlPriv:       "https://priv-api-nvb-demo-prod.etnasoft.us/api/",
		WSUrlPub:          "wss://md-str-nvb-demo-prod.etnasoft.us",
		WSUrlPriv:         "wss://oms-str-nvb-demo-prod.etnasoft.us",
		WSPingTimeout:     30,
		WSMaxSilentPeriod: 7000000000,
	}
}

var DefaultConfig = defaultConfig()
