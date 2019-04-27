package collector

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/procfs/sysfs"
	"github.com/prometheus/procfs/xfs"
)

type xfsCollector struct{ fs sysfs.FS }

func init() {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	registerCollector("xfs", defaultEnabled, NewXFSCollector)
}
func NewXFSCollector() (Collector, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	fs, err := sysfs.NewFS(*sysPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sysfs: %v", err)
	}
	return &xfsCollector{fs: fs}, nil
}
func (c *xfsCollector) Update(ch chan<- prometheus.Metric) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	stats, err := c.fs.XFSStats()
	if err != nil {
		return fmt.Errorf("failed to retrieve XFS stats: %v", err)
	}
	for _, s := range stats {
		c.updateXFSStats(ch, s)
	}
	return nil
}
func (c *xfsCollector) updateXFSStats(ch chan<- prometheus.Metric, s *xfs.Stats) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	const (
		subsystem = "xfs"
	)
	var (
		labels = []string{"device"}
	)
	metrics := []struct {
		name	string
		desc	string
		value	float64
	}{{name: "extent_allocation_extents_allocated_total", desc: "Number of extents allocated for a filesystem.", value: float64(s.ExtentAllocation.ExtentsAllocated)}, {name: "extent_allocation_blocks_allocated_total", desc: "Number of blocks allocated for a filesystem.", value: float64(s.ExtentAllocation.BlocksAllocated)}, {name: "extent_allocation_extents_freed_total", desc: "Number of extents freed for a filesystem.", value: float64(s.ExtentAllocation.ExtentsFreed)}, {name: "extent_allocation_blocks_freed_total", desc: "Number of blocks freed for a filesystem.", value: float64(s.ExtentAllocation.BlocksFreed)}, {name: "allocation_btree_lookups_total", desc: "Number of allocation B-tree lookups for a filesystem.", value: float64(s.AllocationBTree.Lookups)}, {name: "allocation_btree_compares_total", desc: "Number of allocation B-tree compares for a filesystem.", value: float64(s.AllocationBTree.Compares)}, {name: "allocation_btree_records_inserted_total", desc: "Number of allocation B-tree records inserted for a filesystem.", value: float64(s.AllocationBTree.RecordsInserted)}, {name: "allocation_btree_records_deleted_total", desc: "Number of allocation B-tree records deleted for a filesystem.", value: float64(s.AllocationBTree.RecordsDeleted)}, {name: "block_mapping_reads_total", desc: "Number of block map for read operations for a filesystem.", value: float64(s.BlockMapping.Reads)}, {name: "block_mapping_writes_total", desc: "Number of block map for write operations for a filesystem.", value: float64(s.BlockMapping.Writes)}, {name: "block_mapping_unmaps_total", desc: "Number of block unmaps (deletes) for a filesystem.", value: float64(s.BlockMapping.Unmaps)}, {name: "block_mapping_extent_list_insertions_total", desc: "Number of extent list insertions for a filesystem.", value: float64(s.BlockMapping.ExtentListInsertions)}, {name: "block_mapping_extent_list_deletions_total", desc: "Number of extent list deletions for a filesystem.", value: float64(s.BlockMapping.ExtentListDeletions)}, {name: "block_mapping_extent_list_lookups_total", desc: "Number of extent list lookups for a filesystem.", value: float64(s.BlockMapping.ExtentListLookups)}, {name: "block_mapping_extent_list_compares_total", desc: "Number of extent list compares for a filesystem.", value: float64(s.BlockMapping.ExtentListCompares)}, {name: "block_map_btree_lookups_total", desc: "Number of block map B-tree lookups for a filesystem.", value: float64(s.BlockMapBTree.Lookups)}, {name: "block_map_btree_compares_total", desc: "Number of block map B-tree compares for a filesystem.", value: float64(s.BlockMapBTree.Compares)}, {name: "block_map_btree_records_inserted_total", desc: "Number of block map B-tree records inserted for a filesystem.", value: float64(s.BlockMapBTree.RecordsInserted)}, {name: "block_map_btree_records_deleted_total", desc: "Number of block map B-tree records deleted for a filesystem.", value: float64(s.BlockMapBTree.RecordsDeleted)}}
	for _, m := range metrics {
		desc := prometheus.NewDesc(prometheus.BuildFQName(namespace, subsystem, m.name), m.desc, labels, nil)
		ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, m.value, s.Name)
	}
}
