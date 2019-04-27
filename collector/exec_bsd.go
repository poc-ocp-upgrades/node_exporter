package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

type execCollector struct{ sysctls []bsdSysctl }

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("exec", defaultEnabled, NewExecCollector)
}
func NewExecCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &execCollector{sysctls: []bsdSysctl{{name: "exec_context_switches_total", description: "Context switches since system boot.  Resets at architecture unsigned integer.", mib: "vm.stats.sys.v_swtch"}, {name: "exec_traps_total", description: "Traps since system boot.  Resets at architecture unsigned integer.", mib: "vm.stats.sys.v_trap"}, {name: "exec_system_calls_total", description: "System calls since system boot.  Resets at architecture unsigned integer.", mib: "vm.stats.sys.v_syscall"}, {name: "exec_device_interrupts_total", description: "Device interrupts since system boot.  Resets at architecture unsigned integer.", mib: "vm.stats.sys.v_intr"}, {name: "exec_software_interrupts_total", description: "Software interrupts since system boot.  Resets at architecture unsigned integer.", mib: "vm.stats.sys.v_soft"}, {name: "exec_forks_total", description: "Number of fork() calls since system boot.  Resets at architecture unsigned integer.", mib: "vm.stats.vm.v_forks"}}}, nil
}
func (c *execCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, m := range c.sysctls {
		v, err := m.Value()
		if err != nil {
			return err
		}
		ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(namespace+"_"+m.name, m.description, nil, nil), prometheus.CounterValue, v)
	}
	return nil
}
