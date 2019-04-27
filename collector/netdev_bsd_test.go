package collector

import "testing"

type uintToStringTest struct {
	in	uint64
	out	string
}

var uinttostringtests = []uintToStringTest{{0, "0"}, {1, "1"}, {12345678, "12345678"}, {1<<31 - 1, "2147483647"}, {1 << 31, "2147483648"}, {1<<31 + 1, "2147483649"}, {1<<32 - 1, "4294967295"}, {1 << 32, "4294967296"}, {1<<32 + 1, "4294967297"}, {1 << 50, "1125899906842624"}, {1<<63 - 1, "9223372036854775807"}, {0x1bf0c640a, "7500227594"}, {0xbee5df75, "3202735989"}}

func TestUintToString(t *testing.T) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	for _, test := range uinttostringtests {
		is := convertFreeBSDCPUTime(test.in)
		if is != test.out {
			t.Errorf("convertFreeBSDCPUTime(%v) = %v want %v", test.in, is, test.out)
		}
	}
}
