package edison

import "github.com/hybridgroup/gobot/sysfs"

type mux struct {
	pin   int
	value int
}
type sysfsPin struct {
	id           int // position, some mraa things refer to this
	name         string
	gpioPin      int // 0 is nil
	i2cPin       int // 0 is nil
	spiPin       int // 0 is nil
	resistor     int // 0 is nil
	levelShifter int
	pwmPin       int // -1 is nil
	mux          []mux
}

var sysfsPinsMiniboard = map[string]sysfsPin{
	"J17-1": sysfsPin{
		id:      0,
		name:    "J17-1",
		gpioPin: 182,
		pwmPin:  2,
	},
	"J17-2": sysfsPin{
		id:     1,
		name:   "J17-2",
		pwmPin: -1,
	},
	"J17-3": sysfsPin{
		id:     2,
		name:   "J17-3",
		pwmPin: -1,
	},
	"J17-4": sysfsPin{
		id:     3,
		name:   "J17-4",
		pwmPin: -1,
	},
	"J17-5": sysfsPin{
		id:      4,
		name:    "J17-5",
		gpioPin: 135,
		pwmPin:  -1,
	},
	"J17-6": sysfsPin{
		id:     5,
		name:   "J17-6",
		pwmPin: -1,
	},
	"J17-7": sysfsPin{
		id:      6,
		name:    "J17-7",
		gpioPin: 27,
		i2cPin:  1,
		pwmPin:  -1,
	},
	"J17-8": sysfsPin{
		id:      7,
		name:    "J17-8",
		gpioPin: 20,
		i2cPin:  1,
		pwmPin:  -1,
	},
	"J17-9": sysfsPin{
		id:      8,
		name:    "J17-9",
		gpioPin: 28,
		i2cPin:  1,
		pwmPin:  -1,
	},
	"J17-10": sysfsPin{
		id:      9,
		name:    "J17-10",
		gpioPin: 111,
		pwmPin:  -1,
	},
	"J17-11": sysfsPin{
		id:      10,
		name:    "J17-11",
		gpioPin: 109,
		spiPin:  5,
		pwmPin:  -1,
	},
	"J17-12": sysfsPin{
		id:      11,
		name:    "J17-12",
		gpioPin: 115,
		spiPin:  5,
		pwmPin:  -1,
	},
	"J17-13": sysfsPin{
		id:     12,
		name:   "J17-13",
		pwmPin: -1,
	},
	"J17-14": sysfsPin{
		id:      13,
		name:    "J17-14",
		gpioPin: 128,
		pwmPin:  -1,
	},
	"J18-1": sysfsPin{
		id:      14,
		name:    "J18-1",
		gpioPin: 13,
		pwmPin:  1,
	},
	"J18-2": sysfsPin{
		id:      15,
		name:    "J18-2",
		gpioPin: 165,
		pwmPin:  -1,
	},
	"J18-3": sysfsPin{
		id:     16,
		name:   "J18-3",
		pwmPin: -1,
	},
	"J18-4": sysfsPin{
		id:     17,
		name:   "J18-4",
		pwmPin: -1,
	},
	"J18-5": sysfsPin{
		id:     18,
		name:   "J18-5",
		pwmPin: -1,
	},
	"J18-6": sysfsPin{
		id:      19,
		name:    "J18-6",
		gpioPin: 19,
		i2cPin:  1,
		pwmPin:  -1,
	},
	"J18-7": sysfsPin{
		id:      20,
		name:    "J18-7",
		gpioPin: 12,
		pwmPin:  0,
	},
	"J18-8": sysfsPin{
		id:      21,
		name:    "J18-8",
		gpioPin: 183,
		pwmPin:  3,
	},
	"J18-9": sysfsPin{
		id:     22,
		name:   "J18-9",
		pwmPin: -1,
	},
	"J18-10": sysfsPin{
		id:      23,
		name:    "J18-10",
		gpioPin: 110,
		spiPin:  5,
		pwmPin:  -1,
	},
	"J18-11": sysfsPin{
		id:      24,
		name:    "J18-11",
		gpioPin: 114,
		spiPin:  5,
		pwmPin:  -1,
	},
	"J18-12": sysfsPin{
		id:      25,
		name:    "J18-12",
		gpioPin: 129,
		pwmPin:  -1,
	},
	"J18-13": sysfsPin{
		id:      26,
		name:    "J18-13",
		gpioPin: 130,
		// uart pinmap = 0
		pwmPin: -1,
	},
	"J18-14": sysfsPin{
		id:     27,
		name:   "J18-14",
		pwmPin: -1,
	},
	"J19-1": sysfsPin{
		id:     28,
		name:   "J19-1",
		pwmPin: -1,
	},
	"J19-2": sysfsPin{
		id:     29,
		name:   "J19-2",
		pwmPin: -1,
	},
	"J19-3": sysfsPin{
		id:     30,
		name:   "J19-3",
		pwmPin: -1,
	},
	"J19-4": sysfsPin{
		id:      31,
		name:    "J19-4",
		gpioPin: 44,
		pwmPin:  -1,
	},
	"J19-5": sysfsPin{
		id:      32,
		name:    "J19-5",
		gpioPin: 46,
		pwmPin:  -1,
	},
	"J19-6": sysfsPin{
		id:      33,
		name:    "J19-6",
		gpioPin: 48,
		pwmPin:  -1,
	},
	"J19-7": sysfsPin{
		id:     34,
		name:   "J19-7",
		pwmPin: -1,
	},
	"J19-8": sysfsPin{
		id:      35,
		name:    "J19-8",
		gpioPin: 131,
		// uart pinmap = 0
		pwmPin: -1,
	},
	"J19-9": sysfsPin{
		id:      36,
		name:    "J19-9",
		gpioPin: 14,
		pwmPin:  -1,
	},
	"J19-10": sysfsPin{
		id:      37,
		name:    "J19-10",
		gpioPin: 40,
		pwmPin:  -1,
	},
	"J19-11": sysfsPin{
		id:      38,
		name:    "J19-11",
		gpioPin: 43,
		pwmPin:  -1,
	},
	"J19-12": sysfsPin{
		id:      39,
		name:    "J19-12",
		gpioPin: 77,
		pwmPin:  -1,
	},
	"J19-13": sysfsPin{
		id:      40,
		name:    "J19-13",
		gpioPin: 82,
		pwmPin:  -1,
	},
	"J19-14": sysfsPin{
		id:      41,
		name:    "J19-14",
		gpioPin: 83,
		pwmPin:  -1,
	},
	"J20-1": sysfsPin{
		id:     42,
		name:   "J20-1",
		pwmPin: -1,
	},
	"J20-2": sysfsPin{
		id:     43,
		name:   "J20-2",
		pwmPin: -1,
	},
	"J20-3": sysfsPin{
		id:     44,
		name:   "J20-3",
		pwmPin: -1,
	},
	"J20-4": sysfsPin{
		id:      45,
		name:    "J20-4",
		gpioPin: 45,
		pwmPin:  -1,
	},
	"J20-5": sysfsPin{
		id:      46,
		name:    "J20-5",
		gpioPin: 47,
		pwmPin:  -1,
	},
	"J20-6": sysfsPin{
		id:      47,
		name:    "J20-6",
		gpioPin: 49,
		pwmPin:  -1,
	},
	"J20-7": sysfsPin{
		id:      48,
		name:    "J20-7",
		gpioPin: 15,
		pwmPin:  -1,
	},
	"J20-8": sysfsPin{
		id:      49,
		name:    "J20-8",
		gpioPin: 84,
		pwmPin:  -1,
	},
	"J20-9": sysfsPin{
		id:      50,
		name:    "J20-9",
		gpioPin: 42,
		pwmPin:  -1,
	},
	"J20-10": sysfsPin{
		id:      51,
		name:    "J20-10",
		gpioPin: 41,
		pwmPin:  -1,
	},
	"J20-11": sysfsPin{
		id:      52,
		name:    "J20-11",
		gpioPin: 78,
		pwmPin:  -1,
	},
	"J20-12": sysfsPin{
		id:      53,
		name:    "J20-12",
		gpioPin: 79,
		pwmPin:  -1,
	},
	"J20-13": sysfsPin{
		id:      54,
		name:    "J20-13",
		gpioPin: 80,
		pwmPin:  -1,
	},
	"J20-14": sysfsPin{
		id:      55,
		name:    "J20-14",
		gpioPin: 81,
		pwmPin:  -1,
	},
}

var sysfsPinsArduino = map[string]sysfsPin{
	"0": sysfsPin{
		name:         "IO0",
		gpioPin:      130,
		resistor:     216,
		levelShifter: 248,
		pwmPin:       -1,
	},
	"1": sysfsPin{
		name:         "IO1",
		gpioPin:      131,
		resistor:     217,
		levelShifter: 249,
		pwmPin:       -1,
	},
	"2": sysfsPin{
		name:         "IO2",
		gpioPin:      128,
		resistor:     218,
		levelShifter: 250,
		pwmPin:       -1,
	},
	"3": sysfsPin{
		name:         "IO3",
		gpioPin:      12,
		resistor:     219,
		levelShifter: 251,
		pwmPin:       0,
	},

	"4": sysfsPin{
		name:         "IO4",
		gpioPin:      129,
		resistor:     220,
		levelShifter: 252,
		pwmPin:       -1,
	},
	"5": sysfsPin{
		name:         "IO5",
		gpioPin:      13,
		resistor:     221,
		levelShifter: 253,
		pwmPin:       1,
	},
	"6": sysfsPin{
		name:         "IO6",
		gpioPin:      182,
		resistor:     222,
		levelShifter: 254,
		pwmPin:       2,
	},
	"7": sysfsPin{
		name:         "IO7",
		gpioPin:      48,
		resistor:     223,
		levelShifter: 255,
		pwmPin:       -1,
	},
	"8": sysfsPin{
		name:         "IO8",
		gpioPin:      49,
		resistor:     224,
		levelShifter: 256,
		pwmPin:       -1,
	},
	"9": sysfsPin{
		name:         "IO9",
		gpioPin:      183,
		resistor:     225,
		levelShifter: 257,
		pwmPin:       3,
	},
	"10": sysfsPin{
		name:         "IO10",
		gpioPin:      41,
		resistor:     226,
		levelShifter: 258,
		pwmPin:       4,
		mux: []mux{
			mux{263, sysfs.HIGH},
			mux{240, sysfs.LOW},
		},
	},
	"11": sysfsPin{
		name:         "IO11",
		gpioPin:      43,
		resistor:     227,
		levelShifter: 259,
		pwmPin:       5,
		mux: []mux{
			mux{262, sysfs.HIGH},
			mux{241, sysfs.LOW},
		},
	},
	"12": sysfsPin{
		name:         "IO12",
		gpioPin:      42,
		resistor:     228,
		levelShifter: 260,
		pwmPin:       -1,
		mux: []mux{
			mux{242, sysfs.LOW},
		},
	},
	"13": sysfsPin{
		name:         "IO13",
		gpioPin:      40,
		resistor:     229,
		levelShifter: 261,
		pwmPin:       -1,
		mux: []mux{
			mux{243, sysfs.LOW},
		},
	},
}
