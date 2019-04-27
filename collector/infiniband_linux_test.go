package collector

import (
	"testing"
)

func TestInfiniBandDevices(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	devices, err := infinibandDevices("fixtures/sys/class/infiniband")
	if err != nil {
		t.Fatal(err)
	}
	if l := len(devices); l != 2 {
		t.Fatalf("Retrieved an unexpected number of InfiniBand devices: %d", l)
	}
}
func TestInfiniBandPorts(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ports, err := infinibandPorts("fixtures/sys/class/infiniband", "mlx4_0")
	if err != nil {
		t.Fatal(err)
	}
	if l := len(ports); l != 2 {
		t.Fatalf("Retrieved an unexpected number of InfiniBand ports: %d", l)
	}
}
