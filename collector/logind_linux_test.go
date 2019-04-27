package collector

import (
	"testing"
	"github.com/godbus/dbus"
	"github.com/prometheus/client_golang/prometheus"
)

type testLogindInterface struct{}

var testSeats = []string{"seat0", ""}

func (c *testLogindInterface) listSeats() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return testSeats, nil
}
func (c *testLogindInterface) listSessions() ([]logindSessionEntry, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return []logindSessionEntry{{SessionID: "1", UserID: 0, UserName: "", SeatID: "", SessionObjectPath: dbus.ObjectPath("/org/freedesktop/login1/session/1")}, {SessionID: "2", UserID: 0, UserName: "", SeatID: "seat0", SessionObjectPath: dbus.ObjectPath("/org/freedesktop/login1/session/2")}}, nil
}
func (c *testLogindInterface) getSession(session logindSessionEntry) *logindSession {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	sessions := map[dbus.ObjectPath]*logindSession{dbus.ObjectPath("/org/freedesktop/login1/session/1"): {seat: session.SeatID, remote: "true", sessionType: knownStringOrOther("tty", attrTypeValues), class: knownStringOrOther("user", attrClassValues)}, dbus.ObjectPath("/org/freedesktop/login1/session/2"): {seat: session.SeatID, remote: "false", sessionType: knownStringOrOther("x11", attrTypeValues), class: knownStringOrOther("greeter", attrClassValues)}}
	return sessions[session.SessionObjectPath]
}
func TestLogindCollectorKnownStringOrOther(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	known := []string{"foo", "bar"}
	actual := knownStringOrOther("foo", known)
	expected := "foo"
	if actual != expected {
		t.Errorf("knownStringOrOther failed: got %q, expected %q.", actual, expected)
	}
	actual = knownStringOrOther("baz", known)
	expected = "other"
	if actual != expected {
		t.Errorf("knownStringOrOther failed: got %q, expected %q.", actual, expected)
	}
}
func TestLogindCollectorCollectMetrics(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ch := make(chan prometheus.Metric)
	go func() {
		collectMetrics(ch, &testLogindInterface{})
		close(ch)
	}()
	count := 0
	for range ch {
		count++
	}
	expected := len(testSeats) * len(attrRemoteValues) * len(attrTypeValues) * len(attrClassValues)
	if count != expected {
		t.Errorf("collectMetrics did not generate the expected number of metrics: got %d, expected %d.", count, expected)
	}
}
