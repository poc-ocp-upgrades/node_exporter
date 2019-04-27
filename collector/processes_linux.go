package collector

import (
	"fmt"
	"os"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/procfs"
)

type processCollector struct {
	threadAlloc	*prometheus.Desc
	threadLimit	*prometheus.Desc
	procsState	*prometheus.Desc
	pidUsed		*prometheus.Desc
	pidMax		*prometheus.Desc
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("processes", defaultDisabled, NewProcessStatCollector)
}
func NewProcessStatCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	subsystem := "processes"
	return &processCollector{threadAlloc: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "threads"), "Allocated threads in system", nil, nil), threadLimit: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_threads"), "Limit of threads in the system", nil, nil), procsState: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "state"), "Number of processes in each state.", []string{"state"}, nil), pidUsed: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "pids"), "Number of PIDs", nil, nil), pidMax: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "max_processes"), "Number of max PIDs limit", nil, nil)}, nil
}
func (t *processCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pids, states, threads, err := getAllocatedThreads()
	if err != nil {
		return fmt.Errorf("unable to retrieve number of allocated threads: %q", err)
	}
	ch <- prometheus.MustNewConstMetric(t.threadAlloc, prometheus.GaugeValue, float64(threads))
	maxThreads, err := readUintFromFile(procFilePath("sys/kernel/threads-max"))
	if err != nil {
		return fmt.Errorf("unable to retrieve limit number of threads: %q", err)
	}
	ch <- prometheus.MustNewConstMetric(t.threadLimit, prometheus.GaugeValue, float64(maxThreads))
	for state := range states {
		ch <- prometheus.MustNewConstMetric(t.procsState, prometheus.GaugeValue, float64(states[state]), state)
	}
	pidM, err := readUintFromFile(procFilePath("sys/kernel/pid_max"))
	if err != nil {
		return fmt.Errorf("unable to retrieve limit number of maximum pids alloved: %q", err)
	}
	ch <- prometheus.MustNewConstMetric(t.pidUsed, prometheus.GaugeValue, float64(pids))
	ch <- prometheus.MustNewConstMetric(t.pidMax, prometheus.GaugeValue, float64(pidM))
	return nil
}
func getAllocatedThreads() (int, map[string]int32, int, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return 0, nil, 0, err
	}
	p, err := fs.AllProcs()
	if err != nil {
		return 0, nil, 0, err
	}
	pids := 0
	thread := 0
	procStates := make(map[string]int32)
	for _, pid := range p {
		stat, err := pid.NewStat()
		if os.IsNotExist(err) {
			log.Debugf("file not found when retrieving stats: %q", err)
			continue
		}
		if err != nil {
			return 0, nil, 0, err
		}
		pids++
		procStates[stat.State]++
		thread += stat.NumThreads
	}
	return pids, procStates, thread, nil
}
