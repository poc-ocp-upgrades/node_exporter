package collector

import (
	"fmt"
	"path"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ksmdFiles = []string{"full_scans", "merge_across_nodes", "pages_shared", "pages_sharing", "pages_to_scan", "pages_unshared", "pages_volatile", "run", "sleep_millisecs"}
)

type ksmdCollector struct{ metricDescs map[string]*prometheus.Desc }

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("ksmd", defaultDisabled, NewKsmdCollector)
}
func getCanonicalMetricName(filename string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch filename {
	case "full_scans":
		return filename + "_total"
	case "sleep_millisecs":
		return "sleep_seconds"
	default:
		return filename
	}
}
func NewKsmdCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	subsystem := "ksmd"
	descs := make(map[string]*prometheus.Desc)
	for _, n := range ksmdFiles {
		descs[n] = prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, getCanonicalMetricName(n)), fmt.Sprintf("ksmd '%s' file.", n), nil, nil)
	}
	return &ksmdCollector{descs}, nil
}
func (c *ksmdCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, n := range ksmdFiles {
		val, err := readUintFromFile(sysFilePath(path.Join("kernel/mm/ksm", n)))
		if err != nil {
			return err
		}
		t := prometheus.GaugeValue
		v := float64(val)
		switch n {
		case "full_scans":
			t = prometheus.CounterValue
		case "sleep_millisecs":
			v /= 1000
		}
		ch <- prometheus.MustNewConstMetric(c.metricDescs[n], t, v)
	}
	return nil
}
