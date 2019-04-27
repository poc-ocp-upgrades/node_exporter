package collector

import "testing"

func TestFileFDStats(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	fileFDStats, err := parseFileFDStats("fixtures/proc/sys/fs/file-nr")
	if err != nil {
		t.Fatal(err)
	}
	if want, got := "1024", fileFDStats["allocated"]; want != got {
		t.Errorf("want filefd allocated %q, got %q", want, got)
	}
	if want, got := "1631329", fileFDStats["maximum"]; want != got {
		t.Errorf("want filefd maximum %q, got %q", want, got)
	}
}
