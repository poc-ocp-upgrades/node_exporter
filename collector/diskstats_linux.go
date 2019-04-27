package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	diskSubsystem		= "disk"
	diskSectorSize		= 512
	diskstatsFilename	= "diskstats"
)

var (
	ignoredDevices = kingpin.Flag("collector.diskstats.ignored-devices", "Regexp of devices to ignore for diskstats.").Default("^(ram|loop|fd|(h|s|v|xv)d[a-z]|nvme\\d+n\\d+p)\\d+$").String()
)

type typedFactorDesc struct {
	desc		*prometheus.Desc
	valueType	prometheus.ValueType
	factor		float64
}

func (d *typedFactorDesc) mustNewConstMetric(value float64, labels ...string) prometheus.Metric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if d.factor != 0 {
		value *= d.factor
	}
	return prometheus.MustNewConstMetric(d.desc, d.valueType, value, labels...)
}

type diskstatsCollector struct {
	ignoredDevicesPattern	*regexp.Regexp
	descs			[]typedFactorDesc
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("diskstats", defaultEnabled, NewDiskstatsCollector)
}
func NewDiskstatsCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var diskLabelNames = []string{"device"}
	return &diskstatsCollector{ignoredDevicesPattern: regexp.MustCompile(*ignoredDevices), descs: []typedFactorDesc{{desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "reads_completed_total"), "The total number of reads completed successfully.", diskLabelNames, nil), valueType: prometheus.CounterValue}, {desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "reads_merged_total"), "The total number of reads merged.", diskLabelNames, nil), valueType: prometheus.CounterValue}, {desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "read_bytes_total"), "The total number of bytes read successfully.", diskLabelNames, nil), valueType: prometheus.CounterValue, factor: diskSectorSize}, {desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "read_time_seconds_total"), "The total number of seconds spent by all reads.", diskLabelNames, nil), valueType: prometheus.CounterValue, factor: .001}, {desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "writes_completed_total"), "The total number of writes completed successfully.", diskLabelNames, nil), valueType: prometheus.CounterValue}, {desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "writes_merged_total"), "The number of writes merged.", diskLabelNames, nil), valueType: prometheus.CounterValue}, {desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "written_bytes_total"), "The total number of bytes written successfully.", diskLabelNames, nil), valueType: prometheus.CounterValue, factor: diskSectorSize}, {desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "write_time_seconds_total"), "This is the total number of seconds spent by all writes.", diskLabelNames, nil), valueType: prometheus.CounterValue, factor: .001}, {desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "io_now"), "The number of I/Os currently in progress.", diskLabelNames, nil), valueType: prometheus.GaugeValue}, {desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "io_time_seconds_total"), "Total seconds spent doing I/Os.", diskLabelNames, nil), valueType: prometheus.CounterValue, factor: .001}, {desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "io_time_weighted_seconds_total"), "The weighted # of seconds spent doing I/Os.", diskLabelNames, nil), valueType: prometheus.CounterValue, factor: .001}, {desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "discards_completed_total"), "The total number of discards completed successfully.", diskLabelNames, nil), valueType: prometheus.CounterValue}, {desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "discards_merged_total"), "The total number of discards merged.", diskLabelNames, nil), valueType: prometheus.CounterValue}, {desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "discarded_sectors_total"), "The total number of sectors discarded successfully.", diskLabelNames, nil), valueType: prometheus.CounterValue}, {desc: prometheus.NewDesc(prometheus.BuildFQName(namespace, diskSubsystem, "discard_time_seconds_total"), "This is the total number of seconds spent by all discards.", diskLabelNames, nil), valueType: prometheus.CounterValue, factor: .001}}}, nil
}
func (c *diskstatsCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	diskStats, err := getDiskStats()
	if err != nil {
		return fmt.Errorf("couldn't get diskstats: %s", err)
	}
	for dev, stats := range diskStats {
		if c.ignoredDevicesPattern.MatchString(dev) {
			log.Debugf("Ignoring device: %s", dev)
			continue
		}
		for i, value := range stats {
			if i >= len(c.descs) {
				break
			}
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid value %s in diskstats: %s", value, err)
			}
			ch <- c.descs[i].mustNewConstMetric(v, dev)
		}
	}
	return nil
}
func getDiskStats() (map[string][]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := os.Open(procFilePath(diskstatsFilename))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseDiskStats(file)
}
func parseDiskStats(r io.Reader) (map[string][]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		diskStats	= map[string][]string{}
		scanner		= bufio.NewScanner(r)
	)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) < 4 {
			return nil, fmt.Errorf("invalid line in %s: %s", procFilePath(diskstatsFilename), scanner.Text())
		}
		dev := parts[2]
		diskStats[dev] = parts[3:]
	}
	return diskStats, scanner.Err()
}
