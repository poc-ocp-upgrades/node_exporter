package collector

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/prometheus/common/log"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	textFileDirectory	= kingpin.Flag("collector.textfile.directory", "Directory to read text files with metrics from.").Default("").String()
	mtimeDesc		= prometheus.NewDesc("node_textfile_mtime_seconds", "Unixtime mtime of textfiles successfully read.", []string{"file"}, nil)
)

type textFileCollector struct {
	path	string
	mtime	*float64
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("textfile", defaultEnabled, NewTextFileCollector)
}
func NewTextFileCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	c := &textFileCollector{path: *textFileDirectory}
	return c, nil
}
func convertMetricFamily(metricFamily *dto.MetricFamily, ch chan<- prometheus.Metric) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var valType prometheus.ValueType
	var val float64
	allLabelNames := map[string]struct{}{}
	for _, metric := range metricFamily.Metric {
		labels := metric.GetLabel()
		for _, label := range labels {
			if _, ok := allLabelNames[label.GetName()]; !ok {
				allLabelNames[label.GetName()] = struct{}{}
			}
		}
	}
	for _, metric := range metricFamily.Metric {
		if metric.TimestampMs != nil {
			log.Warnf("Ignoring unsupported custom timestamp on textfile collector metric %v", metric)
		}
		labels := metric.GetLabel()
		var names []string
		var values []string
		for _, label := range labels {
			names = append(names, label.GetName())
			values = append(values, label.GetValue())
		}
		for k := range allLabelNames {
			present := false
			for _, name := range names {
				if k == name {
					present = true
					break
				}
			}
			if present == false {
				names = append(names, k)
				values = append(values, "")
			}
		}
		metricType := metricFamily.GetType()
		switch metricType {
		case dto.MetricType_COUNTER:
			valType = prometheus.CounterValue
			val = metric.Counter.GetValue()
		case dto.MetricType_GAUGE:
			valType = prometheus.GaugeValue
			val = metric.Gauge.GetValue()
		case dto.MetricType_UNTYPED:
			valType = prometheus.UntypedValue
			val = metric.Untyped.GetValue()
		case dto.MetricType_SUMMARY:
			quantiles := map[float64]float64{}
			for _, q := range metric.Summary.Quantile {
				quantiles[q.GetQuantile()] = q.GetValue()
			}
			ch <- prometheus.MustNewConstSummary(prometheus.NewDesc(*metricFamily.Name, metricFamily.GetHelp(), names, nil), metric.Summary.GetSampleCount(), metric.Summary.GetSampleSum(), quantiles, values...)
		case dto.MetricType_HISTOGRAM:
			buckets := map[float64]uint64{}
			for _, b := range metric.Histogram.Bucket {
				buckets[b.GetUpperBound()] = b.GetCumulativeCount()
			}
			ch <- prometheus.MustNewConstHistogram(prometheus.NewDesc(*metricFamily.Name, metricFamily.GetHelp(), names, nil), metric.Histogram.GetSampleCount(), metric.Histogram.GetSampleSum(), buckets, values...)
		default:
			panic("unknown metric type")
		}
		if metricType == dto.MetricType_GAUGE || metricType == dto.MetricType_COUNTER || metricType == dto.MetricType_UNTYPED {
			ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(*metricFamily.Name, metricFamily.GetHelp(), names, nil), valType, val, values...)
		}
	}
}
func (c *textFileCollector) exportMTimes(mtimes map[string]time.Time, ch chan<- prometheus.Metric) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if len(mtimes) > 0 {
		filenames := make([]string, 0, len(mtimes))
		for filename := range mtimes {
			filenames = append(filenames, filename)
		}
		sort.Strings(filenames)
		for _, filename := range filenames {
			mtime := float64(mtimes[filename].UnixNano() / 1e9)
			if c.mtime != nil {
				mtime = *c.mtime
			}
			ch <- prometheus.MustNewConstMetric(mtimeDesc, prometheus.GaugeValue, mtime, filename)
		}
	}
}
func (c *textFileCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	error := 0.0
	mtimes := map[string]time.Time{}
	files, err := ioutil.ReadDir(c.path)
	if err != nil && c.path != "" {
		log.Errorf("Error reading textfile collector directory %q: %s", c.path, err)
		error = 1.0
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".prom") {
			continue
		}
		path := filepath.Join(c.path, f.Name())
		file, err := os.Open(path)
		if err != nil {
			log.Errorf("Error opening %q: %v", path, err)
			error = 1.0
			continue
		}
		var parser expfmt.TextParser
		parsedFamilies, err := parser.TextToMetricFamilies(file)
		file.Close()
		if err != nil {
			log.Errorf("Error parsing %q: %v", path, err)
			error = 1.0
			continue
		}
		if hasTimestamps(parsedFamilies) {
			log.Errorf("Textfile %q contains unsupported client-side timestamps, skipping entire file", path)
			error = 1.0
			continue
		}
		for _, mf := range parsedFamilies {
			if mf.Help == nil {
				help := fmt.Sprintf("Metric read from %s", path)
				mf.Help = &help
			}
		}
		mtimes[f.Name()] = f.ModTime()
		for _, mf := range parsedFamilies {
			convertMetricFamily(mf, ch)
		}
	}
	c.exportMTimes(mtimes, ch)
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc("node_textfile_scrape_error", "1 if there was an error opening or reading a file, 0 otherwise", nil, nil), prometheus.GaugeValue, error)
	return nil
}
func hasTimestamps(parsedFamilies map[string]*dto.MetricFamily) bool {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, mf := range parsedFamilies {
		for _, m := range mf.Metric {
			if m.TimestampMs != nil {
				return true
			}
		}
	}
	return false
}
