package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)

type entropyCollector struct{ entropyAvail *prometheus.Desc }

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("entropy", defaultEnabled, NewEntropyCollector)
}
func NewEntropyCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &entropyCollector{entropyAvail: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "entropy_available_bits"), "Bits of available entropy.", nil, nil)}, nil
}
func (c *entropyCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	value, err := readUintFromFile(procFilePath("sys/kernel/random/entropy_avail"))
	if err != nil {
		return fmt.Errorf("couldn't get entropy_avail: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(c.entropyAvail, prometheus.GaugeValue, float64(value))
	return nil
}
