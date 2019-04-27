package collector

import (
	"fmt"
	"syscall"
	"unsafe"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)
import "C"

type bsdSysctlType uint8

const (
	bsdSysctlTypeUint32	bsdSysctlType	= iota
	bsdSysctlTypeUint64
	bsdSysctlTypeStructTimeval
	bsdSysctlTypeCLong
)

type bsdSysctl struct {
	name		string
	description	string
	valueType	prometheus.ValueType
	mib		string
	dataType	bsdSysctlType
	conversion	func(float64) float64
}

func (b bsdSysctl) Value() (float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var tmp32 uint32
	var tmp64 uint64
	var tmpf64 float64
	var err error
	switch b.dataType {
	case bsdSysctlTypeUint32:
		tmp32, err = unix.SysctlUint32(b.mib)
		tmpf64 = float64(tmp32)
	case bsdSysctlTypeUint64:
		tmp64, err = unix.SysctlUint64(b.mib)
		tmpf64 = float64(tmp64)
	case bsdSysctlTypeStructTimeval:
		tmpf64, err = b.getStructTimeval()
	case bsdSysctlTypeCLong:
		tmpf64, err = b.getCLong()
	}
	if err != nil {
		return 0, err
	}
	if b.conversion != nil {
		return b.conversion(tmpf64), nil
	}
	return tmpf64, nil
}
func (b bsdSysctl) getStructTimeval() (float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	raw, err := unix.SysctlRaw(b.mib)
	if err != nil {
		return 0, err
	}
	if len(raw) != int(unsafe.Sizeof(syscall.Timeval{})) {
		return 0, fmt.Errorf("length of bytes received from sysctl (%d) does not match expected bytes (%d)", len(raw), unsafe.Sizeof(syscall.Timeval{}))
	}
	tv := *(*syscall.Timeval)(unsafe.Pointer(&raw[0]))
	return (float64(tv.Sec) + (float64(tv.Usec) / float64(1000*1000))), nil
}
func (b bsdSysctl) getCLong() (float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	raw, err := unix.SysctlRaw(b.mib)
	if err != nil {
		return 0, err
	}
	if len(raw) == C.sizeof_long {
		return float64(*(*C.long)(unsafe.Pointer(&raw[0]))), nil
	}
	if len(raw) == C.sizeof_int {
		return float64(*(*C.int)(unsafe.Pointer(&raw[0]))), nil
	}
	return 0, fmt.Errorf("length of bytes received from sysctl (%d) does not match expected bytes (long: %d), (int: %d)", len(raw), C.sizeof_long, C.sizeof_int)
}
