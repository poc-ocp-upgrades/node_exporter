package collector

import (
	"fmt"
	"github.com/prometheus/procfs"
	"github.com/prometheus/client_golang/prometheus"
)

type statCollector struct {
	intr		*prometheus.Desc
	ctxt		*prometheus.Desc
	forks		*prometheus.Desc
	btime		*prometheus.Desc
	procsRunning	*prometheus.Desc
	procsBlocked	*prometheus.Desc
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("stat", defaultEnabled, NewStatCollector)
}
func NewStatCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &statCollector{intr: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "intr_total"), "Total number of interrupts serviced.", nil, nil), ctxt: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "context_switches_total"), "Total number of context switches.", nil, nil), forks: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "forks_total"), "Total number of forks.", nil, nil), btime: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "boot_time_seconds"), "Node boot time, in unixtime.", nil, nil), procsRunning: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "procs_running"), "Number of processes in runnable state.", nil, nil), procsBlocked: prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "procs_blocked"), "Number of processes blocked waiting for I/O to complete.", nil, nil)}, nil
}
func (c *statCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return fmt.Errorf("failed to open procfs: %v", err)
	}
	stats, err := fs.NewStat()
	if err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(c.intr, prometheus.CounterValue, float64(stats.IRQTotal))
	ch <- prometheus.MustNewConstMetric(c.ctxt, prometheus.CounterValue, float64(stats.ContextSwitches))
	ch <- prometheus.MustNewConstMetric(c.forks, prometheus.CounterValue, float64(stats.ProcessCreated))
	ch <- prometheus.MustNewConstMetric(c.btime, prometheus.GaugeValue, float64(stats.BootTime))
	ch <- prometheus.MustNewConstMetric(c.procsRunning, prometheus.GaugeValue, float64(stats.ProcessesRunning))
	ch <- prometheus.MustNewConstMetric(c.procsBlocked, prometheus.GaugeValue, float64(stats.ProcessesBlocked))
	return nil
}
