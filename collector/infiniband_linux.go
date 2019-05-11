package collector

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

const infinibandPath = "class/infiniband"

var (
	errInfinibandNoDevicesFound	= errors.New("no InfiniBand devices detected")
	errInfinibandNoPortsFound	= errors.New("no InfiniBand ports detected")
)

type infinibandCollector struct {
	metricDescs		map[string]*prometheus.Desc
	counters		map[string]infinibandMetric
	legacyCounters	map[string]infinibandMetric
}
type infinibandMetric struct {
	File	string
	Help	string
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("infiniband", defaultEnabled, NewInfiniBandCollector)
}
func NewInfiniBandCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var i infinibandCollector
	i.counters = map[string]infinibandMetric{"link_downed_total": {"link_downed", "Number of times the link failed to recover from an error state and went down"}, "link_error_recovery_total": {"link_error_recovery", "Number of times the link successfully recovered from an error state"}, "multicast_packets_received_total": {"multicast_rcv_packets", "Number of multicast packets received (including errors)"}, "multicast_packets_transmitted_total": {"multicast_xmit_packets", "Number of multicast packets transmitted (including errors)"}, "port_data_received_bytes_total": {"port_rcv_data", "Number of data octets received on all links"}, "port_data_transmitted_bytes_total": {"port_xmit_data", "Number of data octets transmitted on all links"}, "unicast_packets_received_total": {"unicast_rcv_packets", "Number of unicast packets received (including errors)"}, "unicast_packets_transmitted_total": {"unicast_xmit_packets", "Number of unicast packets transmitted (including errors)"}}
	i.legacyCounters = map[string]infinibandMetric{"legacy_multicast_packets_received_total": {"port_multicast_rcv_packets", "Number of multicast packets received"}, "legacy_multicast_packets_transmitted_total": {"port_multicast_xmit_packets", "Number of multicast packets transmitted"}, "legacy_data_received_bytes_total": {"port_rcv_data_64", "Number of data octets received on all links"}, "legacy_packets_received_total": {"port_rcv_packets_64", "Number of data packets received on all links"}, "legacy_unicast_packets_received_total": {"port_unicast_rcv_packets", "Number of unicast packets received"}, "legacy_unicast_packets_transmitted_total": {"port_unicast_xmit_packets", "Number of unicast packets transmitted"}, "legacy_data_transmitted_bytes_total": {"port_xmit_data_64", "Number of data octets transmitted on all links"}, "legacy_packets_transmitted_total": {"port_xmit_packets_64", "Number of data packets received on all links"}}
	subsystem := "infiniband"
	i.metricDescs = make(map[string]*prometheus.Desc)
	for metricName, infinibandMetric := range i.counters {
		i.metricDescs[metricName] = prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, metricName), infinibandMetric.Help, []string{"device", "port"}, nil)
	}
	for metricName, infinibandMetric := range i.legacyCounters {
		i.metricDescs[metricName] = prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, metricName), infinibandMetric.Help, []string{"device", "port"}, nil)
	}
	return &i, nil
}
func infinibandDevices(infinibandPath string) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	devices, err := filepath.Glob(filepath.Join(infinibandPath, "/*"))
	if err != nil {
		return nil, err
	}
	if len(devices) < 1 {
		log.Debugf("Unable to detect InfiniBand devices")
		err = errInfinibandNoDevicesFound
		return nil, err
	}
	for i, device := range devices {
		devices[i] = filepath.Base(device)
	}
	return devices, nil
}
func infinibandPorts(infinibandPath, device string) ([]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ports, err := filepath.Glob(filepath.Join(infinibandPath, device, "ports/*"))
	if err != nil {
		return nil, err
	}
	if len(ports) < 1 {
		log.Debugf("Unable to detect ports for %s", device)
		err = errInfinibandNoPortsFound
		return nil, err
	}
	for i, port := range ports {
		ports[i] = filepath.Base(port)
	}
	return ports, nil
}
func readMetric(directory, metricFile string) (uint64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	metric, err := readUintFromFile(filepath.Join(directory, metricFile))
	if err != nil {
		if strings.Contains(err.Error(), "N/A (no PMA)") {
			log.Debugf("%q value is N/A", metricFile)
			return 0, nil
		}
		log.Debugf("Error reading %q file", metricFile)
		return 0, err
	}
	switch metricFile {
	case "port_rcv_data", "port_xmit_data", "port_rcv_data_64", "port_xmit_data_64":
		metric *= 4
	}
	return metric, nil
}
func (c *infinibandCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	devices, err := infinibandDevices(sysFilePath(infinibandPath))
	switch err {
	case nil:
	case errInfinibandNoDevicesFound:
		return nil
	default:
		return err
	}
	for _, device := range devices {
		ports, err := infinibandPorts(sysFilePath(infinibandPath), device)
		switch err {
		case nil:
		case errInfinibandNoPortsFound:
			continue
		default:
			return err
		}
		for _, port := range ports {
			portFiles := sysFilePath(filepath.Join(infinibandPath, device, "ports", port))
			for metricName, infinibandMetric := range c.counters {
				if _, err := os.Stat(filepath.Join(portFiles, "counters", infinibandMetric.File)); os.IsNotExist(err) {
					continue
				}
				metric, err := readMetric(filepath.Join(portFiles, "counters"), infinibandMetric.File)
				if err != nil {
					return err
				}
				ch <- prometheus.MustNewConstMetric(c.metricDescs[metricName], prometheus.CounterValue, float64(metric), device, port)
			}
			for metricName, infinibandMetric := range c.legacyCounters {
				if _, err := os.Stat(filepath.Join(portFiles, "counters_ext", infinibandMetric.File)); os.IsNotExist(err) {
					continue
				}
				metric, err := readMetric(filepath.Join(portFiles, "counters_ext"), infinibandMetric.File)
				if err != nil {
					return err
				}
				ch <- prometheus.MustNewConstMetric(c.metricDescs[metricName], prometheus.CounterValue, float64(metric), device, port)
			}
		}
	}
	return nil
}
