package collector

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func getLoad() (loads []float64, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	data, err := ioutil.ReadFile(procFilePath("loadavg"))
	if err != nil {
		return nil, err
	}
	loads, err = parseLoad(string(data))
	if err != nil {
		return nil, err
	}
	return loads, nil
}
func parseLoad(data string) (loads []float64, err error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	loads = make([]float64, 3)
	parts := strings.Fields(data)
	if len(parts) < 3 {
		return nil, fmt.Errorf("unexpected content in %s", procFilePath("loadavg"))
	}
	for i, load := range parts[0:3] {
		loads[i], err = strconv.ParseFloat(load, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse load '%s': %s", load, err)
		}
	}
	return loads, nil
}
