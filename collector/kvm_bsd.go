package collector

import (
	"fmt"
	"sync"
)
import "C"

type kvm struct {
	mu		sync.Mutex
	hasErr	bool
}

func (k *kvm) SwapUsedPages() (value uint64, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	k.mu.Lock()
	defer k.mu.Unlock()
	if C._kvm_swap_used_pages((*C.uint64_t)(&value)) == -1 {
		k.hasErr = true
		return 0, fmt.Errorf("couldn't get kvm stats")
	}
	return value, nil
}
