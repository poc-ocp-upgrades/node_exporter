package collector

import (
	"fmt"
	"os"
	"reflect"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/nfs"
)

const (
	nfsSubsystem = "nfs"
)

type nfsCollector struct {
	fs					procfs.FS
	nfsNetReadsDesc				*prometheus.Desc
	nfsNetConnectionsDesc			*prometheus.Desc
	nfsRPCOperationsDesc			*prometheus.Desc
	nfsRPCRetransmissionsDesc		*prometheus.Desc
	nfsRPCAuthenticationRefreshesDesc	*prometheus.Desc
	nfsProceduresDesc			*prometheus.Desc
}

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("nfs", defaultEnabled, NewNfsCollector)
}
func NewNfsCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	fs, err := procfs.NewFS(*procPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open procfs: %v", err)
	}
	return &nfsCollector{fs: fs, nfsNetReadsDesc: prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsSubsystem, "packets_total"), "Total NFSd network packets (sent+received) by protocol type.", []string{"protocol"}, nil), nfsNetConnectionsDesc: prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsSubsystem, "connections_total"), "Total number of NFSd TCP connections.", nil, nil), nfsRPCOperationsDesc: prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsSubsystem, "rpcs_total"), "Total number of RPCs performed.", nil, nil), nfsRPCRetransmissionsDesc: prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsSubsystem, "rpc_retransmissions_total"), "Number of RPC transmissions performed.", nil, nil), nfsRPCAuthenticationRefreshesDesc: prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsSubsystem, "rpc_authentication_refreshes_total"), "Number of RPC authentication refreshes performed.", nil, nil), nfsProceduresDesc: prometheus.NewDesc(prometheus.BuildFQName(namespace, nfsSubsystem, "requests_total"), "Number of NFS procedures invoked.", []string{"proto", "method"}, nil)}, nil
}
func (c *nfsCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	stats, err := c.fs.NFSClientRPCStats()
	if err != nil {
		if os.IsNotExist(err) {
			log.Debugf("Not collecting NFS metrics: %s", err)
			return nil
		}
		return fmt.Errorf("failed to retrieve nfs stats: %v", err)
	}
	c.updateNFSNetworkStats(ch, &stats.Network)
	c.updateNFSClientRPCStats(ch, &stats.ClientRPC)
	c.updateNFSRequestsv2Stats(ch, &stats.V2Stats)
	c.updateNFSRequestsv3Stats(ch, &stats.V3Stats)
	c.updateNFSRequestsv4Stats(ch, &stats.ClientV4Stats)
	return nil
}
func (c *nfsCollector) updateNFSNetworkStats(ch chan<- prometheus.Metric, s *nfs.Network) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ch <- prometheus.MustNewConstMetric(c.nfsNetReadsDesc, prometheus.CounterValue, float64(s.UDPCount), "udp")
	ch <- prometheus.MustNewConstMetric(c.nfsNetReadsDesc, prometheus.CounterValue, float64(s.TCPCount), "tcp")
	ch <- prometheus.MustNewConstMetric(c.nfsNetConnectionsDesc, prometheus.CounterValue, float64(s.TCPConnect))
}
func (c *nfsCollector) updateNFSClientRPCStats(ch chan<- prometheus.Metric, s *nfs.ClientRPC) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ch <- prometheus.MustNewConstMetric(c.nfsRPCOperationsDesc, prometheus.CounterValue, float64(s.RPCCount))
	ch <- prometheus.MustNewConstMetric(c.nfsRPCRetransmissionsDesc, prometheus.CounterValue, float64(s.Retransmissions))
	ch <- prometheus.MustNewConstMetric(c.nfsRPCAuthenticationRefreshesDesc, prometheus.CounterValue, float64(s.AuthRefreshes))
}
func (c *nfsCollector) updateNFSRequestsv2Stats(ch chan<- prometheus.Metric, s *nfs.V2Stats) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	const proto = "2"
	v := reflect.ValueOf(s).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue, float64(field.Uint()), proto, v.Type().Field(i).Name)
	}
}
func (c *nfsCollector) updateNFSRequestsv3Stats(ch chan<- prometheus.Metric, s *nfs.V3Stats) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	const proto = "3"
	v := reflect.ValueOf(s).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue, float64(field.Uint()), proto, v.Type().Field(i).Name)
	}
}
func (c *nfsCollector) updateNFSRequestsv4Stats(ch chan<- prometheus.Metric, s *nfs.ClientV4Stats) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	const proto = "4"
	v := reflect.ValueOf(s).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		ch <- prometheus.MustNewConstMetric(c.nfsProceduresDesc, prometheus.CounterValue, float64(field.Uint()), proto, v.Type().Field(i).Name)
	}
}
