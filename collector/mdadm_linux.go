package collector

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	statuslineRE			= regexp.MustCompile(`(\d+) blocks .*\[(\d+)/(\d+)\] \[[U_]+\]`)
	raid0lineRE			= regexp.MustCompile(`(\d+) blocks .*\d+k (chunks|rounding)`)
	buildlineRE			= regexp.MustCompile(`\((\d+)/\d+\)`)
	unknownPersonalityLineRE	= regexp.MustCompile(`(\d+) blocks (.*)`)
	raidPersonalityRE		= regexp.MustCompile(`^(linear|raid[0-9]+)$`)
)

type mdStatus struct {
	name		string
	active		bool
	disksActive	int64
	disksTotal	int64
	blocksTotal	int64
	blocksSynced	int64
}
type mdadmCollector struct{}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("mdadm", defaultEnabled, NewMdadmCollector)
}
func evalStatusline(statusline string) (active, total, size int64, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	matches := statuslineRE.FindStringSubmatch(statusline)
	if len(matches) < 3+1 {
		return 0, 0, 0, fmt.Errorf("too few matches found in statusline: %s", statusline)
	} else if len(matches) > 3+1 {
		return 0, 0, 0, fmt.Errorf("too many matches found in statusline: %s", statusline)
	}
	size, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("%s in statusline: %s", err, statusline)
	}
	total, err = strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("%s in statusline: %s", err, statusline)
	}
	active, err = strconv.ParseInt(matches[3], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("%s in statusline: %s", err, statusline)
	}
	return active, total, size, nil
}
func evalRaid0line(statusline string) (size int64, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	matches := raid0lineRE.FindStringSubmatch(statusline)
	if len(matches) < 2 {
		return 0, fmt.Errorf("invalid raid0 status line: %s", statusline)
	}
	size, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s in statusline: %s", err, statusline)
	}
	return size, nil
}
func evalUnknownPersonalitylineRE(statusline string) (size int64, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	matches := unknownPersonalityLineRE.FindStringSubmatch(statusline)
	if len(matches) != 2+1 {
		return 0, fmt.Errorf("invalid unknown personality status line: %s", statusline)
	}
	size, err = strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s in statusline: %s", err, statusline)
	}
	return size, nil
}
func evalBuildline(buildline string) (int64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	matches := buildlineRE.FindStringSubmatch(buildline)
	if len(matches) < 1+1 {
		return 0, fmt.Errorf("too few matches found in buildline: %s", buildline)
	}
	if len(matches) > 1+1 {
		return 0, fmt.Errorf("too many matches found in buildline: %s", buildline)
	}
	syncedSize, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s in buildline: %s", err, buildline)
	}
	return syncedSize, nil
}
func parseMdstat(mdStatusFilePath string) ([]mdStatus, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	content, err := ioutil.ReadFile(mdStatusFilePath)
	if err != nil {
		return []mdStatus{}, err
	}
	lines := strings.Split(string(content), "\n")
	mdStates := make([]mdStatus, 0, len(lines)/3)
	for i, line := range lines {
		if line == "" {
			continue
		}
		if line[0] == ' ' || line[0] == '\t' {
			continue
		}
		if strings.HasPrefix(line, "Personalities") || strings.HasPrefix(line, "unused") {
			continue
		}
		mainLine := strings.Split(line, " ")
		if len(mainLine) < 4 {
			return mdStates, fmt.Errorf("error parsing mdline: %s", line)
		}
		md := mdStatus{name: mainLine[0], active: mainLine[2] == "active"}
		if len(lines) <= i+3 {
			return mdStates, fmt.Errorf("error parsing mdstat: entry for %s has fewer lines than expected", md.name)
		}
		personality := ""
		for _, possiblePersonality := range mainLine[3:] {
			if raidPersonalityRE.MatchString(possiblePersonality) {
				personality = possiblePersonality
				break
			}
		}
		switch {
		case personality == "raid0" || personality == "linear":
			md.disksActive = int64(len(mainLine) - 4)
			md.disksTotal = md.disksActive
			md.blocksTotal, err = evalRaid0line(lines[i+1])
		case raidPersonalityRE.MatchString(personality):
			md.disksActive, md.disksTotal, md.blocksTotal, err = evalStatusline(lines[i+1])
		default:
			log.Debugf("Personality unknown: %s", mainLine)
			md.disksTotal = int64(len(mainLine) - 3)
			md.blocksTotal, err = evalUnknownPersonalitylineRE(lines[i+1])
		}
		if err != nil {
			return mdStates, fmt.Errorf("error parsing mdstat: %s", err)
		}
		if !md.active {
			md.disksActive = 0
		}
		syncLine := lines[i+2]
		if strings.Contains(syncLine, "bitmap") {
			syncLine = lines[i+3]
		}
		if strings.Contains(syncLine, "recovery") || strings.Contains(syncLine, "resync") && !strings.Contains(syncLine, "\tresync=") {
			md.blocksSynced, err = evalBuildline(syncLine)
			if err != nil {
				return mdStates, fmt.Errorf("error parsing mdstat: %s", err)
			}
		} else {
			md.blocksSynced = md.blocksTotal
		}
		mdStates = append(mdStates, md)
	}
	return mdStates, nil
}
func NewMdadmCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &mdadmCollector{}, nil
}

var (
	isActiveDesc		= prometheus.NewDesc(prometheus.BuildFQName(namespace, "md", "is_active"), "Indicator whether the md-device is active or not.", []string{"device"}, nil)
	disksActiveDesc		= prometheus.NewDesc(prometheus.BuildFQName(namespace, "md", "disks_active"), "Number of active disks of device.", []string{"device"}, nil)
	disksTotalDesc		= prometheus.NewDesc(prometheus.BuildFQName(namespace, "md", "disks"), "Total number of disks of device.", []string{"device"}, nil)
	blocksTotalDesc		= prometheus.NewDesc(prometheus.BuildFQName(namespace, "md", "blocks"), "Total number of blocks on device.", []string{"device"}, nil)
	blocksSyncedDesc	= prometheus.NewDesc(prometheus.BuildFQName(namespace, "md", "blocks_synced"), "Number of blocks synced on device.", []string{"device"}, nil)
)

func (c *mdadmCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	statusfile := procFilePath("mdstat")
	mdstate, err := parseMdstat(statusfile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debugf("Not collecting mdstat, file does not exist: %s", statusfile)
			return nil
		}
		return fmt.Errorf("error parsing mdstatus: %s", err)
	}
	for _, mds := range mdstate {
		log.Debugf("collecting metrics for device %s", mds.name)
		var active float64
		if mds.active {
			active = 1
		}
		ch <- prometheus.MustNewConstMetric(isActiveDesc, prometheus.GaugeValue, active, mds.name)
		ch <- prometheus.MustNewConstMetric(disksActiveDesc, prometheus.GaugeValue, float64(mds.disksActive), mds.name)
		ch <- prometheus.MustNewConstMetric(disksTotalDesc, prometheus.GaugeValue, float64(mds.disksTotal), mds.name)
		ch <- prometheus.MustNewConstMetric(blocksTotalDesc, prometheus.GaugeValue, float64(mds.blocksTotal), mds.name)
		ch <- prometheus.MustNewConstMetric(blocksSyncedDesc, prometheus.GaugeValue, float64(mds.blocksSynced), mds.name)
	}
	return nil
}
