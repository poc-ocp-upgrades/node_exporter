package collector

import (
	"os"
	"testing"
)

func TestNetStats(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	testNetStats(t, "fixtures/proc/net/netstat")
	testSNMP6Stats(t, "fixtures/proc/net/snmp6")
}
func testNetStats(t *testing.T, fileName string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	netStats, err := parseNetStats(file, fileName)
	if err != nil {
		t.Fatal(err)
	}
	if want, got := "102471", netStats["TcpExt"]["DelayedACKs"]; want != got {
		t.Errorf("want netstat TCP DelayedACKs %s, got %s", want, got)
	}
	if want, got := "2786264347", netStats["IpExt"]["OutOctets"]; want != got {
		t.Errorf("want netstat IP OutOctets %s, got %s", want, got)
	}
}
func testSNMP6Stats(t *testing.T, fileName string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	snmp6Stats, err := parseSNMP6Stats(file)
	if err != nil {
		t.Fatal(err)
	}
	if want, got := "460", snmp6Stats["Ip6"]["InOctets"]; want != got {
		t.Errorf("want netstat IPv6 InOctets %s, got %s", want, got)
	}
	if want, got := "8", snmp6Stats["Icmp6"]["OutMsgs"]; want != got {
		t.Errorf("want netstat ICPM6 OutMsgs %s, got %s", want, got)
	}
}
