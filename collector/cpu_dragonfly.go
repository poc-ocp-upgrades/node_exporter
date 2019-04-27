package collector

import (
	"errors"
	"fmt"
	"unsafe"
	"github.com/prometheus/client_golang/prometheus"
)
import "C"

const maxCPUTimesLen = C.MAXCPU * C.CPUSTATES

type statCollector struct{ cpu *prometheus.Desc }

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("cpu", defaultEnabled, NewStatCollector)
}
func NewStatCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &statCollector{cpu: nodeCPUSecondsDesc}, nil
}
func getDragonFlyCPUTimes() ([]float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		cpuTimesC	*C.uint64_t
		cpuTimerFreq	C.long
		cpuTimesLength	C.size_t
	)
	if C.getCPUTimes(&cpuTimesC, &cpuTimesLength, &cpuTimerFreq) == -1 {
		return nil, errors.New("could not retrieve CPU times")
	}
	defer C.free(unsafe.Pointer(cpuTimesC))
	cput := (*[maxCPUTimesLen]C.uint64_t)(unsafe.Pointer(cpuTimesC))[:cpuTimesLength:cpuTimesLength]
	cpuTimes := make([]float64, cpuTimesLength)
	for i, value := range cput {
		cpuTimes[i] = float64(value) / float64(cpuTimerFreq)
	}
	return cpuTimes, nil
}
func (c *statCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var fieldsCount = 5
	cpuTimes, err := getDragonFlyCPUTimes()
	if err != nil {
		return err
	}
	cpuFields := []string{"user", "nice", "sys", "interrupt", "idle"}
	for i, value := range cpuTimes {
		cpux := fmt.Sprintf("%d", i/fieldsCount)
		ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, value, cpux, cpuFields[i%fieldsCount])
	}
	return nil
}
