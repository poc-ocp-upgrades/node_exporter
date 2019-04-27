package collector

import (
	"os"
	"strconv"
	"testing"
)

func TestSockStats(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	testSockStats(t, "fixtures/proc/net/sockstat")
	testSockStats(t, "fixtures/proc/net/sockstat_rhe4")
}
func testSockStats(t *testing.T, fixture string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := os.Open(fixture)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	sockStats, err := parseSockStats(file, fixture)
	if err != nil {
		t.Fatal(err)
	}
	if want, got := "229", sockStats["sockets"]["used"]; want != got {
		t.Errorf("want sockstat sockets used %s, got %s", want, got)
	}
	if want, got := "4", sockStats["TCP"]["tw"]; want != got {
		t.Errorf("want sockstat TCP tw %s, got %s", want, got)
	}
	if want, got := "17", sockStats["TCP"]["alloc"]; want != got {
		t.Errorf("want sockstat TCP alloc %s, got %s", want, got)
	}
	if want, got := strconv.Itoa(os.Getpagesize()), sockStats["TCP"]["mem_bytes"]; want != got {
		t.Errorf("want sockstat TCP mem_bytes %s, got %s", want, got)
	}
}
