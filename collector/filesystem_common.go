package collector

import (
	"regexp"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	ignoredMountPoints		= kingpin.Flag("collector.filesystem.ignored-mount-points", "Regexp of mount points to ignore for filesystem collector.").Default(defIgnoredMountPoints).String()
	ignoredFSTypes			= kingpin.Flag("collector.filesystem.ignored-fs-types", "Regexp of filesystem types to ignore for filesystem collector.").Default(defIgnoredFSTypes).String()
	filesystemLabelNames	= []string{"device", "mountpoint", "fstype"}
)

type filesystemCollector struct {
	ignoredMountPointsPattern		*regexp.Regexp
	ignoredFSTypesPattern			*regexp.Regexp
	sizeDesc, freeDesc, availDesc	*prometheus.Desc
	filesDesc, filesFreeDesc		*prometheus.Desc
	roDesc, deviceErrorDesc			*prometheus.Desc
}
type filesystemLabels struct{ device, mountPoint, fsType, options string }
type filesystemStats struct {
	labels				filesystemLabels
	size, free, avail	float64
	files, filesFree	float64
	ro, deviceError		float64
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("filesystem", defaultEnabled, NewFilesystemCollector)
}
func NewFilesystemCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	subsystem := "filesystem"
	mountPointPattern := regexp.MustCompile(*ignoredMountPoints)
	filesystemsTypesPattern := regexp.MustCompile(*ignoredFSTypes)
	sizeDesc := prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "size_bytes"), "Filesystem size in bytes.", filesystemLabelNames, nil)
	freeDesc := prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "free_bytes"), "Filesystem free space in bytes.", filesystemLabelNames, nil)
	availDesc := prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "avail_bytes"), "Filesystem space available to non-root users in bytes.", filesystemLabelNames, nil)
	filesDesc := prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "files"), "Filesystem total file nodes.", filesystemLabelNames, nil)
	filesFreeDesc := prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "files_free"), "Filesystem total free file nodes.", filesystemLabelNames, nil)
	roDesc := prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "readonly"), "Filesystem read-only status.", filesystemLabelNames, nil)
	deviceErrorDesc := prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "device_error"), "Whether an error occurred while getting statistics for the given device.", filesystemLabelNames, nil)
	return &filesystemCollector{ignoredMountPointsPattern: mountPointPattern, ignoredFSTypesPattern: filesystemsTypesPattern, sizeDesc: sizeDesc, freeDesc: freeDesc, availDesc: availDesc, filesDesc: filesDesc, filesFreeDesc: filesFreeDesc, roDesc: roDesc, deviceErrorDesc: deviceErrorDesc}, nil
}
func (c *filesystemCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	stats, err := c.GetStats()
	if err != nil {
		return err
	}
	seen := map[filesystemLabels]bool{}
	for _, s := range stats {
		if seen[s.labels] {
			continue
		}
		seen[s.labels] = true
		ch <- prometheus.MustNewConstMetric(c.deviceErrorDesc, prometheus.GaugeValue, s.deviceError, s.labels.device, s.labels.mountPoint, s.labels.fsType)
		if s.deviceError > 0 {
			continue
		}
		ch <- prometheus.MustNewConstMetric(c.sizeDesc, prometheus.GaugeValue, s.size, s.labels.device, s.labels.mountPoint, s.labels.fsType)
		ch <- prometheus.MustNewConstMetric(c.freeDesc, prometheus.GaugeValue, s.free, s.labels.device, s.labels.mountPoint, s.labels.fsType)
		ch <- prometheus.MustNewConstMetric(c.availDesc, prometheus.GaugeValue, s.avail, s.labels.device, s.labels.mountPoint, s.labels.fsType)
		ch <- prometheus.MustNewConstMetric(c.filesDesc, prometheus.GaugeValue, s.files, s.labels.device, s.labels.mountPoint, s.labels.fsType)
		ch <- prometheus.MustNewConstMetric(c.filesFreeDesc, prometheus.GaugeValue, s.filesFree, s.labels.device, s.labels.mountPoint, s.labels.fsType)
		ch <- prometheus.MustNewConstMetric(c.roDesc, prometheus.GaugeValue, s.ro, s.labels.device, s.labels.mountPoint, s.labels.fsType)
	}
	return nil
}
