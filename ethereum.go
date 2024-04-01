package ethereum

import (
	"xk6-eth/testmetrics"

	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/ethereum", &RootModule{})
}

type (
	RootModule struct{}
	Module     struct {
		vu      modules.VU
		metrics testmetrics.Metrics
	}
)

func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &Module{
		vu:      vu,
		metrics: testmetrics.RegisterMetrics(vu),
	}
}

func (module *Module) Exports() modules.Exports {
	return modules.Exports{Named: map[string]interface{}{
		"Premine": module.Premine,
		"Client":  module.NewClient,
	}}
}
