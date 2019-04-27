package collector

import (
	"strconv"
	"unsafe"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)
import "C"

type cpuCollector struct{ cpu typedDesc }

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("cpu", defaultEnabled, NewCpuCollector)
}
func NewCpuCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &cpuCollector{cpu: typedDesc{nodeCPUSecondsDesc, prometheus.CounterValue}}, nil
}
func (c *cpuCollector) Update(ch chan<- prometheus.Metric) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	clockb, err := unix.SysctlRaw("kern.clockrate")
	if err != nil {
		return err
	}
	clock := *(*C.struct_clockinfo)(unsafe.Pointer(&clockb[0]))
	hz := float64(clock.stathz)
	ncpus, err := unix.SysctlUint32("hw.ncpu")
	if err != nil {
		return err
	}
	var cp_time [][C.CPUSTATES]C.int64_t
	for i := 0; i < int(ncpus); i++ {
		cp_timeb, err := unix.SysctlRaw("kern.cp_time2", i)
		if err != nil && err != unix.ENODEV {
			return err
		}
		if err != unix.ENODEV {
			cp_time = append(cp_time, *(*[C.CPUSTATES]C.int64_t)(unsafe.Pointer(&cp_timeb[0])))
		}
	}
	for cpu, time := range cp_time {
		lcpu := strconv.Itoa(cpu)
		ch <- c.cpu.mustNewConstMetric(float64(time[C.CP_USER])/hz, lcpu, "user")
		ch <- c.cpu.mustNewConstMetric(float64(time[C.CP_NICE])/hz, lcpu, "nice")
		ch <- c.cpu.mustNewConstMetric(float64(time[C.CP_SYS])/hz, lcpu, "system")
		ch <- c.cpu.mustNewConstMetric(float64(time[C.CP_INTR])/hz, lcpu, "interrupt")
		ch <- c.cpu.mustNewConstMetric(float64(time[C.CP_IDLE])/hz, lcpu, "idle")
	}
	return err
}
