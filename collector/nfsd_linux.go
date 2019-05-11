package collector

import (
	"fmt"
	"os"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/nfs"
)

type nfsdCollector struct {
	fs				procfs.FS
	requestsDesc	*prometheus.Desc
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("nfsd", defaultEnabled, NewNFSdCollector)
}

const (
	nfsdSubsystem = "nfsd"
)

func NewNFSdCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %v", err)
	}
	return &nfsdCollector{fs: fs, requestsDesc: prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsdSubsystem, "requests_total"), "Total number NFSd Requests by method and protocol.", []string{"proto", "method"}, nil)}, nil
}
func (c *nfsdCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	stats, err := c.fs.NFSdServerRPCStats()
	if err != nil {
		if os.IsNotExist(err) {
			log.Debugf("Not collecting NFSd metrics: %s", err)
			return nil
		}
		return fmt.Errorf("failed to retrieve nfsd stats: %v", err)
	}
	c.updateNFSdReplyCacheStats(ch, &stats.ReplyCache)
	c.updateNFSdFileHandlesStats(ch, &stats.FileHandles)
	c.updateNFSdInputOutputStats(ch, &stats.InputOutput)
	c.updateNFSdThreadsStats(ch, &stats.Threads)
	c.updateNFSdReadAheadCacheStats(ch, &stats.ReadAheadCache)
	c.updateNFSdNetworkStats(ch, &stats.Network)
	c.updateNFSdServerRPCStats(ch, &stats.ServerRPC)
	c.updateNFSdRequestsv2Stats(ch, &stats.V2Stats)
	c.updateNFSdRequestsv3Stats(ch, &stats.V3Stats)
	c.updateNFSdRequestsv4Stats(ch, &stats.V4Ops)
	return nil
}
func (c *nfsdCollector) updateNFSdReplyCacheStats(ch chan<- prometheus.Metric, s *nfs.ReplyCache) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsdSubsystem, "reply_cache_hits_total"), "Total number of NFSd Reply Cache hits (client lost server response).", nil, nil), prometheus.CounterValue, float64(s.Hits))
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsdSubsystem, "reply_cache_misses_total"), "Total number of NFSd Reply Cache an operation that requires caching (idempotent).", nil, nil), prometheus.CounterValue, float64(s.Misses))
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsdSubsystem, "reply_cache_nocache_total"), "Total number of NFSd Reply Cache non-idempotent operations (rename/delete/…).", nil, nil), prometheus.CounterValue, float64(s.NoCache))
}
func (c *nfsdCollector) updateNFSdFileHandlesStats(ch chan<- prometheus.Metric, s *nfs.FileHandles) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsdSubsystem, "file_handles_stale_total"), "Total number of NFSd stale file handles", nil, nil), prometheus.CounterValue, float64(s.Stale))
}
func (c *nfsdCollector) updateNFSdInputOutputStats(ch chan<- prometheus.Metric, s *nfs.InputOutput) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsdSubsystem, "disk_bytes_read_total"), "Total NFSd bytes read.", nil, nil), prometheus.CounterValue, float64(s.Read))
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsdSubsystem, "disk_bytes_written_total"), "Total NFSd bytes written.", nil, nil), prometheus.CounterValue, float64(s.Write))
}
func (c *nfsdCollector) updateNFSdThreadsStats(ch chan<- prometheus.Metric, s *nfs.Threads) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsdSubsystem, "server_threads"), "Total number of NFSd kernel threads that are running.", nil, nil), prometheus.GaugeValue, float64(s.Threads))
}
func (c *nfsdCollector) updateNFSdReadAheadCacheStats(ch chan<- prometheus.Metric, s *nfs.ReadAheadCache) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsdSubsystem, "read_ahead_cache_size_blocks"), "How large the read ahead cache is in blocks.", nil, nil), prometheus.GaugeValue, float64(s.CacheSize))
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsdSubsystem, "read_ahead_cache_not_found_total"), "Total number of NFSd read ahead cache not found.", nil, nil), prometheus.CounterValue, float64(s.NotFound))
}
func (c *nfsdCollector) updateNFSdNetworkStats(ch chan<- prometheus.Metric, s *nfs.Network) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	packetDesc := prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsdSubsystem, "packets_total"), "Total NFSd network packets (sent+received) by protocol type.", []string{"proto"}, nil)
	ch <- prometheus.MustNewConstMetric(packetDesc, prometheus.CounterValue, float64(s.UDPCount), "udp")
	ch <- prometheus.MustNewConstMetric(packetDesc, prometheus.CounterValue, float64(s.TCPCount), "tcp")
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsdSubsystem, "connections_total"), "Total number of NFSd TCP connections.", nil, nil), prometheus.CounterValue, float64(s.TCPConnect))
}
func (c *nfsdCollector) updateNFSdServerRPCStats(ch chan<- prometheus.Metric, s *nfs.ServerRPC) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	badRPCDesc := prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsdSubsystem, "rpc_errors_total"), "Total number of NFSd RPC errors by error type.", []string{"error"}, nil)
	ch <- prometheus.MustNewConstMetric(badRPCDesc, prometheus.CounterValue, float64(s.BadFmt), "fmt")
	ch <- prometheus.MustNewConstMetric(badRPCDesc, prometheus.CounterValue, float64(s.BadAuth), "auth")
	ch <- prometheus.MustNewConstMetric(badRPCDesc, prometheus.CounterValue, float64(s.BadcInt), "cInt")
	ch <- prometheus.MustNewConstMetric(prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsdSubsystem, "server_rpcs_total"), "Total number of NFSd RPCs.", nil, nil), prometheus.CounterValue, float64(s.RPCCount))
}
func (c *nfsdCollector) updateNFSdRequestsv2Stats(ch chan<- prometheus.Metric, s *nfs.V2Stats) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	const proto = "2"
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.GetAttr), proto, "GetAttr")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.SetAttr), proto, "SetAttr")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Root), proto, "Root")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Lookup), proto, "Lookup")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.ReadLink), proto, "ReadLink")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Read), proto, "Read")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.WrCache), proto, "WrCache")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Write), proto, "Write")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Create), proto, "Create")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Remove), proto, "Remove")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Rename), proto, "Rename")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Link), proto, "Link")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.SymLink), proto, "SymLink")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.MkDir), proto, "MkDir")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.RmDir), proto, "RmDir")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.ReadDir), proto, "ReadDir")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.FsStat), proto, "FsStat")
}
func (c *nfsdCollector) updateNFSdRequestsv3Stats(ch chan<- prometheus.Metric, s *nfs.V3Stats) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	const proto = "3"
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.GetAttr), proto, "GetAttr")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.SetAttr), proto, "SetAttr")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Lookup), proto, "Lookup")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Access), proto, "Access")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.ReadLink), proto, "ReadLink")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Read), proto, "Read")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Write), proto, "Write")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Create), proto, "Create")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.MkDir), proto, "MkDir")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.SymLink), proto, "SymLink")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.MkNod), proto, "MkNod")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Remove), proto, "Remove")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.RmDir), proto, "RmDir")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Rename), proto, "Rename")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Link), proto, "Link")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.ReadDir), proto, "ReadDir")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.ReadDirPlus), proto, "ReadDirPlus")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.FsStat), proto, "FsStat")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.FsInfo), proto, "FsInfo")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.PathConf), proto, "PathConf")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Commit), proto, "Commit")
}
func (c *nfsdCollector) updateNFSdRequestsv4Stats(ch chan<- prometheus.Metric, s *nfs.V4Ops) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	const proto = "4"
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Access), proto, "Access")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Close), proto, "Close")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Commit), proto, "Commit")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Create), proto, "Create")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.DelegPurge), proto, "DelegPurge")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.DelegReturn), proto, "DelegReturn")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.GetAttr), proto, "GetAttr")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.GetFH), proto, "GetFH")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Link), proto, "Link")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Lock), proto, "Lock")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Lockt), proto, "Lockt")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Locku), proto, "Locku")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Lookup), proto, "Lookup")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.LookupRoot), proto, "LookupRoot")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Nverify), proto, "Nverify")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Open), proto, "Open")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.OpenAttr), proto, "OpenAttr")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.OpenConfirm), proto, "OpenConfirm")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.OpenDgrd), proto, "OpenDgrd")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.PutFH), proto, "PutFH")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Read), proto, "Read")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.ReadDir), proto, "ReadDir")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.ReadLink), proto, "ReadLink")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Remove), proto, "Remove")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Rename), proto, "Rename")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Renew), proto, "Renew")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.RestoreFH), proto, "RestoreFH")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.SaveFH), proto, "SaveFH")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.SecInfo), proto, "SecInfo")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.SetAttr), proto, "SetAttr")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Verify), proto, "Verify")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.Write), proto, "Write")
	ch <- prometheus.MustNewConstMetric(c.requestsDesc, prometheus.CounterValue, float64(s.RelLockOwner), proto, "RelLockOwner")
}
