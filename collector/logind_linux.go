package collector

import (
	"fmt"
	"os"
	"strconv"
	"github.com/godbus/dbus"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	logindSubsystem	= "logind"
	dbusObject	= "org.freedesktop.login1"
	dbusPath	= "/org/freedesktop/login1"
)

var (
	attrRemoteValues	= []string{"true", "false"}
	attrTypeValues		= []string{"other", "unspecified", "tty", "x11", "wayland", "mir", "web"}
	attrClassValues		= []string{"other", "user", "greeter", "lock-screen", "background"}
	sessionsDesc		= prometheus.NewDesc(prometheus.BuildFQName(namespace, logindSubsystem, "sessions"), "Number of sessions registered in logind.", []string{"seat", "remote", "type", "class"}, nil)
)

type logindCollector struct{}
type logindDbus struct {
	conn	*dbus.Conn
	object	dbus.BusObject
}
type logindInterface interface {
	listSeats() ([]string, error)
	listSessions() ([]logindSessionEntry, error)
	getSession(logindSessionEntry) *logindSession
}
type logindSession struct {
	seat		string
	remote		string
	sessionType	string
	class		string
}
type logindSessionEntry struct {
	SessionID		string
	UserID			uint32
	UserName		string
	SeatID			string
	SessionObjectPath	dbus.ObjectPath
}
type logindSeatEntry struct {
	SeatID		string
	SeatObjectPath	dbus.ObjectPath
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("logind", defaultDisabled, NewLogindCollector)
}
func NewLogindCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &logindCollector{}, nil
}
func (lc *logindCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c, err := newDbus()
	if err != nil {
		return fmt.Errorf("unable to connect to dbus: %s", err)
	}
	defer c.conn.Close()
	return collectMetrics(ch, c)
}
func collectMetrics(ch chan<- prometheus.Metric, c logindInterface) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	seats, err := c.listSeats()
	if err != nil {
		return fmt.Errorf("unable to get seats: %s", err)
	}
	sessionList, err := c.listSessions()
	if err != nil {
		return fmt.Errorf("unable to get sessions: %s", err)
	}
	sessions := make(map[logindSession]float64)
	for _, s := range sessionList {
		session := c.getSession(s)
		if session != nil {
			sessions[*session]++
		}
	}
	for _, remote := range attrRemoteValues {
		for _, sessionType := range attrTypeValues {
			for _, class := range attrClassValues {
				for _, seat := range seats {
					count := sessions[logindSession{seat, remote, sessionType, class}]
					ch <- prometheus.MustNewConstMetric(sessionsDesc, prometheus.GaugeValue, count, seat, remote, sessionType, class)
				}
			}
		}
	}
	return nil
}
func knownStringOrOther(value string, known []string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for i := range known {
		if value == known[i] {
			return value
		}
	}
	return "other"
}
func newDbus() (*logindDbus, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	conn, err := dbus.SystemBusPrivate()
	if err != nil {
		return nil, err
	}
	methods := []dbus.Auth{dbus.AuthExternal(strconv.Itoa(os.Getuid()))}
	err = conn.Auth(methods)
	if err != nil {
		conn.Close()
		return nil, err
	}
	err = conn.Hello()
	if err != nil {
		conn.Close()
		return nil, err
	}
	object := conn.Object(dbusObject, dbus.ObjectPath(dbusPath))
	return &logindDbus{conn: conn, object: object}, nil
}
func (c *logindDbus) listSeats() ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var result [][]interface{}
	err := c.object.Call(dbusObject+".Manager.ListSeats", 0).Store(&result)
	if err != nil {
		return nil, err
	}
	resultInterface := make([]interface{}, len(result))
	for i := range result {
		resultInterface[i] = result[i]
	}
	seats := make([]logindSeatEntry, len(result))
	seatsInterface := make([]interface{}, len(seats))
	for i := range seats {
		seatsInterface[i] = &seats[i]
	}
	err = dbus.Store(resultInterface, seatsInterface...)
	if err != nil {
		return nil, err
	}
	ret := make([]string, len(seats)+1)
	for i := range seats {
		ret[i] = seats[i].SeatID
	}
	ret[len(seats)] = ""
	return ret, nil
}
func (c *logindDbus) listSessions() ([]logindSessionEntry, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var result [][]interface{}
	err := c.object.Call(dbusObject+".Manager.ListSessions", 0).Store(&result)
	if err != nil {
		return nil, err
	}
	resultInterface := make([]interface{}, len(result))
	for i := range result {
		resultInterface[i] = result[i]
	}
	sessions := make([]logindSessionEntry, len(result))
	sessionsInterface := make([]interface{}, len(sessions))
	for i := range sessions {
		sessionsInterface[i] = &sessions[i]
	}
	err = dbus.Store(resultInterface, sessionsInterface...)
	if err != nil {
		return nil, err
	}
	return sessions, nil
}
func (c *logindDbus) getSession(session logindSessionEntry) *logindSession {
	_logClusterCodePath()
	defer _logClusterCodePath()
	object := c.conn.Object(dbusObject, session.SessionObjectPath)
	remote, err := object.GetProperty(dbusObject + ".Session.Remote")
	if err != nil {
		return nil
	}
	sessionType, err := object.GetProperty(dbusObject + ".Session.Type")
	if err != nil {
		return nil
	}
	sessionTypeStr, ok := sessionType.Value().(string)
	if !ok {
		return nil
	}
	class, err := object.GetProperty(dbusObject + ".Session.Class")
	if err != nil {
		return nil
	}
	classStr, ok := class.Value().(string)
	if !ok {
		return nil
	}
	return &logindSession{seat: session.SeatID, remote: remote.String(), sessionType: knownStringOrOther(sessionTypeStr, attrTypeValues), class: knownStringOrOther(classStr, attrClassValues)}
}
