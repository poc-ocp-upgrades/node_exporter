package collector

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type collectorAdapter struct{ Collector }

func (a collectorAdapter) Describe(ch chan<- *prometheus.Desc) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	ch <- prometheus.NewDesc("dummy_metric", "Dummy metric.", nil, nil)
}
func (a collectorAdapter) Collect(ch chan<- prometheus.Metric) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	err := a.Update(ch)
	if err != nil {
		panic(fmt.Sprintf("failed to update collector: %v", err))
	}
}
func TestTextfileCollector(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	tests := []struct {
		path	string
		out	string
	}{{path: "fixtures/textfile/no_metric_files", out: "fixtures/textfile/no_metric_files.out"}, {path: "fixtures/textfile/two_metric_files", out: "fixtures/textfile/two_metric_files.out"}, {path: "fixtures/textfile/nonexistent_path", out: "fixtures/textfile/nonexistent_path.out"}, {path: "fixtures/textfile/client_side_timestamp", out: "fixtures/textfile/client_side_timestamp.out"}, {path: "fixtures/textfile/different_metric_types", out: "fixtures/textfile/different_metric_types.out"}, {path: "fixtures/textfile/inconsistent_metrics", out: "fixtures/textfile/inconsistent_metrics.out"}, {path: "fixtures/textfile/histogram", out: "fixtures/textfile/histogram.out"}, {path: "fixtures/textfile/histogram_extra_dimension", out: "fixtures/textfile/histogram_extra_dimension.out"}, {path: "fixtures/textfile/summary", out: "fixtures/textfile/summary.out"}, {path: "fixtures/textfile/summary_extra_dimension", out: "fixtures/textfile/summary_extra_dimension.out"}}
	for i, test := range tests {
		mtime := 1.0
		c := &textFileCollector{path: test.path, mtime: &mtime}
		log.AddFlags(kingpin.CommandLine)
		_, err := kingpin.CommandLine.Parse([]string{"--log.level", "fatal"})
		if err != nil {
			t.Fatal(err)
		}
		registry := prometheus.NewRegistry()
		registry.MustRegister(collectorAdapter{c})
		rw := httptest.NewRecorder()
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{}).ServeHTTP(rw, &http.Request{})
		got := string(rw.Body.String())
		want, err := ioutil.ReadFile(test.out)
		if err != nil {
			t.Fatalf("%d. error reading fixture file %s: %s", i, test.out, err)
		}
		if string(want) != got {
			t.Fatalf("%d.%q want:\n\n%s\n\ngot:\n\n%s", i, test.path, string(want), got)
		}
	}
}
