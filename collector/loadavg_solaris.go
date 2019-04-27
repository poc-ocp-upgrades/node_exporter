package collector

import (
	"errors"
)
import "C"

func getLoad() ([]float64, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	var loadavg [3]C.double
	samples := C.getloadavg(&loadavg[0], 3)
	if samples != 3 {
		return nil, errors.New("failed to get load average")
	}
	return []float64{float64(loadavg[0]), float64(loadavg[1]), float64(loadavg[2])}, nil
}
