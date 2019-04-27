package collector

import (
	"testing"
	"github.com/prometheus/procfs"
	"gopkg.in/alecthomas/kingpin.v2"
)

func TestDefaultProcPath(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if _, err := kingpin.CommandLine.Parse([]string{"--path.procfs", procfs.DefaultMountPoint}); err != nil {
		t.Fatal(err)
	}
	if got, want := procFilePath("somefile"), "/proc/somefile"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
	if got, want := procFilePath("some/file"), "/proc/some/file"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
}
func TestCustomProcPath(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if _, err := kingpin.CommandLine.Parse([]string{"--path.procfs", "./../some/./place/"}); err != nil {
		t.Fatal(err)
	}
	if got, want := procFilePath("somefile"), "../some/place/somefile"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
	if got, want := procFilePath("some/file"), "../some/place/some/file"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
}
func TestDefaultSysPath(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if _, err := kingpin.CommandLine.Parse([]string{"--path.sysfs", "/sys"}); err != nil {
		t.Fatal(err)
	}
	if got, want := sysFilePath("somefile"), "/sys/somefile"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
	if got, want := sysFilePath("some/file"), "/sys/some/file"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
}
func TestCustomSysPath(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	if _, err := kingpin.CommandLine.Parse([]string{"--path.sysfs", "./../some/./place/"}); err != nil {
		t.Fatal(err)
	}
	if got, want := sysFilePath("somefile"), "../some/place/somefile"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
	if got, want := sysFilePath("some/file"), "../some/place/some/file"; got != want {
		t.Errorf("Expected: %s, Got: %s", want, got)
	}
}
