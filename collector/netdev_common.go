package collector

import (
	"fmt"
	"regexp"
	"strconv"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	netdevIgnoredDevices = kingpin.Flag("collector.netdev.ignored-devices", "Regexp of net devices to ignore for netdev collector.").Default("^$").String()
)

type netDevCollector struct {
	subsystem				string
	ignoredDevicesPattern	*regexp.Regexp
	metricDescs				map[string]*prometheus.Desc
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("netdev", defaultEnabled, NewNetDevCollector)
}
func NewNetDevCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	pattern := regexp.MustCompile(*netdevIgnoredDevices)
	return &netDevCollector{subsystem: "network", ignoredDevicesPattern: pattern, metricDescs: map[string]*prometheus.Desc{}}, nil
}
func (c *netDevCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	netDev, err := getNetDevStats(c.ignoredDevicesPattern)
	if err != nil {
		return fmt.Errorf("couldn't get netstats: %s", err)
	}
	for dev, devStats := range netDev {
		for key, value := range devStats {
			desc, ok := c.metricDescs[key]
			if !ok {
				desc = prometheus.NewDesc(prometheus.BuildFQName(namespace, c.subsystem, key+"_total"), fmt.Sprintf("Network device statistic %s.", key), []string{"device"}, nil)
				c.metricDescs[key] = desc
			}
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid value %s in netstats: %s", value, err)
			}
			ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, v, dev)
		}
	}
	return nil
}
