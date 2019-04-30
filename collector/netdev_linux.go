package collector

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"github.com/prometheus/common/log"
)

var (
	procNetDevInterfaceRE	= regexp.MustCompile(`^(.+): *(.+)$`)
	procNetDevFieldSep	= regexp.MustCompile(` +`)
)

func getNetDevStats(ignore *regexp.Regexp) (map[string]map[string]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := os.Open(procFilePath("net/dev"))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return parseNetDevStats(file, ignore)
}
func parseNetDevStats(r io.Reader, ignore *regexp.Regexp) (map[string]map[string]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	scanner.Scan()
	parts := strings.Split(scanner.Text(), "|")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid header line in net/dev: %s", scanner.Text())
	}
	receiveHeader := strings.Fields(parts[1])
	transmitHeader := strings.Fields(parts[2])
	headerLength := len(receiveHeader) + len(transmitHeader)
	netDev := map[string]map[string]string{}
	for scanner.Scan() {
		line := strings.TrimLeft(scanner.Text(), " ")
		parts := procNetDevInterfaceRE.FindStringSubmatch(line)
		if len(parts) != 3 {
			return nil, fmt.Errorf("couldn't get interface name, invalid line in net/dev: %q", line)
		}
		dev := parts[1]
		if ignore.MatchString(dev) {
			log.Debugf("Ignoring device: %s", dev)
			continue
		}
		values := procNetDevFieldSep.Split(strings.TrimLeft(parts[2], " "), -1)
		if len(values) != headerLength {
			return nil, fmt.Errorf("couldn't get values, invalid line in net/dev: %q", parts[2])
		}
		netDev[dev] = map[string]string{}
		for i := 0; i < len(receiveHeader); i++ {
			netDev[dev]["receive_"+receiveHeader[i]] = values[i]
		}
		for i := 0; i < len(transmitHeader); i++ {
			netDev[dev]["transmit_"+transmitHeader[i]] = values[i+len(receiveHeader)]
		}
	}
	return netDev, scanner.Err()
}
