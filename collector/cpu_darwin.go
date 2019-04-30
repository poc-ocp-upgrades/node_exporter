package collector

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"unsafe"
	"github.com/prometheus/client_golang/prometheus"
)
import "C"

const ClocksPerSec = float64(C.CLK_TCK)

type statCollector struct{ cpu *prometheus.Desc }

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("cpu", defaultEnabled, NewCPUCollector)
}
func NewCPUCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &statCollector{cpu: nodeCPUSecondsDesc}, nil
}
func (c *statCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		count	C.mach_msg_type_number_t
		cpuload	*C.processor_cpu_load_info_data_t
		ncpu	C.natural_t
	)
	status := C.host_processor_info(C.host_t(C.mach_host_self()), C.PROCESSOR_CPU_LOAD_INFO, &ncpu, (*C.processor_info_array_t)(unsafe.Pointer(&cpuload)), &count)
	if status != C.KERN_SUCCESS {
		return fmt.Errorf("host_processor_info error=%d", status)
	}
	target := C.vm_map_t(C.mach_task_self_)
	address := C.vm_address_t(uintptr(unsafe.Pointer(cpuload)))
	defer C.vm_deallocate(target, address, C.vm_size_t(ncpu))
	var cpuTicks [C.CPU_STATE_MAX]uint32
	size := int(ncpu) * binary.Size(cpuTicks)
	buf := (*[1 << 30]byte)(unsafe.Pointer(cpuload))[:size:size]
	bbuf := bytes.NewBuffer(buf)
	for i := 0; i < int(ncpu); i++ {
		err := binary.Read(bbuf, binary.LittleEndian, &cpuTicks)
		if err != nil {
			return err
		}
		for k, v := range map[string]int{"user": C.CPU_STATE_USER, "system": C.CPU_STATE_SYSTEM, "nice": C.CPU_STATE_NICE, "idle": C.CPU_STATE_IDLE} {
			ch <- prometheus.MustNewConstMetric(c.cpu, prometheus.CounterValue, float64(cpuTicks[v])/ClocksPerSec, strconv.Itoa(i), k)
		}
	}
	return nil
}
