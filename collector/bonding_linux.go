package collector

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

type bondingCollector struct{ slaves, active typedDesc }

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("bonding", defaultEnabled, NewBondingCollector)
}
func NewBondingCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &bondingCollector{slaves: typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, "bonding", "slaves"), "Number of configured slaves per bonding interface.", []string{"master"}, nil), prometheus.GaugeValue}, active: typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, "bonding", "active"), "Number of active slaves per bonding interface.", []string{"master"}, nil), prometheus.GaugeValue}}, nil
}
func (c *bondingCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	statusfile := sysFilePath("class/net")
	bondingStats, err := readBondingStats(statusfile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debugf("Not collecting bonding, file does not exist: %s", statusfile)
			return nil
		}
		return err
	}
	for master, status := range bondingStats {
		ch <- c.slaves.mustNewConstMetric(float64(status[0]), master)
		ch <- c.active.mustNewConstMetric(float64(status[1]), master)
	}
	return nil
}
func readBondingStats(root string) (status map[string][2]int, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	status = map[string][2]int{}
	masters, err := ioutil.ReadFile(path.Join(root, "bonding_masters"))
	if err != nil {
		return nil, err
	}
	for _, master := range strings.Fields(string(masters)) {
		slaves, err := ioutil.ReadFile(path.Join(root, master, "bonding", "slaves"))
		if err != nil {
			return nil, err
		}
		sstat := [2]int{0, 0}
		for _, slave := range strings.Fields(string(slaves)) {
			state, err := ioutil.ReadFile(path.Join(root, master, fmt.Sprintf("lower_%s", slave), "operstate"))
			if os.IsNotExist(err) {
				state, err = ioutil.ReadFile(path.Join(root, master, fmt.Sprintf("slave_%s", slave), "operstate"))
			}
			if err != nil {
				return nil, err
			}
			sstat[0]++
			if strings.TrimSpace(string(state)) == "up" {
				sstat[1]++
			}
		}
		status[master] = sstat
	}
	return status, err
}
