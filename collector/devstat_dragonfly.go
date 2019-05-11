package collector

import (
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
)
import "C"

const (
	devstatSubsystem = "devstat"
)

type devstatCollector struct {
	bytesDesc		*prometheus.Desc
	transfersDesc	*prometheus.Desc
	blocksDesc		*prometheus.Desc
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("devstat", defaultDisabled, NewDevstatCollector)
}
func NewDevstatCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &devstatCollector{bytesDesc: prometheus.NewDesc(prometheus.BuildFQName(namespace, devstatSubsystem, "bytes_total"), "The total number of bytes transferred for reads and writes on the device.", []string{"device"}, nil), transfersDesc: prometheus.NewDesc(prometheus.BuildFQName(namespace, devstatSubsystem, "transfers_total"), "The total number of transactions completed.", []string{"device"}, nil), blocksDesc: prometheus.NewDesc(prometheus.BuildFQName(namespace, devstatSubsystem, "blocks_total"), "The total number of bytes given in terms of the devices blocksize.", []string{"device"}, nil)}, nil
}
func (c *devstatCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	count := C._get_ndevs()
	if count == -1 {
		return errors.New("getdevs() failed")
	}
	if count == -2 {
		return errors.New("calloc() failed")
	}
	for i := C.int(0); i < count; i++ {
		stats := C._get_stats(i)
		device := fmt.Sprintf("%s%d", C.GoString(&stats.device[0]), stats.unit)
		ch <- prometheus.MustNewConstMetric(c.bytesDesc, prometheus.CounterValue, float64(stats.bytes), device)
		ch <- prometheus.MustNewConstMetric(c.transfersDesc, prometheus.CounterValue, float64(stats.transfers), device)
		ch <- prometheus.MustNewConstMetric(c.blocksDesc, prometheus.CounterValue, float64(stats.blocks), device)
	}
	return nil
}
