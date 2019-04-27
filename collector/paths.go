package collector

import (
	"path"
	"github.com/prometheus/procfs"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	procPath	= kingpin.Flag("path.procfs", "procfs mountpoint.").Default(procfs.DefaultMountPoint).String()
	sysPath		= kingpin.Flag("path.sysfs", "sysfs mountpoint.").Default("/sys").String()
	rootfsPath	= kingpin.Flag("path.rootfs", "rootfs mountpoint.").Default("/").String()
)

func procFilePath(name string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return path.Join(*procPath, name)
}
func sysFilePath(name string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return path.Join(*sysPath, name)
}
func rootfsFilePath(name string) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return path.Join(*rootfsPath, name)
}
