package collector

import (
	"testing"
	"gopkg.in/alecthomas/kingpin.v2"
)

func TestReadProcessStatus(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	if _, err := kingpin.CommandLine.Parse([]string{"--path.procfs", "fixtures/proc"}); err != nil {
		t.Fatal(err)
	}
	want := 1
	pids, states, threads, err := getAllocatedThreads()
	if err != nil {
		t.Fatalf("Cannot retrieve data from procfs getAllocatedThreads function: %v ", err)
	}
	if threads < want {
		t.Fatalf("Current threads: %d Shouldn't be less than wanted %d", threads, want)
	}
	if states == nil {
		t.Fatalf("Process states cannot be nil %v:", states)
	}
	maxPid, err := readUintFromFile(procFilePath("sys/kernel/pid_max"))
	if err != nil {
		t.Fatalf("Unable to retrieve limit number of maximum pids alloved %v\n", err)
	}
	if uint64(pids) > maxPid || pids == 0 {
		t.Fatalf("Total running pids cannot be greater than %d or equals to 0", maxPid)
	}
}
