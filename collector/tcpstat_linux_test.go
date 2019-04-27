package collector

import (
	"os"
	"testing"
)

func TestTCPStat(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := os.Open("fixtures/proc/net/tcpstat")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	tcpStats, err := parseTCPStats(file)
	if err != nil {
		t.Fatal(err)
	}
	if want, got := 1, int(tcpStats[tcpEstablished]); want != got {
		t.Errorf("want tcpstat number of established state %d, got %d", want, got)
	}
	if want, got := 1, int(tcpStats[tcpListen]); want != got {
		t.Errorf("want tcpstat number of listen state %d, got %d", want, got)
	}
}
