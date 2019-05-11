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
			devStats["receive_packets"] = strconv.FormatUint(uint64(data.ifi_ipackets), 10)
			devStats["transmit_packets"] = strconv.FormatUint(uint64(data.ifi_opackets), 10)
			devStats["receive_errs"] = strconv.FormatUint(uint64(data.ifi_ierrors), 10)
			devStats["transmit_errs"] = strconv.FormatUint(uint64(data.ifi_oerrors), 10)
			devStats["receive_bytes"] = strconv.FormatUint(uint64(data.ifi_ibytes), 10)
			devStats["transmit_bytes"] = strconv.FormatUint(uint64(data.ifi_obytes), 10)
			devStats["receive_multicast"] = strconv.FormatUint(uint64(data.ifi_imcasts), 10)
			devStats["transmit_multicast"] = strconv.FormatUint(uint64(data.ifi_omcasts), 10)
			netDev[dev] = devStats
		}
	}
	return netDev, nil
}
