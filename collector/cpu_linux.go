package collector

import (
	"fmt"
	"path/filepath"
	"strconv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/sysfs"
)

type cpuCollector struct {
	cpu			*prometheus.Desc
	cpuGuest		*prometheus.Desc
	cpuFreq			*prometheus.Desc
	cpuFreqMin		*prometheus.Desc
	cpuFreqMax		*prometheus.Desc
	cpuCoreThrottle		*prometheus.Desc
	cpuPackageThrottle	*prometheus.Desc
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("cpu", defaultEnabled, NewCPUCollector)
}
func NewCPUCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &cpuCollector{cpu: nodeCPUSecondsDesc, cpuGuest: prometheus.NewDesc(prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "guest_seconds_total"), "Seconds the cpus spent in guests (VMs) for each mode.", []string{"cpu", "mode"}, nil), cpuFreq: prometheus.NewDesc(prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_hertz"), "Current cpu thread frequency in hertz.", []string{"cpu"}, nil), cpuFreqMin: prometheus.NewDesc(prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_min_hertz"), "Minimum cpu thread frequency in hertz.", []string{"cpu"}, nil), cpuFreqMax: prometheus.NewDesc(prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "frequency_max_hertz"), "Maximum cpu thread frequency in hertz.", []string{"cpu"}, nil), cpuCoreThrottle: prometheus.NewDesc(prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "core_throttles_total"), "Number of times this cpu core has been throttled.", []string{"package", "core"}, nil), cpuPackageThrottle: prometheus.NewDesc(prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "package_throttles_total"), "Number of times this cpu package has been throttled.", []string{"package"}, nil)}, nil
}
func (c *cpuCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if err := c.updateStat(ch); err != nil {
		return err
	}
	if err := c.updateCPUfreq(ch); err != nil {
		return err
	}
	if err := c.updateThermalThrottle(ch); err != nil {
		return err
	}
	return nil
}
func (c *cpuCollector) updateCPUfreq(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return fmt.Errorf("failed to open sysfs: %v", err)
	}
	cpuFreqs, err := fs.NewSystemCpufreq()
	if err != nil {
		return err
	}
	for _, stats := range cpuFreqs {
		ch <- prometheus.MustNewConstMetric(c.cpuFreq, prometheus.GaugeValue, float64(stats.CurrentFrequency)*1000.0, stats.Name)
		ch <- prometheus.MustNewConstMetric(c.cpuFreqMin, prometheus.GaugeValue, float64(stats.MinimumFrequency)*1000.0, stats.Name)
		ch <- prometheus.MustNewConstMetric(c.cpuFreqMax, prometheus.GaugeValue, float64(stats.MaximumFrequency)*1000.0, stats.Name)
	}
	return nil
}
func (c *cpuCollector) updateThermalThrottle(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	cpus, err := filepath.Glob(sysFilePath("devices/system/cpu/cpu[0-9]*"))
	if err != nil {
		return err
	}
	packageThrottles := make(map[uint64]uint64)
	packageCoreThrottles := make(map[uint64]map[uint64]uint64)
	for _, cpu := range cpus {
		var err error
		var physicalPackageID, coreID uint64
		if physicalPackageID, err = readUintFromFile(filepath.Join(cpu, "topology", "physical_package_id")); err != nil {
			log.Debugf("CPU %v is missing physical_package_id", cpu)
			continue
		}
		if coreID, err = readUintFromFile(filepath.Join(cpu, "topology", "core_id")); err != nil {
			log.Debugf("CPU %v is missing core_id", cpu)
			continue
		}
		if _, present := packageCoreThrottles[physicalPackageID]; !present {
			packageCoreThrottles[physicalPackageID] = make(map[uint64]uint64)
		}
		if _, present := packageCoreThrottles[physicalPackageID][coreID]; !present {
			if coreThrottleCount, err := readUintFromFile(filepath.Join(cpu, "thermal_throttle", "core_throttle_count")); err == nil {
				packageCoreThrottles[physicalPackageID][coreID] = coreThrottleCount
			} else {
				log.Debugf("CPU %v is missing core_throttle_count", cpu)
			}
		}
		if _, present := packageThrottles[physicalPackageID]; !present {
			if packageThrottleCount, err := readUintFromFile(filepath.Join(cpu, "thermal_throttle", "package_throttle_count")); err == nil {
				packageThrottles[physicalPackageID] = packageThrottleCount
			} else {
				log.Debugf("CPU %v is missing package_throttle_count", cpu)
			}
		}
	}
	for physicalPackageID, packageThrottleCount := range packageThrottles {
		ch <- prometheus.MustNewConstMetric(c.cpuPackageThrottle, prometheus.CounterValue, float64(packageThrottleCount), strconv.FormatUint(physicalPackageID, 10))
	}
	for physicalPackageID, coreMap := range packageCoreThrottles {
		for coreID, coreThrottleCount := range coreMap {
			ch <- prometheus.MustNewConstMetric(c.cpuCoreThrottle, prometheus.CounterValue, float64(coreThrottleCount), strconv.FormatUint(physicalPackageID, 10), strconv.FormatUint(coreID, 10))
		}
	}
	return nil
}
func (c *cpuCollector) updateStat(ch chan<- prometheus.Metric) error {
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
	for cpuID, cpuStat := range stats.CPU {
		cpuNum := fmt.Sprintf("%d", cpuID)
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.User, cpuNum, "user")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.Nice, cpuNum, "nice")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.System, cpuNum, "system")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.Idle, cpuNum, "idle")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.Iowait, cpuNum, "iowait")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.IRQ, cpuNum, "irq")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.SoftIRQ, cpuNum, "softirq")
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, cpuStat.Steal, cpuNum, "steal")
		ch <- prometheus.MustNewConstMetric(c.cpuGuest, prometheus.CounterValue, cpuStat.Guest, cpuNum, "user")
		ch <- prometheus.MustNewConstMetric(c.cpuGuest, prometheus.CounterValue, cpuStat.GuestNice, cpuNum, "nice")
	}
	return nil
}
