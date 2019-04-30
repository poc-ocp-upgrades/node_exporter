package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
)

type tcpConnectionState int

const (
	tcpEstablished	tcpConnectionState	= iota + 1
	tcpSynSent
	tcpSynRecv
	tcpFinWait1
	tcpFinWait2
	tcpTimeWait
	tcpClose
	tcpCloseWait
	tcpLastAck
	tcpListen
	tcpClosing
)

type tcpStatCollector struct{ desc typedDesc }

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("tcpstat", defaultDisabled, NewTCPStatCollector)
}
func NewTCPStatCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &tcpStatCollector{desc: typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, "tcp", "connection_states"), "Number of connection states.", []string{"state"}, nil), prometheus.GaugeValue}}, nil
}
func (c *tcpStatCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	tcpStats, err := getTCPStats(procFilePath("net/tcp"))
	if err != nil {
		return fmt.Errorf("couldn't get tcpstats: %s", err)
	}
	tcp6File := procFilePath("net/tcp6")
	if _, hasIPv6 := os.Stat(tcp6File); hasIPv6 == nil {
		tcp6Stats, err := getTCPStats(tcp6File)
		if err != nil {
			return fmt.Errorf("couldn't get tcp6stats: %s", err)
		}
		for st, value := range tcp6Stats {
			tcpStats[st] += value
		}
	}
	for st, value := range tcpStats {
		ch <- c.desc.mustNewConstMetric(value, st.String())
	}
	return nil
}
func getTCPStats(statsFile string) (map[tcpConnectionState]float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := os.Open(statsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseTCPStats(file)
}
func parseTCPStats(r io.Reader) (map[tcpConnectionState]float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		tcpStats	= map[tcpConnectionState]float64{}
		scanner		= bufio.NewScanner(r)
	)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 0 {
			continue
		}
		if strings.HasPrefix(parts[0], "sl") {
			continue
		}
		st, err := strconv.ParseInt(parts[3], 16, 8)
		if err != nil {
			return nil, err
		}
		tcpStats[tcpConnectionState(st)]++
	}
	return tcpStats, scanner.Err()
}
func (st tcpConnectionState) String() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch st {
	case tcpEstablished:
		return "established"
	case tcpSynSent:
		return "syn_sent"
	case tcpSynRecv:
		return "syn_recv"
	case tcpFinWait1:
		return "fin_wait1"
	case tcpFinWait2:
		return "fin_wait2"
	case tcpTimeWait:
		return "time_wait"
	case tcpClose:
		return "close"
	case tcpCloseWait:
		return "close_wait"
	case tcpLastAck:
		return "last_ack"
	case tcpListen:
		return "listen"
	case tcpClosing:
		return "closing"
	default:
		return "unknown"
	}
}
