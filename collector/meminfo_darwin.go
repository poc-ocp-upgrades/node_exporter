package collector

import "C"
import (
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"
	"golang.org/x/sys/unix"
)

func (c *meminfoCollector) getMemInfo() (map[string]float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	infoCount := C.mach_msg_type_number_t(C.HOST_VM_INFO_COUNT)
	vmstat := C.vm_statistics_data_t{}
	ret := C.host_statistics(C.host_t(C.mach_host_self()), C.HOST_VM_INFO, C.host_info_t(unsafe.Pointer(&vmstat)), &infoCount)
	if ret != C.KERN_SUCCESS {
		return nil, fmt.Errorf("Couldn't get memory statistics, host_statistics returned %d", ret)
	}
	totalb, err := unix.Sysctl("hw.memsize")
	if err != nil {
		return nil, err
	}
	total := binary.LittleEndian.Uint64([]byte(totalb + "\x00"))
	ps := float64(C.natural_t(syscall.Getpagesize()))
	return map[string]float64{"active_bytes": ps * float64(vmstat.active_count), "inactive_bytes": ps * float64(vmstat.inactive_count), "wired_bytes": ps * float64(vmstat.wire_count), "free_bytes": ps * float64(vmstat.free_count), "swapped_in_bytes_total": ps * float64(vmstat.pageins), "swapped_out_bytes_total": ps * float64(vmstat.pageouts), "total_bytes": float64(total)}, nil
}
