package collector

import "github.com/prometheus/client_golang/prometheus"

type interruptsCollector struct{ desc typedDesc }

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("interrupts", defaultDisabled, NewInterruptsCollector)
}
func NewInterruptsCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &interruptsCollector{desc: typedDesc{prometheus.NewDesc(namespace+"_interrupts_total", "Interrupt details.", interruptLabelNames, nil), prometheus.CounterValue}}, nil
}
