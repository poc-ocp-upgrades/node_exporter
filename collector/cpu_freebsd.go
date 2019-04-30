package collector

import (
	"fmt"
	"math"
	"strconv"
	"unsafe"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"golang.org/x/sys/unix"
)

type clockinfo struct {
	hz	int32
	tick	int32
	spare	int32
	stathz	int32
	profhz	int32
}
type cputime struct {
	user	float64
	nice	float64
	sys	float64
	intr	float64
	idle	float64
}

func getCPUTimes() ([]cputime, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	const states = 5
	clockb, err := unix.SysctlRaw("kern.clockrate")
	if err != nil {
		return nil, err
	}
	clock := *(*clockinfo)(unsafe.Pointer(&clockb[0]))
	cpb, err := unix.SysctlRaw("kern.cp_times")
	if err != nil {
		return nil, err
	}
	var cpufreq float64
	if clock.stathz > 0 {
		cpufreq = float64(clock.stathz)
	} else {
		cpufreq = float64(clock.hz)
	}
	var times []float64
	for len(cpb) >= int(unsafe.Sizeof(int(0))) {
		t := *(*int)(unsafe.Pointer(&cpb[0]))
		times = append(times, float64(t)/cpufreq)
		cpb = cpb[unsafe.Sizeof(int(0)):]
	}
	cpus := make([]cputime, len(times)/states)
	for i := 0; i < len(times); i += states {
		cpu := &cpus[i/states]
		cpu.user = times[i]
		cpu.nice = times[i+1]
		cpu.sys = times[i+2]
		cpu.intr = times[i+3]
		cpu.idle = times[i+4]
	}
	return cpus, nil
}

type statCollector struct {
	cpu	typedDesc
	temp	typedDesc
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("cpu", defaultEnabled, NewStatCollector)
}
func NewStatCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &statCollector{cpu: typedDesc{nodeCPUSecondsDesc, prometheus.CounterValue}, temp: typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "temperature_celsius"), "CPU temperature", []string{"cpu"}, nil), prometheus.GaugeValue}}, nil
}
func (c *statCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	cpuTimes, err := getCPUTimes()
	if err != nil {
		return err
	}
	for cpu, t := range cpuTimes {
		lcpu := strconv.Itoa(cpu)
		ch <- c.cpu.mustNewConstMetric(float64(t.user), lcpu, "user")
		ch <- c.cpu.mustNewConstMetric(float64(t.nice), lcpu, "nice")
		ch <- c.cpu.mustNewConstMetric(float64(t.sys), lcpu, "system")
		ch <- c.cpu.mustNewConstMetric(float64(t.intr), lcpu, "interrupt")
		ch <- c.cpu.mustNewConstMetric(float64(t.idle), lcpu, "idle")
		temp, err := unix.SysctlUint32(fmt.Sprintf("dev.cpu.%d.temperature", cpu))
		if err != nil {
			if err == unix.ENOENT {
				log.Debugf("no temperature information for CPU %d", cpu)
			} else {
				ch <- c.temp.mustNewConstMetric(math.NaN(), lcpu)
				log.Errorf("failed to query CPU temperature for CPU %d: %s", cpu, err)
			}
			continue
		}
		ch <- c.temp.mustNewConstMetric(float64(int32(temp)-2732)/10, lcpu)
	}
	return err
}
