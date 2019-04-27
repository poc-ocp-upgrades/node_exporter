package collector

import (
	"time"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type timeCollector struct{ desc *prometheus.Desc }

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("time", defaultEnabled, NewTimeCollector)
}
func NewTimeCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &timeCollector{desc: prometheus.NewDesc(namespace+"_time_seconds", "System time in seconds since epoch (1970).", nil, nil)}, nil
}
func (c *timeCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	now := float64(time.Now().UnixNano()) / 1e9
	log.Debugf("Return time: %f", now)
	ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, now)
	return nil
}
