package collector

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
)

var (
	hwmonInvalidMetricChars	= regexp.MustCompile("[^a-z0-9:_]")
	hwmonFilenameFormat	= regexp.MustCompile(`^(?P<type>[^0-9]+)(?P<id>[0-9]*)?(_(?P<property>.+))?$`)
	hwmonLabelDesc		= []string{"chip", "sensor"}
	hwmonChipNameLabelDesc	= []string{"chip", "chip_name"}
	hwmonSensorTypes	= []string{"vrm", "beep_enable", "update_interval", "in", "cpu", "fan", "pwm", "temp", "curr", "power", "energy", "humidity", "intrusion"}
)

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("hwmon", defaultEnabled, NewHwMonCollector)
}

type hwMonCollector struct{}

func NewHwMonCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &hwMonCollector{}, nil
}
func cleanMetricName(name string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	lower := strings.ToLower(name)
	replaced := hwmonInvalidMetricChars.ReplaceAllLiteralString(lower, "_")
	cleaned := strings.Trim(replaced, "_")
	return cleaned
}
func addValueFile(data map[string]map[string]string, sensor string, prop string, file string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	raw, err := sysReadFile(file)
	if err != nil {
		return
	}
	value := strings.Trim(string(raw), "\n")
	if _, ok := data[sensor]; !ok {
		data[sensor] = make(map[string]string)
	}
	data[sensor][prop] = value
}
func sysReadFile(file string) ([]byte, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b := make([]byte, 128)
	n, err := syscall.Read(int(f.Fd()), b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func explodeSensorFilename(filename string) (ok bool, sensorType string, sensorNum int, sensorProperty string) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	matches := hwmonFilenameFormat.FindStringSubmatch(filename)
	if len(matches) == 0 {
		return false, sensorType, sensorNum, sensorProperty
	}
	for i, match := range hwmonFilenameFormat.SubexpNames() {
		if i >= len(matches) {
			return true, sensorType, sensorNum, sensorProperty
		}
		if match == "type" {
			sensorType = matches[i]
		}
		if match == "property" {
			sensorProperty = matches[i]
		}
		if match == "id" && len(matches[i]) > 0 {
			if num, err := strconv.Atoi(matches[i]); err == nil {
				sensorNum = num
			} else {
				return false, sensorType, sensorNum, sensorProperty
			}
		}
	}
	return true, sensorType, sensorNum, sensorProperty
}
func collectSensorData(dir string, data map[string]map[string]string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sensorFiles, dirError := ioutil.ReadDir(dir)
	if dirError != nil {
		return dirError
	}
	for _, file := range sensorFiles {
		filename := file.Name()
		ok, sensorType, sensorNum, sensorProperty := explodeSensorFilename(filename)
		if !ok {
			continue
		}
		for _, t := range hwmonSensorTypes {
			if t == sensorType {
				addValueFile(data, sensorType+strconv.Itoa(sensorNum), sensorProperty, path.Join(dir, file.Name()))
				break
			}
		}
	}
	return nil
}
func (c *hwMonCollector) updateHwmon(ch chan<- prometheus.Metric, dir string) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	hwmonName, err := c.hwmonName(dir)
	if err != nil {
		return err
	}
	data := make(map[string]map[string]string)
	err = collectSensorData(dir, data)
	if err != nil {
		return err
	}
	if _, err := os.Stat(path.Join(dir, "device")); err == nil {
		err := collectSensorData(path.Join(dir, "device"), data)
		if err != nil {
			return err
		}
	}
	hwmonChipName, err := c.hwmonHumanReadableChipName(dir)
	if err == nil {
		desc := prometheus.NewDesc("node_hwmon_chip_names", "Annotation metric for human-readable chip names", hwmonChipNameLabelDesc, nil)
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, 1.0, hwmonName, hwmonChipName)
	}
	for sensor, sensorData := range data {
		_, sensorType, _, _ := explodeSensorFilename(sensor)
		labels := []string{hwmonName, sensor}
		if labelText, ok := sensorData["label"]; ok {
			label := cleanMetricName(labelText)
			if label != "" {
				desc := prometheus.NewDesc("node_hwmon_sensor_label", "Label for given chip and sensor", []string{"chip", "sensor", "label"}, nil)
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, 1.0, hwmonName, sensor, label)
			}
		}
		if sensorType == "beep_enable" {
			value := 0.0
			if sensorData[""] == "1" {
				value = 1.0
			}
			metricName := "node_hwmon_beep_enabled"
			desc := prometheus.NewDesc(metricName, "Hardware beep enabled", hwmonLabelDesc, nil)
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, value, labels...)
			continue
		}
		if sensorType == "vrm" {
			parsedValue, err := strconv.ParseFloat(sensorData[""], 64)
			if err != nil {
				continue
			}
			metricName := "node_hwmon_voltage_regulator_version"
			desc := prometheus.NewDesc(metricName, "Hardware voltage regulator", hwmonLabelDesc, nil)
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, parsedValue, labels...)
			continue
		}
		if sensorType == "update_interval" {
			parsedValue, err := strconv.ParseFloat(sensorData[""], 64)
			if err != nil {
				continue
			}
			metricName := "node_hwmon_update_interval_seconds"
			desc := prometheus.NewDesc(metricName, "Hardware monitor update interval", hwmonLabelDesc, nil)
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, parsedValue*0.001, labels...)
			continue
		}
		prefix := "node_hwmon_" + sensorType
		for element, value := range sensorData {
			if element == "label" {
				continue
			}
			name := prefix
			if element == "input" {
				if _, ok := sensorData[""]; ok {
					name = name + "_input"
				}
			} else if element != "" {
				name = name + "_" + cleanMetricName(element)
			}
			parsedValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				continue
			}
			if element == "fault" || element == "alarm" {
				desc := prometheus.NewDesc(name, "Hardware sensor "+element+" status ("+sensorType+")", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, parsedValue, labels...)
				continue
			}
			if element == "beep" {
				desc := prometheus.NewDesc(name+"_enabled", "Hardware monitor sensor has beeping enabled", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, parsedValue, labels...)
				continue
			}
			if sensorType == "in" || sensorType == "cpu" {
				desc := prometheus.NewDesc(name+"_volts", "Hardware monitor for voltage ("+element+")", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, parsedValue*0.001, labels...)
				continue
			}
			if sensorType == "temp" && element != "type" {
				if element == "" {
					element = "input"
				}
				desc := prometheus.NewDesc(name+"_celsius", "Hardware monitor for temperature ("+element+")", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, parsedValue*0.001, labels...)
				continue
			}
			if sensorType == "curr" {
				desc := prometheus.NewDesc(name+"_amps", "Hardware monitor for current ("+element+")", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, parsedValue*0.001, labels...)
				continue
			}
			if sensorType == "energy" {
				desc := prometheus.NewDesc(name+"_joule_total", "Hardware monitor for joules used so far ("+element+")", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, parsedValue/1000000.0, labels...)
				continue
			}
			if sensorType == "power" && element == "accuracy" {
				desc := prometheus.NewDesc(name, "Hardware monitor power meter accuracy, as a ratio", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, parsedValue/1000000.0, labels...)
				continue
			}
			if sensorType == "power" && (element == "average_interval" || element == "average_interval_min" || element == "average_interval_max") {
				desc := prometheus.NewDesc(name+"_seconds", "Hardware monitor power usage update interval ("+element+")", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, parsedValue*0.001, labels...)
				continue
			}
			if sensorType == "power" {
				desc := prometheus.NewDesc(name+"_watt", "Hardware monitor for power usage in watts ("+element+")", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, parsedValue/1000000.0, labels...)
				continue
			}
			if sensorType == "humidity" {
				desc := prometheus.NewDesc(name, "Hardware monitor for humidity, as a ratio (multiply with 100.0 to get the humidity as a percentage) ("+element+")", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, parsedValue/1000000.0, labels...)
				continue
			}
			if sensorType == "fan" && (element == "input" || element == "min" || element == "max" || element == "target") {
				desc := prometheus.NewDesc(name+"_rpm", "Hardware monitor for fan revolutions per minute ("+element+")", hwmonLabelDesc, nil)
				ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, parsedValue, labels...)
				continue
			}
			desc := prometheus.NewDesc(name, "Hardware monitor "+sensorType+" element "+element, hwmonLabelDesc, nil)
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, parsedValue, labels...)
		}
	}
	return nil
}
func (c *hwMonCollector) hwmonName(dir string) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	devicePath, devErr := filepath.EvalSymlinks(path.Join(dir, "device"))
	if devErr == nil {
		devPathPrefix, devName := path.Split(devicePath)
		_, devType := path.Split(strings.TrimRight(devPathPrefix, "/"))
		cleanDevName := cleanMetricName(devName)
		cleanDevType := cleanMetricName(devType)
		if cleanDevType != "" && cleanDevName != "" {
			return cleanDevType + "_" + cleanDevName, nil
		}
		if cleanDevName != "" {
			return cleanDevName, nil
		}
	}
	sysnameRaw, nameErr := ioutil.ReadFile(path.Join(dir, "name"))
	if nameErr == nil && string(sysnameRaw) != "" {
		cleanName := cleanMetricName(string(sysnameRaw))
		if cleanName != "" {
			return cleanName, nil
		}
	}
	realDir, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return "", err
	}
	_, name := path.Split(realDir)
	cleanName := cleanMetricName(name)
	if cleanName != "" {
		return cleanName, nil
	}
	return "", errors.New("Could not derive a monitoring name for " + dir)
}
func (c *hwMonCollector) hwmonHumanReadableChipName(dir string) (string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	sysnameRaw, nameErr := ioutil.ReadFile(path.Join(dir, "name"))
	if nameErr != nil {
		return "", nameErr
	}
	if string(sysnameRaw) != "" {
		cleanName := cleanMetricName(string(sysnameRaw))
		if cleanName != "" {
			return cleanName, nil
		}
	}
	return "", errors.New("Could not derive a human-readable chip type for " + dir)
}
func (c *hwMonCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	hwmonPathName := path.Join(sysFilePath("class"), "hwmon")
	hwmonFiles, err := ioutil.ReadDir(hwmonPathName)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debug("hwmon collector metrics are not available for this system")
			return nil
		}
		return err
	}
	for _, hwDir := range hwmonFiles {
		hwmonXPathName := path.Join(hwmonPathName, hwDir.Name())
		if hwDir.Mode()&os.ModeSymlink > 0 {
			hwDir, err = os.Stat(hwmonXPathName)
			if err != nil {
				continue
			}
		}
		if !hwDir.IsDir() {
			continue
		}
		if lastErr := c.updateHwmon(ch, hwmonXPathName); lastErr != nil {
			err = lastErr
		}
	}
	return err
}
