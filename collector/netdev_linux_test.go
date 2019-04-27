package collector

import (
	"os"
	"regexp"
	"testing"
)

func TestNetDevStats(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := os.Open("fixtures/proc/net/dev")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	netStats, err := parseNetDevStats(file, regexp.MustCompile("^veth"))
	if err != nil {
		t.Fatal(err)
	}
	if want, got := "10437182923", netStats["wlan0"]["receive_bytes"]; want != got {
		t.Errorf("want netstat wlan0 bytes %s, got %s", want, got)
	}
	if want, got := "68210035552", netStats["eth0"]["receive_bytes"]; want != got {
		t.Errorf("want netstat eth0 bytes %s, got %s", want, got)
	}
	if want, got := "934", netStats["tun0"]["transmit_packets"]; want != got {
		t.Errorf("want netstat tun0 packets %s, got %s", want, got)
	}
	if want, got := 9, len(netStats); want != got {
		t.Errorf("want count of devices to be %d, got %d", want, got)
	}
	if _, ok := netStats["veth4B09XN"]["transmit_bytes"]; ok {
		t.Error("want fixture interface veth4B09XN to not exist, but it does")
	}
	if want, got := "0", netStats["ibr10:30"]["receive_fifo"]; want != got {
		t.Error("want fixture interface ibr10:30 to exist, but it does not")
	}
	if want, got := "72", netStats["💩0"]["receive_multicast"]; want != got {
		t.Error("want fixture interface 💩0 to exist, but it does not")
	}
}
