package testmetrics

import (
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/metrics"
)

type Metrics struct {
	Block *metrics.Metric // Number of blocks processed during the test
	TPS   *metrics.Metric // Transactions per second
}

func RegisterMetrics(vu modules.VU) Metrics {
	registry := vu.InitEnv().Registry

	metric := Metrics{
		Block: registry.MustNewMetric("ethereum_block", metrics.Counter, metrics.Default),
		TPS:   registry.MustNewMetric("ethereum_tps", metrics.Trend, metrics.Default),
	}

	return metric
}
