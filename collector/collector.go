package collector

import (
	"fmt"
	"sync"
	"time"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"gopkg.in/alecthomas/kingpin.v2"
)

const namespace = "node"

var (
	scrapeDurationDesc	= prometheus.NewDesc(prometheus.BuildFQName(namespace, "scrape", "collector_duration_seconds"), "node_exporter: Duration of a collector scrape.", []string{"collector"}, nil)
	scrapeSuccessDesc	= prometheus.NewDesc(prometheus.BuildFQName(namespace, "scrape", "collector_success"), "node_exporter: Whether a collector succeeded.", []string{"collector"}, nil)
)

const (
	defaultEnabled	= true
	defaultDisabled	= false
)

var (
	factories	= make(map[string]func() (Collector, error))
	collectorState	= make(map[string]*bool)
)

func registerCollector(collector string, isDefaultEnabled bool, factory func() (Collector, error)) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	var helpDefaultState string
	if isDefaultEnabled {
		helpDefaultState = "enabled"
	} else {
		helpDefaultState = "disabled"
	}
	flagName := fmt.Sprintf("collector.%s", collector)
	flagHelp := fmt.Sprintf("Enable the %s collector (default: %s).", collector, helpDefaultState)
	defaultValue := fmt.Sprintf("%v", isDefaultEnabled)
	flag := kingpin.Flag(flagName, flagHelp).Default(defaultValue).Bool()
	collectorState[collector] = flag
	factories[collector] = factory
}

type nodeCollector struct{ Collectors map[string]Collector }

func NewNodeCollector(filters ...string) (*nodeCollector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	f := make(map[string]bool)
	for _, filter := range filters {
		enabled, exist := collectorState[filter]
		if !exist {
			return nil, fmt.Errorf("missing collector: %s", filter)
		}
		if !*enabled {
			return nil, fmt.Errorf("disabled collector: %s", filter)
		}
		f[filter] = true
	}
	collectors := make(map[string]Collector)
	for key, enabled := range collectorState {
		if *enabled {
			collector, err := factories[key]()
			if err != nil {
				return nil, err
			}
			if len(f) == 0 || f[key] {
				collectors[key] = collector
			}
		}
	}
	return &nodeCollector{Collectors: collectors}, nil
}
func (n nodeCollector) Describe(ch chan<- *prometheus.Desc) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}
func (n nodeCollector) Collect(ch chan<- prometheus.Metric) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	wg := sync.WaitGroup{}
	wg.Add(len(n.Collectors))
	for name, c := range n.Collectors {
		go func(name string, c Collector) {
			execute(name, c, ch)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
}
func execute(name string, c Collector, ch chan<- prometheus.Metric) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	begin := time.Now()
	err := c.Update(ch)
	duration := time.Since(begin)
	var success float64
	if err != nil {
		log.Errorf("ERROR: %s collector failed after %fs: %s", name, duration.Seconds(), err)
		success = 0
	} else {
		log.Debugf("OK: %s collector succeeded after %fs.", name, duration.Seconds())
		success = 1
	}
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)
}

type Collector interface {
	Update(ch chan<- prometheus.Metric) error
}
type typedDesc struct {
	desc		*prometheus.Desc
	valueType	prometheus.ValueType
}

func (d *typedDesc) mustNewConstMetric(value float64, labels ...string) prometheus.Metric {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return prometheus.MustNewConstMetric(d.desc, d.valueType, value, labels...)
}
