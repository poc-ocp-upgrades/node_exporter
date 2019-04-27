package collector

import (
	"bufio"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"
	"github.com/prometheus/common/log"
)

const (
	defIgnoredMountPoints	= "^/(dev|proc|sys|var/lib/docker/.+)($|/)"
	defIgnoredFSTypes	= "^(autofs|binfmt_misc|bpf|cgroup2?|configfs|debugfs|devpts|devtmpfs|fusectl|hugetlbfs|mqueue|nsfs|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|selinuxfs|squashfs|sysfs|tracefs)$"
	readOnly		= 0x1
	mountTimeout		= 30 * time.Second
)

var stuckMounts = make(map[string]struct{})
var stuckMountsMtx = &sync.Mutex{}

func (c *filesystemCollector) GetStats() ([]filesystemStats, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	mps, err := mountPointDetails()
	if err != nil {
		return nil, err
	}
	stats := []filesystemStats{}
	for _, labels := range mps {
		if c.ignoredMountPointsPattern.MatchString(labels.mountPoint) {
			log.Debugf("Ignoring mount point: %s", labels.mountPoint)
			continue
		}
		if c.ignoredFSTypesPattern.MatchString(labels.fsType) {
			log.Debugf("Ignoring fs type: %s", labels.fsType)
			continue
		}
		stuckMountsMtx.Lock()
		if _, ok := stuckMounts[labels.mountPoint]; ok {
			stats = append(stats, filesystemStats{labels: labels, deviceError: 1})
			log.Debugf("Mount point %q is in an unresponsive state", labels.mountPoint)
			stuckMountsMtx.Unlock()
			continue
		}
		stuckMountsMtx.Unlock()
		success := make(chan struct{})
		go stuckMountWatcher(labels.mountPoint, success)
		buf := new(syscall.Statfs_t)
		err = syscall.Statfs(rootfsFilePath(labels.mountPoint), buf)
		stuckMountsMtx.Lock()
		close(success)
		if _, ok := stuckMounts[labels.mountPoint]; ok {
			log.Debugf("Mount point %q has recovered, monitoring will resume", labels.mountPoint)
			delete(stuckMounts, labels.mountPoint)
		}
		stuckMountsMtx.Unlock()
		if err != nil {
			stats = append(stats, filesystemStats{labels: labels, deviceError: 1})
			log.Debugf("Error on statfs() system call for %q: %s", rootfsFilePath(labels.mountPoint), err)
			continue
		}
		var ro float64
		for _, option := range strings.Split(labels.options, ",") {
			if option == "ro" {
				ro = 1
				break
			}
		}
		stats = append(stats, filesystemStats{labels: labels, size: float64(buf.Blocks) * float64(buf.Bsize), free: float64(buf.Bfree) * float64(buf.Bsize), avail: float64(buf.Bavail) * float64(buf.Bsize), files: float64(buf.Files), filesFree: float64(buf.Ffree), ro: ro})
	}
	return stats, nil
}
func stuckMountWatcher(mountPoint string, success chan struct{}) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	select {
	case <-success:
	case <-time.After(mountTimeout):
		stuckMountsMtx.Lock()
		select {
		case <-success:
		default:
			log.Debugf("Mount point %q timed out, it is being labeled as stuck and will not be monitored", mountPoint)
			stuckMounts[mountPoint] = struct{}{}
		}
		stuckMountsMtx.Unlock()
	}
}
func mountPointDetails() ([]filesystemLabels, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	file, err := os.Open(procFilePath("1/mounts"))
	if os.IsNotExist(err) {
		log.Debugf("Got %q reading root mounts, falling back to system mounts", err)
		file, err = os.Open(procFilePath("mounts"))
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()
	filesystems := []filesystemLabels{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		parts[1] = strings.Replace(parts[1], "\\040", " ", -1)
		parts[1] = strings.Replace(parts[1], "\\011", "\t", -1)
		filesystems = append(filesystems, filesystemLabels{device: parts[0], mountPoint: parts[1], fsType: parts[2], options: parts[3]})
	}
	return filesystems, scanner.Err()
}
