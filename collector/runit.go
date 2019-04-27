package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/soundcloud/go-runit/runit"
	"gopkg.in/alecthomas/kingpin.v2"
)

var runitServiceDir = kingpin.Flag("collector.runit.servicedir", "Path to runit service directory.").Default("/etc/service").String()

type runitCollector struct{ state, stateDesired, stateNormal, stateTimestamp typedDesc }

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("runit", defaultDisabled, NewRunitCollector)
}
func NewRunitCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		subsystem	= "service"
		constLabels	= prometheus.Labels{"supervisor": "runit"}
		labelNames	= []string{"service"}
	)
	return &runitCollector{state: typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "state"), "State of runit service.", labelNames, constLabels), prometheus.GaugeValue}, stateDesired: typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "desired_state"), "Desired state of runit service.", labelNames, constLabels), prometheus.GaugeValue}, stateNormal: typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "normal_state"), "Normal state of runit service.", labelNames, constLabels), prometheus.GaugeValue}, stateTimestamp: typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "state_last_change_timestamp_seconds"), "Unix timestamp of the last runit service state change.", labelNames, constLabels), prometheus.GaugeValue}}, nil
}
func (c *runitCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	services, err := runit.GetServices(*runitServiceDir)
	if err != nil {
		return err
	}
	for _, service := range services {
		status, err := service.Status()
		if err != nil {
			log.Debugf("Couldn't get status for %s: %s, skipping...", service.Name, err)
			continue
		}
		log.Debugf("%s is %d on pid %d for %d seconds", service.Name, status.State, status.Pid, status.Duration)
		ch <- c.state.mustNewConstMetric(float64(status.State), service.Name)
		ch <- c.stateDesired.mustNewConstMetric(float64(status.Want), service.Name)
		ch <- c.stateTimestamp.mustNewConstMetric(float64(status.Timestamp.Unix()), service.Name)
		if status.NormallyUp {
			ch <- c.stateNormal.mustNewConstMetric(1, service.Name)
		} else {
			ch <- c.stateNormal.mustNewConstMetric(0, service.Name)
		}
	}
	return nil
}
