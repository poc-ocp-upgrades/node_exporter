package collector

import (
	"fmt"
	"strconv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/procfs"
)

const (
	buddyInfoSubsystem = "buddyinfo"
)

type buddyinfoCollector struct{ desc *prometheus.Desc }

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("buddyinfo", defaultDisabled, NewBuddyinfoCollector)
}
func NewBuddyinfoCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	desc := prometheus.NewDesc(prometheus.BuildFQName(namespace, buddyInfoSubsystem, "blocks"), "Count of free blocks according to size.", []string{"node", "zone", "size"}, nil)
	return &buddyinfoCollector{desc}, nil
}
func (c *buddyinfoCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return fmt.Errorf("failed to open procfs: %v", err)
	}
	buddyInfo, err := fs.NewBuddyInfo()
	if err != nil {
		return fmt.Errorf("couldn't get buddyinfo: %s", err)
	}
	log.Debugf("Set node_buddy: %#v", buddyInfo)
	for _, entry := range buddyInfo {
		for size, value := range entry.Sizes {
			ch <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, value, entry.Node, entry.Zone, strconv.Itoa(size))
		}
	}
	return nil
}
