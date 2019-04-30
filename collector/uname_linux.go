package collector

import (
	"bytes"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/sys/unix"
)

var unameDesc = prometheus.NewDesc(prometheus.BuildFQName(namespace, "uname", "info"), "Labeled system information as provided by the uname system call.", []string{"sysname", "release", "version", "machine", "nodename", "domainname"}, nil)

type unameCollector struct{}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("uname", defaultEnabled, newUnameCollector)
}
func newUnameCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &unameCollector{}, nil
}
func (c unameCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var uname unix.Utsname
	if err := unix.Uname(&uname); err != nil {
		return err
	}
	ch <- prometheus.MustNewConstMetric(unameDesc, prometheus.GaugeValue, 1, string(uname.Sysname[:bytes.IndexByte(uname.Sysname[:], 0)]), string(uname.Release[:bytes.IndexByte(uname.Release[:], 0)]), string(uname.Version[:bytes.IndexByte(uname.Version[:], 0)]), string(uname.Machine[:bytes.IndexByte(uname.Machine[:], 0)]), string(uname.Nodename[:bytes.IndexByte(uname.Nodename[:], 0)]), string(uname.Domainname[:bytes.IndexByte(uname.Domainname[:], 0)]))
	return nil
}
