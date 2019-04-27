package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	sockStatSubsystem = "sockstat"
)

var pageSize = os.Getpagesize()

type sockStatCollector struct{}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector(sockStatSubsystem, defaultEnabled, NewSockStatCollector)
}
func NewSockStatCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &sockStatCollector{}, nil
}
func (c *sockStatCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sockStats, err := getSockStats(procFilePath("net/sockstat"))
	if err != nil {
		return fmt.Errorf("couldn't get sockstats: %s", err)
	}
	for protocol, protocolStats := range sockStats {
		for name, value := range protocolStats {
			v, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid value %s in sockstats: %s", value, err)
			}
			ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, sockStatSubsystem, protocol+"_"+name), fmt.Sprintf("Number of %s sockets in state %s.", protocol, name), nil, nil), prometheus.GaugeValue, v)
		}
	}
	return err
}
func getSockStats(fileName string) (map[string]map[string]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseSockStats(file, fileName)
}
func parseSockStats(r io.Reader, fileName string) (map[string]map[string]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		sockStat	= map[string]map[string]string{}
		scanner		= bufio.NewScanner(r)
	)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " ")
		protocol := line[0][:len(line[0])-1]
		sockStat[protocol] = map[string]string{}
		for i := 1; i < len(line) && i+1 < len(line); i++ {
			sockStat[protocol][line[i]] = line[i+1]
			i++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	pageCount, err := strconv.Atoi(sockStat["TCP"]["mem"])
	if err != nil {
		return nil, fmt.Errorf("invalid value %s in sockstats: %s", sockStat["TCP"]["mem"], err)
	}
	sockStat["TCP"]["mem_bytes"] = strconv.Itoa(pageCount * pageSize)
	if udpMem := sockStat["UDP"]["mem"]; udpMem != "" {
		pageCount, err = strconv.Atoi(udpMem)
		if err != nil {
			return nil, fmt.Errorf("invalid value %s in sockstats: %s", sockStat["UDP"]["mem"], err)
		}
		sockStat["UDP"]["mem_bytes"] = strconv.Itoa(pageCount * pageSize)
	}
	return sockStat, nil
}
