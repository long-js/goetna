package goetna

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type EtnaConfig struct {
	RestUrlPub, RestUrlNonRTH, RestUrlPriv string
	WSUrlPub, WSUrlPubFMP, WSUrlPriv       string
	RestTimeout                            time.Duration
	WSPingTimeout                          time.Duration
	WSMaxSilentPeriod                      int64 // Maximum period of silence, seconds
}

func defaultConfig() *EtnaConfig {
	_ = godotenv.Load()
	cfg := EtnaConfig{
		RestTimeout: 12000000000, WSPingTimeout: 30, WSMaxSilentPeriod: 7000000000,
		RestUrlNonRTH: "https://back-dev2.nvbrokerage.com/api/",
	}
	if isTest, err := strconv.ParseBool(os.Getenv("TEST_ENV")); err != nil || !isTest {
		cfg.RestUrlPub = "https://pub-api-nvb-live-prod.etnasoft.us/api/"
		cfg.RestUrlPriv = "https://priv-api-nvb-live-prod.etnasoft.us/api/"
		cfg.WSUrlPub = "wss://md-str-nvb-live-prod.etnasoft.us"
		cfg.WSUrlPubFMP = "wss://websockets.financialmodelingprep.com"
		cfg.WSUrlPriv = "wss://oms-str-nvb-live-prod.etnasoft.us"
	} else {
		cfg.RestUrlPub = "https://pub-api-nvb-demo-prod.etnasoft.us/api/"
		cfg.RestUrlPriv = "https://priv-api-nvb-demo-prod.etnasoft.us/api/"
		cfg.WSUrlPub = "wss://md-str-nvb-demo-prod.etnasoft.us"
		cfg.WSUrlPubFMP = "wss://websockets.financialmodelingprep.com"
		cfg.WSUrlPriv = "wss://oms-str-nvb-demo-prod.etnasoft.us"
	}
	return &cfg
}

var DefaultConfig = defaultConfig()
