package collector

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"github.com/ema/qdisc"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
)

type qdiscStatCollector struct {
	bytes		typedDesc
	packets		typedDesc
	drops		typedDesc
	requeues	typedDesc
	overlimits	typedDesc
}

var (
	collectorQdisc = kingpin.Flag("collector.qdisc.fixtures", "test fixtures to use for qdisc collector end-to-end testing").Default("").String()
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("qdisc", defaultDisabled, NewQdiscStatCollector)
}
func NewQdiscStatCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &qdiscStatCollector{bytes: typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, "qdisc", "bytes_total"), "Number of bytes sent.", []string{"device", "kind"}, nil), prometheus.CounterValue}, packets: typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, "qdisc", "packets_total"), "Number of packets sent.", []string{"device", "kind"}, nil), prometheus.CounterValue}, drops: typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, "qdisc", "drops_total"), "Number of packets dropped.", []string{"device", "kind"}, nil), prometheus.CounterValue}, requeues: typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, "qdisc", "requeues_total"), "Number of packets dequeued, not transmitted, and requeued.", []string{"device", "kind"}, nil), prometheus.CounterValue}, overlimits: typedDesc{prometheus.NewDesc(prometheus.BuildFQName(namespace, "qdisc", "overlimits_total"), "Number of overlimit packets.", []string{"device", "kind"}, nil), prometheus.CounterValue}}, nil
}
func testQdiscGet(fixtures string) ([]qdisc.QdiscInfo, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var res []qdisc.QdiscInfo
	b, err := ioutil.ReadFile(filepath.Join(fixtures, "results.json"))
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(b, &res)
	return res, err
}
func (c *qdiscStatCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var msgs []qdisc.QdiscInfo
	var err error
	fixtures := *collectorQdisc
	if fixtures == "" {
		msgs, err = qdisc.Get()
	} else {
		msgs, err = testQdiscGet(fixtures)
	}
	if err != nil {
		return err
	}
	for _, msg := range msgs {
		if msg.Parent != 0 {
			continue
		}
		ch <- c.bytes.mustNewConstMetric(float64(msg.Bytes), msg.IfaceName, msg.Kind)
		ch <- c.packets.mustNewConstMetric(float64(msg.Packets), msg.IfaceName, msg.Kind)
		ch <- c.drops.mustNewConstMetric(float64(msg.Drops), msg.IfaceName, msg.Kind)
		ch <- c.requeues.mustNewConstMetric(float64(msg.Requeues), msg.IfaceName, msg.Kind)
		ch <- c.overlimits.mustNewConstMetric(float64(msg.Overlimits), msg.IfaceName, msg.Kind)
	}
	return nil
}
