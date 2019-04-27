package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

type conntrackCollector struct {
	current	*prometheus.Desc
	limit	*prometheus.Desc
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("conntrack", defaultEnabled, NewConntrackCollector)
}
func NewConntrackCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &conntrackCollector{current: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "nf_conntrack_entries"), "Number of currently allocated flow entries for connection tracking.", nil, nil), limit: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "nf_conntrack_entries_limit"), "Maximum size of connection tracking table.", nil, nil)}, nil
}
func (c *conntrackCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	value, err := readUintFromFile(procFilePath("sys/net/netfilter/nf_conntrack_count"))
	if err != nil {
		return nil
	}
	ch <- prometheus.MustNewConstMetric(c.current, prometheus.GaugeValue, float64(value))
	value, err = readUintFromFile(procFilePath("sys/net/netfilter/nf_conntrack_max"))
	if err != nil {
		return nil
	}
	ch <- prometheus.MustNewConstMetric(c.limit, prometheus.GaugeValue, float64(value))
	return nil
}
