package collector

import (
	"runtime"
	"testing"
)

func TestCPU(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		fieldsCount	= 5
		times, err	= getDragonFlyCPUTimes()
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(times) == 0 {
		t.Fatalf("no cputimes found")
	}
	want := runtime.NumCPU() * fieldsCount
	if len(times) != want {
		t.Fatalf("should have %d cpuTimes: got %d", want, len(times))
	}
}
