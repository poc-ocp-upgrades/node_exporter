package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	cpuCollectorSubsystem = "cpu"
)

var (
	nodeCPUSecondsDesc = prometheus.NewDesc(prometheus.BuildFQName(namespace, cpuCollectorSubsystem, "seconds_total"), "Seconds the cpus spent in each mode.", []string{"cpu", "mode"}, nil)
)
