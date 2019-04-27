package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

type bootTimeCollector struct{ boottime bsdSysctl }

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("boottime", defaultEnabled, newBootTimeCollector)
}
func newBootTimeCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &bootTimeCollector{boottime: bsdSysctl{name: "boot_time_seconds", description: "Unix time of last boot, including microseconds.", mib: "kern.boottime", dataType: bsdSysctlTypeStructTimeval}}, nil
}
func (c *bootTimeCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	v, err := c.boottime.Value()
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, "", c.boottime.name), c.boottime.description, nil, nil), prometheus.GaugeValue, v)
	return nil
}
