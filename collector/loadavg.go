package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type loadavgCollector struct{ metric []typedDesc }

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("loadavg", defaultEnabled, NewLoadavgCollector)
}
func NewLoadavgCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &loadavgCollector{metric: []typedDesc{{prometheus.NewDesc(namespace+"_load1", "1m load average.", nil, nil), prometheus.GaugeValue}, {prometheus.NewDesc(namespace+"_load5", "5m load average.", nil, nil), prometheus.GaugeValue}, {prometheus.NewDesc(namespace+"_load15", "15m load average.", nil, nil), prometheus.GaugeValue}}}, nil
}
func (c *loadavgCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	loads, err := getLoad()
	if err != nil {
		return fmt.Errorf("couldn't get load: %s", err)
	}
	for i, load := range loads {
		log.Debugf("return load %d: %f", i, load)
		ch <- c.metric[i].mustNewConstMetric(load)
	}
	return err
}
