package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const (
	KSTAT_DATA_CHAR		= "0"
	KSTAT_DATA_INT32	= "1"
	KSTAT_DATA_UINT32	= "2"
	KSTAT_DATA_INT64	= "3"
	KSTAT_DATA_UINT64	= "4"
	KSTAT_DATA_LONG		= "5"
	KSTAT_DATA_ULONG	= "6"
	KSTAT_DATA_STRING	= "7"
)

func (c *zfsCollector) openProcFile(path string) (*os.File, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := os.Open(procFilePath(path))
	if err != nil {
		log.Debugf("Cannot open %q for reading", procFilePath(path))
		return nil, errZFSNotAvailable
	}
	return file, nil
}
func (c *zfsCollector) updateZfsStats(subsystem string, ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := c.openProcFile(filepath.Join(c.linuxProcpathBase, c.linuxPathMap[subsystem]))
	if err != nil {
		return err
	}
	defer file.Close()
	return c.parseProcfsFile(file, c.linuxPathMap[subsystem], func(s zfsSysctl, v uint64) {
		ch <- c.constSysctlMetric(subsystem, s, v)
	})
}
func (c *zfsCollector) updatePoolStats(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	zpoolPaths, err := filepath.Glob(procFilePath(filepath.Join(c.linuxProcpathBase, c.linuxZpoolIoPath)))
	if err != nil {
		return err
	}
	if zpoolPaths == nil {
		return nil
	}
	for _, zpoolPath := range zpoolPaths {
		file, err := os.Open(zpoolPath)
		if err != nil {
			log.Debugf("Cannot open %q for reading", zpoolPath)
			return errZFSNotAvailable
		}
		err = c.parsePoolProcfsFile(file, zpoolPath, func(poolName string, s zfsSysctl, v uint64) {
			ch <- c.constPoolMetric(poolName, s, v)
		})
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
func (c *zfsCollector) parseProcfsFile(reader io.Reader, fmtExt string, handler func(zfsSysctl, uint64)) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	scanner := bufio.NewScanner(reader)
	parseLine := false
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if !parseLine && len(parts) == 3 && parts[0] == "name" && parts[1] == "type" && parts[2] == "data" {
			parseLine = true
			continue
		}
		if !parseLine || len(parts) < 3 {
			continue
		}
		if parts[1] == KSTAT_DATA_UINT64 {
			key := fmt.Sprintf("kstat.zfs.misc.%s.%s", fmtExt, parts[0])
			value, err := strconv.ParseUint(parts[2], 10, 64)
			if err != nil {
				return fmt.Errorf("could not parse expected integer value for %q", key)
			}
			handler(zfsSysctl(key), value)
		}
	}
	if !parseLine {
		return fmt.Errorf("did not parse a single %q metric", fmtExt)
	}
	return scanner.Err()
}
func (c *zfsCollector) parsePoolProcfsFile(reader io.Reader, zpoolPath string, handler func(string, zfsSysctl, uint64)) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	scanner := bufio.NewScanner(reader)
	parseLine := false
	var fields []string
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		if !parseLine && len(line) >= 12 && line[0] == "nread" {
			parseLine = true
			fields = make([]string, len(line))
			copy(fields, line)
			continue
		}
		if !parseLine {
			continue
		}
		zpoolPathElements := strings.Split(zpoolPath, "/")
		pathLen := len(zpoolPathElements)
		if pathLen < 2 {
			return fmt.Errorf("zpool path did not return at least two elements")
		}
		zpoolName := zpoolPathElements[pathLen-2]
		zpoolFile := zpoolPathElements[pathLen-1]
		for i, field := range fields {
			key := fmt.Sprintf("kstat.zfs.misc.%s.%s", zpoolFile, field)
			value, err := strconv.ParseUint(line[i], 10, 64)
			if err != nil {
				return fmt.Errorf("could not parse expected integer value for %q: %v", key, err)
			}
			handler(zpoolName, zfsSysctl(key), value)
		}
	}
	return scanner.Err()
}
