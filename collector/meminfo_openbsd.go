package collector

import (
	"fmt"
)
import "C"

func (c *meminfoCollector) getMemInfo() (map[string]float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var uvmexp C.struct_uvmexp
	if _, err := C.sysctl_uvmexp(&uvmexp); err != nil {
		return nil, fmt.Errorf("sysctl CTL_VM VM_UVMEXP failed: %v", err)
	}
	ps := float64(uvmexp.pagesize)
	return map[string]float64{"active_bytes": ps * float64(uvmexp.active), "cache_bytes": ps * float64(uvmexp.vnodepages), "free_bytes": ps * float64(uvmexp.free), "inactive_bytes": ps * float64(uvmexp.inactive), "size_bytes": ps * float64(uvmexp.npages), "swap_size_bytes": ps * float64(uvmexp.swpages), "swap_used_bytes": ps * float64(uvmexp.swpginuse), "swapped_in_pages_bytes_total": ps * float64(uvmexp.pgswapin), "swapped_out_pages_bytes_total": ps * float64(uvmexp.pgswapout), "wired_bytes": ps * float64(uvmexp.wired)}, nil
}
