package collector

import (
	"errors"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var errZFSNotAvailable = errors.New("ZFS / ZFS statistics are not available")

type zfsSysctl string

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("zfs", defaultEnabled, NewZFSCollector)
}

type zfsCollector struct {
	linuxProcpathBase	string
	linuxZpoolIoPath	string
	linuxPathMap		map[string]string
}

func NewZFSCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &zfsCollector{linuxProcpathBase: "spl/kstat/zfs", linuxZpoolIoPath: "/*/io", linuxPathMap: map[string]string{"zfs_abd": "abdstats", "zfs_arc": "arcstats", "zfs_dbuf": "dbuf_stats", "zfs_dmu_tx": "dmu_tx", "zfs_dnode": "dnodestats", "zfs_fm": "fm", "zfs_vdev_cache": "vdev_cache_stats", "zfs_vdev_mirror": "vdev_mirror_stats", "zfs_xuio": "xuio_stats", "zfs_zfetch": "zfetchstats", "zfs_zil": "zil"}}, nil
}
func (c *zfsCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for subsystem := range c.linuxPathMap {
		if err := c.updateZfsStats(subsystem, ch); err != nil {
			if err == errZFSNotAvailable {
				log.Debug(err)
				continue
			}
			return err
		}
	}
	return c.updatePoolStats(ch)
}
func (s zfsSysctl) metricName() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	parts := strings.Split(string(s), ".")
	return strings.Replace(parts[len(parts)-1], "-", "_", -1)
}
func (c *zfsCollector) constSysctlMetric(subsystem string, sysctl zfsSysctl, value uint64) prometheus.Metric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	metricName := sysctl.metricName()
	return prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, metricName), string(sysctl), nil, nil), prometheus.UntypedValue, float64(value))
}
func (c *zfsCollector) constPoolMetric(poolName string, sysctl zfsSysctl, value uint64) prometheus.Metric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	metricName := sysctl.metricName()
	return prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, "zfs_zpool", metricName), string(sysctl), []string{"zpool"}, nil), prometheus.UntypedValue, float64(value), poolName)
}
