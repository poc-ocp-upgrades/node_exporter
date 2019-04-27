package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"github.com/mdlayher/wifi"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

type wifiCollector struct {
	interfaceFrequencyHertz		*prometheus.Desc
	stationInfo			*prometheus.Desc
	stationConnectedSecondsTotal	*prometheus.Desc
	stationInactiveSeconds		*prometheus.Desc
	stationReceiveBitsPerSecond	*prometheus.Desc
	stationTransmitBitsPerSecond	*prometheus.Desc
	stationSignalDBM		*prometheus.Desc
	stationTransmitRetriesTotal	*prometheus.Desc
	stationTransmitFailedTotal	*prometheus.Desc
	stationBeaconLossTotal		*prometheus.Desc
}

var (
	collectorWifi = kingpin.Flag("collector.wifi.fixtures", "test fixtures to use for wifi collector metrics").Default("").String()
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("wifi", defaultDisabled, NewWifiCollector)
}

var _ wifiStater = &wifi.Client{}

type wifiStater interface {
	BSS(ifi *wifi.Interface) (*wifi.BSS, error)
	Close() error
	Interfaces() ([]*wifi.Interface, error)
	StationInfo(ifi *wifi.Interface) ([]*wifi.StationInfo, error)
}

func NewWifiCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	const (
		subsystem = "wifi"
	)
	var (
		labels = []string{"device", "mac_address"}
	)
	return &wifiCollector{interfaceFrequencyHertz: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "interface_frequency_hertz"), "The current frequency a WiFi interface is operating at, in hertz.", []string{"device"}, nil), stationInfo: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "station_info"), "Labeled WiFi interface station information as provided by the operating system.", []string{"device", "bssid", "ssid", "mode"}, nil), stationConnectedSecondsTotal: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "station_connected_seconds_total"), "The total number of seconds a station has been connected to an access point.", labels, nil), stationInactiveSeconds: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "station_inactive_seconds"), "The number of seconds since any wireless activity has occurred on a station.", labels, nil), stationReceiveBitsPerSecond: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "station_receive_bits_per_second"), "The current WiFi receive bitrate of a station, in bits per second.", labels, nil), stationTransmitBitsPerSecond: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "station_transmit_bits_per_second"), "The current WiFi transmit bitrate of a station, in bits per second.", labels, nil), stationSignalDBM: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "station_signal_dbm"), "The current WiFi signal strength, in decibel-milliwatts (dBm).", labels, nil), stationTransmitRetriesTotal: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "station_transmit_retries_total"), "The total number of times a station has had to retry while sending a packet.", labels, nil), stationTransmitFailedTotal: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "station_transmit_failed_total"), "The total number of times a station has failed to send a packet.", labels, nil), stationBeaconLossTotal: prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, "station_beacon_loss_total"), "The total number of times a station has detected a beacon loss.", labels, nil)}, nil
}
func (c *wifiCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	stat, err := newWifiStater(*collectorWifi)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debug("wifi collector metrics are not available for this system")
			return nil
		}
		if os.IsPermission(err) {
			log.Debug("wifi collector got permission denied when accessing metrics")
			return nil
		}
		return fmt.Errorf("failed to access wifi data: %v", err)
	}
	defer stat.Close()
	ifis, err := stat.Interfaces()
	if err != nil {
		return fmt.Errorf("failed to retrieve wifi interfaces: %v", err)
	}
	for _, ifi := range ifis {
		if ifi.Name == "" {
			continue
		}
		log.Debugf("probing wifi device %q with type %q", ifi.Name, ifi.Type)
		ch <- prometheus.MustNewConstMetric(c.interfaceFrequencyHertz, prometheus.GaugeValue, mHzToHz(ifi.Frequency), ifi.Name)
		bss, err := stat.BSS(ifi)
		switch {
		case err == nil:
			c.updateBSSStats(ch, ifi.Name, bss)
		case os.IsNotExist(err):
			log.Debugf("BSS information not found for wifi device %q", ifi.Name)
		default:
			return fmt.Errorf("failed to retrieve BSS for device %s: %v", ifi.Name, err)
		}
		stations, err := stat.StationInfo(ifi)
		switch {
		case err == nil:
			for _, station := range stations {
				c.updateStationStats(ch, ifi.Name, station)
			}
		case os.IsNotExist(err):
			log.Debugf("station information not found for wifi device %q", ifi.Name)
		default:
			return fmt.Errorf("failed to retrieve station info for device %q: %v", ifi.Name, err)
		}
	}
	return nil
}
func (c *wifiCollector) updateBSSStats(ch chan<- prometheus.Metric, device string, bss *wifi.BSS) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ch <- prometheus.MustNewConstMetric(c.stationInfo, prometheus.GaugeValue, 1, device, bss.BSSID.String(), bss.SSID, bssStatusMode(bss.Status))
}
func (c *wifiCollector) updateStationStats(ch chan<- prometheus.Metric, device string, info *wifi.StationInfo) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ch <- prometheus.MustNewConstMetric(c.stationConnectedSecondsTotal, prometheus.CounterValue, info.Connected.Seconds(), device, info.HardwareAddr.String())
	ch <- prometheus.MustNewConstMetric(c.stationInactiveSeconds, prometheus.GaugeValue, info.Inactive.Seconds(), device, info.HardwareAddr.String())
	ch <- prometheus.MustNewConstMetric(c.stationReceiveBitsPerSecond, prometheus.GaugeValue, float64(info.ReceiveBitrate), device, info.HardwareAddr.String())
	ch <- prometheus.MustNewConstMetric(c.stationTransmitBitsPerSecond, prometheus.GaugeValue, float64(info.TransmitBitrate), device, info.HardwareAddr.String())
	ch <- prometheus.MustNewConstMetric(c.stationSignalDBM, prometheus.GaugeValue, float64(info.Signal), device, info.HardwareAddr.String())
	ch <- prometheus.MustNewConstMetric(c.stationTransmitRetriesTotal, prometheus.CounterValue, float64(info.TransmitRetries), device, info.HardwareAddr.String())
	ch <- prometheus.MustNewConstMetric(c.stationTransmitFailedTotal, prometheus.CounterValue, float64(info.TransmitFailed), device, info.HardwareAddr.String())
	ch <- prometheus.MustNewConstMetric(c.stationBeaconLossTotal, prometheus.CounterValue, float64(info.BeaconLoss), device, info.HardwareAddr.String())
}
func mHzToHz(mHz int) float64 {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return float64(mHz) * 1000 * 1000
}
func bssStatusMode(status wifi.BSSStatus) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	switch status {
	case wifi.BSSStatusAuthenticated, wifi.BSSStatusAssociated:
		return "client"
	case wifi.BSSStatusIBSSJoined:
		return "ad-hoc"
	default:
		return "unknown"
	}
}
func newWifiStater(fixtures string) (wifiStater, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if fixtures != "" {
		return &mockWifiStater{fixtures: fixtures}, nil
	}
	return wifi.New()
}

var _ wifiStater = &mockWifiStater{}

type mockWifiStater struct{ fixtures string }

func (s *mockWifiStater) unmarshalJSONFile(filename string, v interface{}) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	b, err := ioutil.ReadFile(filepath.Join(s.fixtures, filename))
	if err != nil {
		return err
	}
	return json.Unmarshal(b, v)
}
func (s *mockWifiStater) Close() error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return nil
}
func (s *mockWifiStater) BSS(ifi *wifi.Interface) (*wifi.BSS, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p := filepath.Join(ifi.Name, "bss.json")
	var bss wifi.BSS
	if err := s.unmarshalJSONFile(p, &bss); err != nil {
		return nil, err
	}
	return &bss, nil
}
func (s *mockWifiStater) Interfaces() ([]*wifi.Interface, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var ifis []*wifi.Interface
	if err := s.unmarshalJSONFile("interfaces.json", &ifis); err != nil {
		return nil, err
	}
	return ifis, nil
}
func (s *mockWifiStater) StationInfo(ifi *wifi.Interface) ([]*wifi.StationInfo, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	p := filepath.Join(ifi.Name, "stationinfo.json")
	var stations []*wifi.StationInfo
	if err := s.unmarshalJSONFile(p, &stations); err != nil {
		return nil, err
	}
	return stations, nil
}
