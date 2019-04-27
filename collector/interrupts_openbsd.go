package collector

import (
	"fmt"
	"strconv"
	"github.com/prometheus/client_golang/prometheus"
)
import "C"

var (
	interruptLabelNames = []string{"cpu", "type", "devices"}
)

func (c *interruptsCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	interrupts, err := getInterrupts()
	if err != nil {
		return fmt.Errorf("couldn't get interrupts: %s", err)
	}
	for dev, interrupt := range interrupts {
		for cpuNo, value := range interrupt.values {
			ch <- c.desc.mustNewConstMetric(value, strconv.Itoa(cpuNo), fmt.Sprintf("%d", interrupt.vector), dev)
		}
	}
	return nil
}

type interrupt struct {
	vector	int
	device	string
	values	[]float64
}

func getInterrupts() (map[string]interrupt, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		cintr		C.struct_intr
		interrupts	= map[string]interrupt{}
	)
	nintr := C.sysctl_nintr()
	for i := C.int(0); i < nintr; i++ {
		_, err := C.sysctl_intr(&cintr, i)
		if err != nil {
			return nil, err
		}
		dev := C.GoString(&cintr.device[0])
		interrupts[dev] = interrupt{vector: int(cintr.vector), device: dev, values: []float64{float64(cintr.count)}}
	}
	return interrupts, nil
}
