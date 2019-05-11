package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func (c *meminfoCollector) getMemInfo() (map[string]float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := os.Open(procFilePath("meminfo"))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseMemInfo(file)
}
func parseMemInfo(r io.Reader) (map[string]float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		memInfo	= map[string]float64{}
		scanner	= bufio.NewScanner(r)
		re		= regexp.MustCompile(`\((.*)\)`)
	)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		fv, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid value in meminfo: %s", err)
		}
		key := parts[0][:len(parts[0])-1]
		key = re.ReplaceAllString(key, "_${1}")
		switch len(parts) {
		case 2:
		case 3:
			fv *= 1024
			key = key + "_bytes"
		default:
			return nil, fmt.Errorf("invalid line in meminfo: %s", line)
		}
		memInfo[key] = fv
	}
	return memInfo, scanner.Err()
}
