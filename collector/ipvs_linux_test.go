package collector

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

func TestIPVSCollector(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if _, err := kingpin.CommandLine.Parse([]string{"--path.procfs", "fixtures/proc"}); err != nil {
		t.Fatal(err)
	}
	collector, err := newIPVSCollector()
	if err != nil {
		t.Fatal(err)
	}
	sink := make(chan prometheus.Metric)
	go func() {
		err = collector.Update(sink)
		if err != nil {
			panic(fmt.Sprintf("failed to update collector: %v", err))
		}
	}()
	for expected, got := range map[string]string{prometheus.NewDesc("node_ipvs_connections_total", "The total number of connections made.", nil, nil).String(): (<-sink).Desc().String(), prometheus.NewDesc("node_ipvs_incoming_packets_total", "The total number of incoming packets.", nil, nil).String(): (<-sink).Desc().String(), prometheus.NewDesc("node_ipvs_outgoing_packets_total", "The total number of outgoing packets.", nil, nil).String(): (<-sink).Desc().String(), prometheus.NewDesc("node_ipvs_incoming_bytes_total", "The total amount of incoming data.", nil, nil).String(): (<-sink).Desc().String(), prometheus.NewDesc("node_ipvs_outgoing_bytes_total", "The total amount of outgoing data.", nil, nil).String(): (<-sink).Desc().String(), prometheus.NewDesc("node_ipvs_backend_connections_active", "The current active connections by local and remote address.", []string{"local_address", "local_port", "remote_address", "remote_port", "proto"}, nil).String(): (<-sink).Desc().String(), prometheus.NewDesc("node_ipvs_backend_connections_inactive", "The current inactive connections by local and remote address.", []string{"local_address", "local_port", "remote_address", "remote_port", "proto"}, nil).String(): (<-sink).Desc().String(), prometheus.NewDesc("node_ipvs_backend_weight", "The current backend weight by local and remote address.", []string{"local_address", "local_port", "remote_address", "remote_port", "proto"}, nil).String(): (<-sink).Desc().String()} {
		if expected != got {
			t.Fatalf("Expected '%s' but got '%s'", expected, got)
		}
	}
}

type miniCollector struct{ c Collector }

func (c miniCollector) Collect(ch chan<- prometheus.Metric) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	c.c.Update(ch)
}
func (c miniCollector) Describe(ch chan<- *prometheus.Desc) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	prometheus.NewGauge(prometheus.GaugeOpts{Namespace: "fake", Subsystem: "fake", Name: "fake", Help: "fake"}).Describe(ch)
}
func TestIPVSCollectorResponse(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if _, err := kingpin.CommandLine.Parse([]string{"--path.procfs", "fixtures/proc"}); err != nil {
		t.Fatal(err)
	}
	collector, err := NewIPVSCollector()
	if err != nil {
		t.Fatal(err)
	}
	prometheus.MustRegister(miniCollector{c: collector})
	rw := httptest.NewRecorder()
	promhttp.Handler().ServeHTTP(rw, &http.Request{})
	metricsFile := "fixtures/ip_vs_result.txt"
	wantMetrics, err := ioutil.ReadFile(metricsFile)
	if err != nil {
		t.Fatalf("unable to read input test file %s: %s", metricsFile, err)
	}
	wantLines := strings.Split(string(wantMetrics), "\n")
	gotLines := strings.Split(string(rw.Body.String()), "\n")
	gotLinesIdx := 0
wantLoop:
	for _, want := range wantLines {
		for _, got := range gotLines[gotLinesIdx:] {
			if want == got {
				continue wantLoop
			} else {
				gotLinesIdx++
			}
		}
		t.Fatalf("Missing expected output line(s), first missing line is %s", want)
	}
}
