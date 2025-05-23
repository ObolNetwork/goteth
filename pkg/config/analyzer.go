package config

import (
	"github.com/attestantio/go-eth2-client/spec/phase0"
	cli "github.com/urfave/cli/v2"
)

type AnalyzerConfig struct {
	LogLevel                 string      `json:"log-level"`
	InitSlot                 phase0.Slot `json:"init-slot"`
	FinalSlot                phase0.Slot `json:"final-slot"`
	RewardsAggregationEpochs int         `json:"rewards-aggregation-epochs"`
	BnEndpoint               string      `json:"bn-endpoint"`
	BnApiKey                 string      `json:"bn-api-key"`
	CfAccessClientID         string      `json:"cf-access-client-id"`
	CfAccessClientSecret     string      `json:"cf-access-client-secret"`
	ElEndpoint               string      `json:"el-endpoint"`
	DBUrl                    string      `json:"db-url"`
	DownloadMode             string      `json:"download-mode"`
	WorkerNum                int         `json:"worker-num"`
	DbWorkerNum              int         `json:"db-worker-num"`
	Metrics                  string      `json:"metrics"`
	PrometheusPort           int         `json:"prometheus-port"`
	MaxRequestRetries        int         `json:"max-request-retries"`
	BeaconContractAddress    string      `json:"beacon-contract-address"`
}

// TODO: read from config-file
func NewAnalyzerConfig() *AnalyzerConfig {
	// Return Default values for the ethereum configuration
	return &AnalyzerConfig{
		LogLevel:                 DefaultLogLevel,
		InitSlot:                 phase0.Slot(DefaultInitSlot),
		FinalSlot:                phase0.Slot(DefaultFinalSlot),
		RewardsAggregationEpochs: DefaultRewardsAggregationEpochs,
		BnEndpoint:               DefaultBnEndpoint,
		BnApiKey:                 DefaultBnApiKey,
		CfAccessClientID:         DefaultCfAccessClientID,
		CfAccessClientSecret:     DefaultCfAccessClientSecret,
		ElEndpoint:               DefaultElEndpoint,
		DBUrl:                    DefaultDBUrl,
		DownloadMode:             DefaultDownloadMode,
		WorkerNum:                DefaultWorkerNum,
		DbWorkerNum:              DefaultDbWorkerNum,
		Metrics:                  DefaultMetrics,
		PrometheusPort:           DefaultPrometheusPort,
		MaxRequestRetries:        DefaultMaxRequestRetries,
		BeaconContractAddress:    DefaultBeaconContractAddress,
	}
}

func (c *AnalyzerConfig) Apply(ctx *cli.Context) {
	// apply to the existing Default configuration the set flags
	// log level
	if ctx.IsSet("log-level") {
		c.LogLevel = ctx.String("log-level")
	}
	// init slot
	if ctx.IsSet("init-slot") {
		c.InitSlot = phase0.Slot(ctx.Int("init-slot"))
	}
	// final slot
	if ctx.IsSet("final-slot") {
		c.FinalSlot = phase0.Slot(ctx.Int("final-slot"))
	}
	// rewards aggregation epochs
	if ctx.IsSet("rewards-aggregation-epochs") {
		c.RewardsAggregationEpochs = ctx.Int("rewards-aggregation-epochs")
	}
	// cl url
	if ctx.IsSet("bn-endpoint") {
		c.BnEndpoint = ctx.String("bn-endpoint")
	}
	// bn api key
	if ctx.IsSet("bn-api-key") {
		c.BnApiKey = ctx.String("bn-api-key")
	}

	// cf access client id
	if ctx.IsSet("cf-access-client-id") {
		c.CfAccessClientID = ctx.String("cf-access-client-id")
	}

	// cf access client secret
	if ctx.IsSet("cf-access-client-secret") {
		c.CfAccessClientSecret = ctx.String("cf-access-client-secret")
	}

	// el url
	if ctx.IsSet("el-endpoint") {
		c.ElEndpoint = ctx.String("el-endpoint")
	}
	// db url
	if ctx.IsSet("db-url") {
		c.DBUrl = ctx.String("db-url")
	}
	// download mode
	if ctx.IsSet("download-mode") {
		c.DownloadMode = ctx.String("download-mode")
	}
	// worker num
	if ctx.IsSet("workers-num") {
		c.WorkerNum = ctx.Int("workers-num")
	}
	// db worker num
	if ctx.IsSet("db-workers-num") {
		c.DbWorkerNum = ctx.Int("db-workers-num")
	}
	// metrics
	if ctx.IsSet("metrics") {
		c.Metrics = ctx.String("metrics")
	}
	// prometheus port
	if ctx.IsSet("prometheus-port") {
		c.PrometheusPort = ctx.Int("prometheus-port")
	}
	// max request retries
	if ctx.IsSet("max-request-retries") {
		c.MaxRequestRetries = ctx.Int("max-request-retries")
	}
	// beacon contract address
	if ctx.IsSet("beacon-contract-address") {
		c.BeaconContractAddress = ctx.String("beacon-contract-address")
	}
}
