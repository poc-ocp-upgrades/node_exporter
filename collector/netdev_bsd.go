package collector

import (
	"errors"
	"regexp"
	"strconv"
	"github.com/prometheus/common/log"
)
import "C"

func getNetDevStats(ignore *regexp.Regexp) (map[string]map[string]string, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	netDev := map[string]map[string]string{}
	var ifap, ifa *C.struct_ifaddrs
	if C.getifaddrs(&ifap) == -1 {
		return nil, errors.New("getifaddrs() failed")
	}
	defer C.freeifaddrs(ifap)
	for ifa = ifap; ifa != nil; ifa = ifa.ifa_next {
		if ifa.ifa_addr.sa_family == C.AF_LINK {
			dev := C.GoString(ifa.ifa_name)
			if ignore.MatchString(dev) {
				log.Debugf("Ignoring device: %s", dev)
				continue
			}
			devStats := map[string]string{}
			data := (*C.struct_if_data)(ifa.ifa_data)
			devStats["receive_packets"] = convertFreeBSDCPUTime(uint64(data.ifi_ipackets))
			devStats["transmit_packets"] = convertFreeBSDCPUTime(uint64(data.ifi_opackets))
			devStats["receive_errs"] = convertFreeBSDCPUTime(uint64(data.ifi_ierrors))
			devStats["transmit_errs"] = convertFreeBSDCPUTime(uint64(data.ifi_oerrors))
			devStats["receive_bytes"] = convertFreeBSDCPUTime(uint64(data.ifi_ibytes))
			devStats["transmit_bytes"] = convertFreeBSDCPUTime(uint64(data.ifi_obytes))
			devStats["receive_multicast"] = convertFreeBSDCPUTime(uint64(data.ifi_imcasts))
			devStats["transmit_multicast"] = convertFreeBSDCPUTime(uint64(data.ifi_omcasts))
			devStats["receive_drop"] = convertFreeBSDCPUTime(uint64(data.ifi_iqdrops))
			devStats["transmit_drop"] = convertFreeBSDCPUTime(uint64(data.ifi_oqdrops))
			netDev[dev] = devStats
		}
	}
	return netDev, nil
}
func convertFreeBSDCPUTime(counter uint64) string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	return strconv.FormatUint(counter, 10)
}
