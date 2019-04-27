package collector

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	vmStatSubsystem = "vmstat"
)

var (
	vmStatFields = kingpin.Flag("collector.vmstat.fields", "Regexp of fields to return for vmstat collector.").Default("^(oom_kill|pgpg|pswp|pg.*fault).*").String()
)

type vmStatCollector struct{ fieldPattern *regexp.Regexp }

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("vmstat", defaultEnabled, NewvmStatCollector)
}
func NewvmStatCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	pattern := regexp.MustCompile(*vmStatFields)
	return &vmStatCollector{fieldPattern: pattern}, nil
}
func (c *vmStatCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := os.Open(procFilePath("vmstat"))
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		value, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return err
		}
		if !c.fieldPattern.MatchString(parts[0]) {
			continue
		}
		ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, vmStatSubsystem, parts[0]), fmt.Sprintf("/proc/vmstat information field %s.", parts[0]), nil, nil), prometheus.UntypedValue, value)
	}
	return scanner.Err()
}
