package collector

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	edacSubsystem = "edac"
)

var (
	edacMemControllerRE	= regexp.MustCompile(`.*devices/system/edac/mc/mc([0-9]*)`)
	edacMemCsrowRE		= regexp.MustCompile(`.*devices/system/edac/mc/mc[0-9]*/csrow([0-9]*)`)
)

type edacCollector struct {
	ceCount			*prometheus.Desc
	ueCount			*prometheus.Desc
	csRowCECount	*prometheus.Desc
	csRowUECount	*prometheus.Desc
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("edac", defaultEnabled, NewEdacCollector)
}
func NewEdacCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return &edacCollector{ceCount: prometheus.NewDesc(prometheus.BuildFQName(namespace, edacSubsystem, "correctable_errors_total"), "Total correctable memory errors.", []string{"controller"}, nil), ueCount: prometheus.NewDesc(prometheus.BuildFQName(namespace, edacSubsystem, "uncorrectable_errors_total"), "Total uncorrectable memory errors.", []string{"controller"}, nil), csRowCECount: prometheus.NewDesc(prometheus.BuildFQName(namespace, edacSubsystem, "csrow_correctable_errors_total"), "Total correctable memory errors for this csrow.", []string{"controller", "csrow"}, nil), csRowUECount: prometheus.NewDesc(prometheus.BuildFQName(namespace, edacSubsystem, "csrow_uncorrectable_errors_total"), "Total uncorrectable memory errors for this csrow.", []string{"controller", "csrow"}, nil)}, nil
}
func (c *edacCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	memControllers, err := filepath.Glob(sysFilePath("devices/system/edac/mc/mc[0-9]*"))
	if err != nil {
		return err
	}
	for _, controller := range memControllers {
		controllerMatch := edacMemControllerRE.FindStringSubmatch(controller)
		if controllerMatch == nil {
			return fmt.Errorf("controller string didn't match regexp: %s", controller)
		}
		controllerNumber := controllerMatch[1]
		value, err := readUintFromFile(path.Join(controller, "ce_count"))
		if err != nil {
			return fmt.Errorf("couldn't get ce_count for controller %s: %s", controllerNumber, err)
		}
		ch <- prometheus.MustNewConstMetric(c.ceCount, prometheus.CounterValue, float64(value), controllerNumber)
		value, err = readUintFromFile(path.Join(controller, "ce_noinfo_count"))
		if err != nil {
			return fmt.Errorf("couldn't get ce_noinfo_count for controller %s: %s", controllerNumber, err)
		}
		ch <- prometheus.MustNewConstMetric(c.csRowCECount, prometheus.CounterValue, float64(value), controllerNumber, "unknown")
		value, err = readUintFromFile(path.Join(controller, "ue_count"))
		if err != nil {
			return fmt.Errorf("couldn't get ue_count for controller %s: %s", controllerNumber, err)
		}
		ch <- prometheus.MustNewConstMetric(c.ueCount, prometheus.CounterValue, float64(value), controllerNumber)
		value, err = readUintFromFile(path.Join(controller, "ue_noinfo_count"))
		if err != nil {
			return fmt.Errorf("couldn't get ue_noinfo_count for controller %s: %s", controllerNumber, err)
		}
		ch <- prometheus.MustNewConstMetric(c.csRowUECount, prometheus.CounterValue, float64(value), controllerNumber, "unknown")
		csrows, err := filepath.Glob(controller + "/csrow[0-9]*")
		if err != nil {
			return err
		}
		for _, csrow := range csrows {
			csrowMatch := edacMemCsrowRE.FindStringSubmatch(csrow)
			if csrowMatch == nil {
				return fmt.Errorf("csrow string didn't match regexp: %s", csrow)
			}
			csrowNumber := csrowMatch[1]
			value, err = readUintFromFile(path.Join(csrow, "ce_count"))
			if err != nil {
				return fmt.Errorf("couldn't get ce_count for controller/csrow %s/%s: %s", controllerNumber, csrowNumber, err)
			}
			ch <- prometheus.MustNewConstMetric(c.csRowCECount, prometheus.CounterValue, float64(value), controllerNumber, csrowNumber)
			value, err = readUintFromFile(path.Join(csrow, "ue_count"))
			if err != nil {
				return fmt.Errorf("couldn't get ue_count for controller/csrow %s/%s: %s", controllerNumber, csrowNumber, err)
			}
			ch <- prometheus.MustNewConstMetric(c.csRowUECount, prometheus.CounterValue, float64(value), controllerNumber, csrowNumber)
		}
	}
	return err
}
