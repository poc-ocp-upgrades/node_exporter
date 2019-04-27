package collector

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	fileFDStatSubsystem = "filefd"
)

type fileFDStatCollector struct{}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector(fileFDStatSubsystem, defaultEnabled, NewFileFDStatCollector)
}
func NewFileFDStatCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &fileFDStatCollector{}, nil
}
func (c *fileFDStatCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	fileFDStat, err := parseFileFDStats(procFilePath("sys/fs/file-nr"))
	if err != nil {
		return fmt.Errorf("couldn't get file-nr: %s", err)
	}
	for name, value := range fileFDStat {
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("invalid value %s in file-nr: %s", value, err)
		}
		ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, fileFDStatSubsystem, name), fmt.Sprintf("File descriptor statistics: %s.", name), nil, nil), prometheus.GaugeValue, v)
	}
	return nil
}
func parseFileFDStats(filename string) (map[string]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	parts := bytes.Split(bytes.TrimSpace(content), []byte("\u0009"))
	if len(parts) < 3 {
		return nil, fmt.Errorf("unexpected number of file stats in %q", filename)
	}
	var fileFDStat = map[string]string{}
	fileFDStat["allocated"] = string(parts[0])
	fileFDStat["maximum"] = string(parts[2])
	return fileFDStat, nil
}
