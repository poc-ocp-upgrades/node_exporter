package collector

import (
	"fmt"
	"github.com/lufia/iostat"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	diskSubsystem = "disk"
)

type typedDescFunc struct {
	typedDesc
	value	func(stat *iostat.DriveStats) float64
}
type diskstatsCollector struct{ descs []typedDescFunc }

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("diskstats", defaultEnabled, NewDiskstatsCollector)
}
func NewDiskstatsCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var diskLabelNames = []string{"device"}
	return &diskstatsCollector{descs: []typedDescFunc{{typedDesc: typedDesc{desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "reads_completed_total"), "The total number of reads completed successfully.", diskLabelNames, nil), valueType: prometheus.CounterValue}, value: func(stat *iostat.DriveStats) float64 {
		return float64(stat.NumRead)
	}}, {typedDesc: typedDesc{desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "read_sectors_total"), "The total number of sectors read successfully.", diskLabelNames, nil), valueType: prometheus.CounterValue}, value: func(stat *iostat.DriveStats) float64 {
		return float64(stat.NumRead) / float64(stat.BlockSize)
	}}, {typedDesc: typedDesc{desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "read_time_seconds_total"), "The total number of seconds spent by all reads.", diskLabelNames, nil), valueType: prometheus.CounterValue}, value: func(stat *iostat.DriveStats) float64 {
		return stat.TotalReadTime.Seconds()
	}}, {typedDesc: typedDesc{desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "writes_completed_total"), "The total number of writes completed successfully.", diskLabelNames, nil), valueType: prometheus.CounterValue}, value: func(stat *iostat.DriveStats) float64 {
		return float64(stat.NumWrite)
	}}, {typedDesc: typedDesc{desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "written_sectors_total"), "The total number of sectors written successfully.", diskLabelNames, nil), valueType: prometheus.CounterValue}, value: func(stat *iostat.DriveStats) float64 {
		return float64(stat.NumWrite) / float64(stat.BlockSize)
	}}, {typedDesc: typedDesc{desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "write_time_seconds_total"), "This is the total number of seconds spent by all writes.", diskLabelNames, nil), valueType: prometheus.CounterValue}, value: func(stat *iostat.DriveStats) float64 {
		return stat.TotalWriteTime.Seconds()
	}}, {typedDesc: typedDesc{desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "read_bytes_total"), "The total number of bytes read successfully.", diskLabelNames, nil), valueType: prometheus.CounterValue}, value: func(stat *iostat.DriveStats) float64 {
		return float64(stat.BytesRead)
	}}, {typedDesc: typedDesc{desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "written_bytes_total"), "The total number of bytes written successfully.", diskLabelNames, nil), valueType: prometheus.CounterValue}, value: func(stat *iostat.DriveStats) float64 {
		return float64(stat.BytesWritten)
	}}}}, nil
}
func (c *diskstatsCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	diskStats, err := iostat.ReadDriveStats()
	if err != nil {
		return fmt.Errorf("couldn't get diskstats: %s", err)
	}
	for _, stats := range diskStats {
		for _, desc := range c.descs {
			v := desc.value(stats)
			ch <- desc.mustNewConstMetric(v, stats.Name)
		}
	}
	return nil
}
