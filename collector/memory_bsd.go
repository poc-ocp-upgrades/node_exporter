package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

const (
	memorySubsystem = "memory"
)

type memoryCollector struct {
	pageSize	uint64
	sysctls		[]bsdSysctl
	kvm			kvm
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("meminfo", defaultEnabled, NewMemoryCollector)
}
func NewMemoryCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tmp32, err := unix.SysctlUint32("vm.stats.vm.v_page_size")
	if err != nil {
		return nil, fmt.Errorf("sysctl(vm.stats.vm.v_page_size) failed: %s", err)
	}
	size := float64(tmp32)
	fromPage := func(v float64) float64 {
		return v * size
	}
	return &memoryCollector{pageSize: uint64(tmp32), sysctls: []bsdSysctl{{name: "active_bytes", description: "Recently used by userland", mib: "vm.stats.vm.v_active_count", conversion: fromPage}, {name: "inactive_bytes", description: "Not recently used by userland", mib: "vm.stats.vm.v_inactive_count", conversion: fromPage}, {name: "wired_bytes", description: "Locked in memory by kernel, mlock, etc", mib: "vm.stats.vm.v_wire_count", conversion: fromPage}, {name: "cache_bytes", description: "Almost free, backed by swap or files, available for re-allocation", mib: "vm.stats.vm.v_cache_count", conversion: fromPage}, {name: "buffer_bytes", description: "Disk IO Cache entries for non ZFS filesystems, only usable by kernel", mib: "vfs.bufspace", dataType: bsdSysctlTypeCLong}, {name: "free_bytes", description: "Unallocated, available for allocation", mib: "vm.stats.vm.v_free_count", conversion: fromPage}, {name: "size_bytes", description: "Total physical memory size", mib: "vm.stats.vm.v_page_count", conversion: fromPage}, {name: "swap_size_bytes", description: "Total swap memory size", mib: "vm.swap_total", dataType: bsdSysctlTypeUint64}, {name: "swap_in_bytes_total", description: "Bytes paged in from swap devices", mib: "vm.stats.vm.v_swappgsin", valueType: prometheus.CounterValue, conversion: fromPage}, {name: "swap_out_bytes_total", description: "Bytes paged out to swap devices", mib: "vm.stats.vm.v_swappgsout", valueType: prometheus.CounterValue, conversion: fromPage}}}, nil
}
func (c *memoryCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, m := range c.sysctls {
		v, err := m.Value()
		if err != nil {
			return fmt.Errorf("couldn't get memory: %s", err)
		}
		if m.valueType == 0 {
			m.valueType = prometheus.GaugeValue
		}
		ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, memorySubsystem, m.name), m.description, nil, nil), m.valueType, v)
	}
	swapUsed, err := c.kvm.SwapUsedPages()
	if err != nil {
		return fmt.Errorf("couldn't get kvm: %s", err)
	}
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, memorySubsystem, "swap_used_bytes"), "Currently allocated swap", nil, nil), prometheus.GaugeValue, float64(swapUsed*c.pageSize))
	return nil
}
