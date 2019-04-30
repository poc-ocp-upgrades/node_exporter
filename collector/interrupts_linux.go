package collector

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	interruptLabelNames = []string{"cpu", "type", "info", "devices"}
)

func (c *interruptsCollector) Update(ch chan<- prometheus.Metric) (err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	interrupts, err := getInterrupts()
	if err != nil {
		return fmt.Errorf("couldn't get interrupts: %s", err)
	}
	for name, interrupt := range interrupts {
		for cpuNo, value := range interrupt.values {
			fv, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("invalid value %s in interrupts: %s", value, err)
			}
			ch <- c.desc.mustNewConstMetric(fv, strconv.Itoa(cpuNo), name, interrupt.info, interrupt.devices)
		}
	}
	return err
}

type interrupt struct {
	info	string
	devices	string
	values	[]string
}

func getInterrupts() (map[string]interrupt, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := os.Open(procFilePath("interrupts"))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseInterrupts(file)
}
func parseInterrupts(r io.Reader) (map[string]interrupt, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var (
		interrupts	= map[string]interrupt{}
		scanner		= bufio.NewScanner(r)
	)
	if !scanner.Scan() {
		return nil, errors.New("interrupts empty")
	}
	cpuNum := len(strings.Fields(scanner.Text()))
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) < cpuNum+2 {
			continue
		}
		intName := parts[0][:len(parts[0])-1]
		intr := interrupt{values: parts[1 : cpuNum+1]}
		if _, err := strconv.Atoi(intName); err == nil {
			intr.info = parts[cpuNum+1]
			intr.devices = strings.Join(parts[cpuNum+2:], " ")
		} else {
			intr.info = strings.Join(parts[cpuNum+1:], " ")
		}
		interrupts[intName] = intr
	}
	return interrupts, scanner.Err()
}
