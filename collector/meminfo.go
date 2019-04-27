package collector

import (
	"fmt"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	memInfoSubsystem = "memory"
)

type meminfoCollector struct{}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("meminfo", defaultEnabled, NewMeminfoCollector)
}
func NewMeminfoCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &meminfoCollector{}, nil
}
func (c *meminfoCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var metricType prometheus.ValueType
	memInfo, err := c.getMemInfo()
	if err != nil {
		return fmt.Errorf("couldn't get meminfo: %s", err)
	}
	log.Debugf("Set node_mem: %#v", memInfo)
	for k, v := range memInfo {
		if strings.HasSuffix(k, "_total") {
			metricType = prometheus.CounterValue
		} else {
			metricType = prometheus.GaugeValue
		}
		ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, memInfoSubsystem, k), fmt.Sprintf("Memory information field %s.", k), nil, nil), metricType, v)
	}
	return nil
}
