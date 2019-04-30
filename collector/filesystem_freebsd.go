package collector

import (
	"bytes"
	"unsafe"
	"github.com/prometheus/common/log"
	"golang.org/x/sys/unix"
)

const (
	defIgnoredMountPoints	= "^/(dev)($|/)"
	defIgnoredFSTypes	= "^devfs$"
	readOnly		= 0x1
	noWait			= 0x2
)

func gostring(b []int8) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	bb := *(*[]byte)(unsafe.Pointer(&b))
	idx := bytes.IndexByte(bb, 0)
	if idx < 0 {
		return ""
	}
	return string(bb[:idx])
}
func (c *filesystemCollector) GetStats() ([]filesystemStats, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	n, err := unix.Getfsstat(nil, noWait)
	if err != nil {
		return nil, err
	}
	buf := make([]unix.Statfs_t, n)
	_, err = unix.Getfsstat(buf, noWait)
	if err != nil {
		return nil, err
	}
	stats := []filesystemStats{}
	for _, fs := range buf {
		mountpoint := gostring(fs.Mntonname[:])
		if c.ignoredMountPointsPattern.MatchString(mountpoint) {
			log.Debugf("Ignoring mount point: %s", mountpoint)
			continue
		}
		device := gostring(fs.Mntfromname[:])
		fstype := gostring(fs.Fstypename[:])
		if c.ignoredFSTypesPattern.MatchString(fstype) {
			log.Debugf("Ignoring fs type: %s", fstype)
			continue
		}
		var ro float64
		if (fs.Flags & readOnly) != 0 {
			ro = 1
		}
		stats = append(stats, filesystemStats{labels: filesystemLabels{device: device, mountPoint: mountpoint, fsType: fstype}, size: float64(fs.Blocks) * float64(fs.Bsize), free: float64(fs.Bfree) * float64(fs.Bsize), avail: float64(fs.Bavail) * float64(fs.Bsize), files: float64(fs.Files), filesFree: float64(fs.Ffree), ro: ro})
	}
	return stats, nil
}
