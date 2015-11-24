package i2c

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*LSM9DS0Driver)(nil)

const LSM9DS0_ADDRESS_XM = 0x1d
const LSM9DS0_ADDRESS_G = 0x6b

// See http://www.adafruit.com/datasheets/LSM9DS0.pdf for details.

const LSM9DS0_OUT_TEMP_L = 0x05
const LSM9DS0_OUT_TEMP_H = 0x06
const LSM9DS0_STATUS_REG_M = 0x07
const LSM9DS0_OUT_X_L_M = 0x08
const LSM9DS0_OUT_X_H_M = 0x09
const LSM9DS0_OUT_Y_L_M = 0x0A
const LSM9DS0_OUT_Y_H_M = 0x0B
const LSM9DS0_OUT_Z_L_M = 0x0C
const LSM9DS0_OUT_Z_H_M = 0x0D
const LSM9DS0_OUT_X_L_A = 0x28
const LSM9DS0_OUT_X_H_A = 0x29
const LSM9DS0_OUT_Y_L_A = 0x2A
const LSM9DS0_OUT_Y_H_A = 0x2B
const LSM9DS0_OUT_Z_L_A = 0x2C
const LSM9DS0_OUT_Z_H_A = 0x2D
const LSM9DS0_OUT_X_L_G = 0x28
const LSM9DS0_OUT_X_H_G = 0x29
const LSM9DS0_OUT_Y_L_G = 0x2A
const LSM9DS0_OUT_Y_H_G = 0x2B
const LSM9DS0_OUT_Z_L_G = 0x2C
const LSM9DS0_OUT_Z_H_G = 0x2D
const LSM9DS0_WHO_AM_I_XMG = 0x0F
const LSM9DS0_INT_CTRL_REG_M = 0x12
const LSM9DS0_INT_SRC_REG_M = 0x13
const LSM9DS0_INT_THS_L_M = 0x14
const LSM9DS0_INT_THS_H_M = 0x15
const LSM9DS0_OFFSET_X_L_M = 0x16
const LSM9DS0_OFFSET_X_H_M = 0x17
const LSM9DS0_OFFSET_Y_L_M = 0x18
const LSM9DS0_OFFSET_Y_H_M = 0x19
const LSM9DS0_OFFSET_Z_L_M = 0x1A
const LSM9DS0_OFFSET_Z_H_M = 0x1B
const LSM9DS0_REFERENCE_X = 0x1C
const LSM9DS0_REFERENCE_Y = 0x1D
const LSM9DS0_REFERENCE_Z = 0x1E
const LSM9DS0_CTRL_REG1_G = 0x20
const LSM9DS0_CTRL_REG2_G = 0x21
const LSM9DS0_CTRL_REG3_G = 0x22
const LSM9DS0_CTRL_REG4_G = 0x23
const LSM9DS0_CTRL_REG5_G = 0x24
const LSM9DS0_CTRL_REG0_XM = 0x1F
const LSM9DS0_CTRL_REG1_XM = 0x20
const LSM9DS0_CTRL_REG2_XM = 0x21
const LSM9DS0_CTRL_REG3_XM = 0x22
const LSM9DS0_CTRL_REG4_XM = 0x23
const LSM9DS0_CTRL_REG5_XM = 0x24
const LSM9DS0_CTRL_REG6_XM = 0x25
const LSM9DS0_CTRL_REG7 = 0x26
const LSM9DS0_REFERENCE_G = 0x25
const LSM9DS0_STATUS_REG_AG = 0x27
const LSM9DS0_FIFO_CTRL_REG = 0x2E
const LSM9DS0_FIFO_SRC_REG = 0x2F
const LSM9DS0_INT_GEN_1_REG = 0x30
const LSM9DS0_INT_GEN_1_SRC = 0x31
const LSM9DS0_INT_GEN_1_THS = 0x32
const LSM9DS0_INT_GEN_1_DURATION = 0x33
const LSM9DS0_INT_GEN_2_REG = 0x34
const LSM9DS0_INT_GEN_2_SRC = 0x35
const LSM9DS0_INT_GEN_2_THS = 0x36
const LSM9DS0_INT_GEN_2_DURATION = 0x37
const LSM9DS0_INT1_CFG_G = 0x30
const LSM9DS0_INT1_SRC_G = 0x31
const LSM9DS0_INT1_THS_XH_G = 0x32
const LSM9DS0_INT1_THS_XL_G = 0x33
const LSM9DS0_INT1_THS_YH_G = 0x34
const LSM9DS0_INT1_THS_YL_G = 0x35
const LSM9DS0_INT1_THS_ZH_G = 0x36
const LSM9DS0_INT1_THS_ZL_G = 0x37
const LSM9DS0_INT1_DURATION_G = 0x38
const LSM9DS0_CLICK_CFG = 0x38
const LSM9DS0_CLICK_SRC = 0x39
const LSM9DS0_CLICK_THS = 0x3A
const LSM9DS0_TIME_LIMIT = 0x3B
const LSM9DS0_TIME_LATENCY = 0x3C
const LSM9DS0_TIME_WINDOW = 0x3D
const LSM9DS0_ACT_THS = 0x3E
const LSM9DS0_ACT_DUR = 0x3F
const LSM9DS0_G_SCALE_245DPS = 0  // 245 degrees per second
const LSM9DS0_G_SCALE_500DPS = 1  // 500 dps
const LSM9DS0_G_SCALE_2000DPS = 2 // 2000 dps

const LSM9DS0_G_ODR_95_BW_125 = 0x0 //   95         12.5
const LSM9DS0_G_ODR_95_BW_25 = 0x1  //   95          25

// 0x2 and 0x3 define the same data rate and bandwidth
const LSM9DS0_G_ODR_190_BW_125 = 0x4 //   190        12.5
const LSM9DS0_G_ODR_190_BW_25 = 0x5  //   190         25
const LSM9DS0_G_ODR_190_BW_50 = 0x6  //   190         50
const LSM9DS0_G_ODR_190_BW_70 = 0x7  //   190         70
const LSM9DS0_G_ODR_380_BW_20 = 0x8  //   380         20
const LSM9DS0_G_ODR_380_BW_25 = 0x9  //   380         25
const LSM9DS0_G_ODR_380_BW_50 = 0xA  //   380         50
const LSM9DS0_G_ODR_380_BW_100 = 0xB //   380         100
const LSM9DS0_G_ODR_760_BW_30 = 0xC  //   760         30
const LSM9DS0_G_ODR_760_BW_35 = 0xD  //   760         35
const LSM9DS0_G_ODR_760_BW_50 = 0xE  //   760         50
const LSM9DS0_G_ODR_760_BW_100 = 0xF //   760         100

type LSM9DS0Driver struct {
	name       string
	connection I2c

	accelEnabled   bool
	gyroEnabled    bool
	magnetoEnabled bool

	gyroScale  string
	gyroScales map[string]lsm9ds0GyroScale

	accelScale  string
	accelScales map[string]lsm9ds0AccelScale
}

type lsm9ds0GyroScale struct {
	sensitivity float64 // millidegrees/second/LSB - least significant bit
	fs1fs0      uint8
}

type lsm9ds0AccelScale struct {
	factor float64
	afs  uint8
}

// NewLSM9DS0Driver creates a new driver with specified name and i2c interface, enabling
// certain features or not.
func NewLSM9DS0Driver(a I2c, name string, accel, gyro, magneto bool) *LSM9DS0Driver {
	return &LSM9DS0Driver{
		name:           name,
		connection:     a,
		accelEnabled:   accel,
		gyroEnabled:    gyro,
		magnetoEnabled: magneto,
		gyroScale:      "245",
		gyroScales: map[string]lsm9ds0GyroScale{
			"245":  {sensitivity: 8.75, fs1fs0: 0x00},
			"500":  {sensitivity: 17.50, fs1fs0: 0x01},
			"2000": {sensitivity: 70, fs1fs0: 0x02},
		},
		accelScale: "2",
		accelScales: map[string]lsm9ds0AccelScale{
			"2":  {factor: 0.061, afs: 0},
			"4":  {factor: 0.122, afs: 1},
			"6":  {factor: 0.183, afs: 2},
			"8":  {factor: 0.244, afs: 3},
			"16": {factor: 0.732, afs: 4},
		},
	}
}

func (l *LSM9DS0Driver) Name() string                 { return l.name }
func (l *LSM9DS0Driver) Connection() gobot.Connection { return l.connection.(gobot.Connection) }

// Start initialized the lsm9ds0
func (l *LSM9DS0Driver) Start() (errs []error) {
	if l.gyroEnabled {
		if err := l.connection.I2cStart(LSM9DS0_ADDRESS_G); err != nil {
			return []error{err}
		}
	}
	if l.accelEnabled || l.magnetoEnabled {
		if err := l.connection.I2cStart(LSM9DS0_ADDRESS_XM); err != nil {
			return []error{err}
		}
	}

	if l.gyroEnabled {
		err := l.initGyroscope()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if l.accelEnabled {
		err := l.initAccelerator()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if l.magnetoEnabled {
		err := l.initMagnetometer()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return
}

// Halt returns true if devices is halted successfully
func (l *LSM9DS0Driver) Halt() (errs []error) { return }

type Vector struct {
	X float64
	Y float64
	Z float64
}

// GetGyro returns a vector, representing angular velocity in degrees/sec.
func (l *LSM9DS0Driver) GetGyro() (av Vector, err error) {
	if !l.gyroEnabled {
		return av, fmt.Errorf("not enabled")
	}

	i2cCtx := l.connection.(I2cReaderBytesData)
	data := make([]byte, 6)
	if err = i2cCtx.I2cReadBytesData(LSM9DS0_ADDRESS_G, byte(LSM9DS0_OUT_X_L_G)|0x80, data); err != nil {
		return av, fmt.Errorf("I2cReadBytesData error: %s", err)
	}
	// if err = l.connection.I2cWrite(LSM9DS0_ADDRESS_G, []byte{LSM9DS0_OUT_X_L_G}); err != nil {
	// }
	// data, err := l.connection.I2cRead(LSM9DS0_ADDRESS_G, 6)
	// if err != nil {
	// 	return av, fmt.Errorf("I2cRead error: %s", err)
	// }

	buf := bytes.NewBuffer(data)
	load := struct {
		X, Y, Z int16
	}{}
	err = binary.Read(buf, binary.LittleEndian, &load)
	//fmt.Printf("loaded int16: %#v binary=%x\n", load, data)

	scale := l.gyroScales[l.gyroScale].sensitivity // millisecond/second/LSB
	av.X = float64(load.X) * scale / 1000.0
	av.Y = float64(load.Y) * scale / 1000.0
	av.Z = float64(load.Z) * scale / 1000.0

	return
}

func (l *LSM9DS0Driver) GetAccel() (v Vector, err error) {
	if !l.accelEnabled {
		return v, fmt.Errorf("not enabled")
	}

	i2cCtx := l.connection.(I2cReaderBytesData)
	data := make([]byte, 6)
	if err = i2cCtx.I2cReadBytesData(LSM9DS0_ADDRESS_XM, byte(LSM9DS0_OUT_X_L_A)|0x80, data); err != nil {
		return v, fmt.Errorf("I2cReadBytesData error: %s", err)
	}

	// if err = l.connection.I2cWrite(LSM9DS0_ADDRESS_XM, []byte{LSM9DS0_OUT_X_L_A}); err != nil {
	// 	return p, fmt.Errorf("I2cWrite error: %s", err)
	// }
	// data, err := l.connection.I2cRead(LSM9DS0_ADDRESS_G, 6)
	// if err != nil {
	// 	return p, fmt.Errorf("I2cRead error: %s", err)
	// }

	buf := bytes.NewBuffer(data)
	load := struct {
		X, Y, Z int16
	}{}
	err = binary.Read(buf, binary.LittleEndian, &load)

	scale := 1.0
	v.X = float64(load.X) * scale / 1000.0
	v.Y = float64(load.Y) * scale / 1000.0
	v.Z = float64(load.Z) * scale / 1000.0

	return
}

// SetAccelScale sets the sensitivity. Allows values: 245, 500 and 2000 dps
func (l *LSM9DS0Driver) SetAccelScale(scale string) error {
	if _, ok := l.accelScales[scale]; !ok {
		return errors.New("invalid value, use 2, 4, 6, 8 or 16")
	}
	l.accelScale = scale
	return l.setAccelScale()
}

func (l *LSM9DS0Driver) setAccelScale() (err error) {
	afs := l.accelScales[l.accelScale].afs
	err = l.connection.I2cWrite(LSM9DS0_ADDRESS_XM, []byte{LSM9DS0_CTRL_REG2_XM, afs << 3})
	return
}

func (l *LSM9DS0Driver) initAccelerator() (err error) {
	// See https://github.com/hybridgroup/cylon-i2c/blob/master/lib/lsm9ds0xm.js

	if err = l.connection.I2cWrite(LSM9DS0_ADDRESS_XM, []byte{LSM9DS0_CTRL_REG0_XM, 0x00}); err != nil {
		return
	}

	if err = l.connection.I2cWrite(LSM9DS0_ADDRESS_XM, []byte{LSM9DS0_CTRL_REG1_XM, 0x57}); err != nil {
		return
	}

	if err = l.setAccelScale(); err != nil {
		return
	}

	// Accelerometer data ready on INT1_XM (0x04)
	if err = l.connection.I2cWrite(LSM9DS0_ADDRESS_XM, []byte{LSM9DS0_CTRL_REG3_XM, 0x04}); err != nil {
		return
	}

	return nil
}

func (l *LSM9DS0Driver) initMagnetometer() error {
	// See https://github.com/hybridgroup/cylon-i2c/blob/master/lib/lsm9ds0xm.js
	return nil
}

// SetGyroScale sets the sensitivity. Allows values: 245, 500 and 2000 dps
func (l *LSM9DS0Driver) SetGyroScale(scale string) error {
	if _, ok := l.gyroScales[scale]; !ok {
		return errors.New("invalid value, use 245, 500 or 2000")
	}
	l.gyroScale = scale
	return l.setGyroScale()
}

func (l *LSM9DS0Driver) setGyroScale() (err error) {
	fs1fs0 := l.gyroScales[l.gyroScale].fs1fs0
	err = l.connection.I2cWrite(LSM9DS0_ADDRESS_G, []byte{LSM9DS0_CTRL_REG4_G, fs1fs0 << 4})
	return
}

func (l *LSM9DS0Driver) initGyroscope() (err error) {
	// See https://github.com/hybridgroup/cylon-i2c/blob/master/lib/lsm9ds0g.js

	// ODR 95Hz, Cutoff 25 + Normal mode + Enable all axes
	if err = l.connection.I2cWrite(LSM9DS0_ADDRESS_G, []byte{LSM9DS0_CTRL_REG1_G, 0x1F}); err != nil {
		return
	}

	// Normal mode, high cutoff frequency
	if err = l.connection.I2cWrite(LSM9DS0_ADDRESS_G, []byte{LSM9DS0_CTRL_REG2_G, 0x00}); err != nil {
		return
	}

	// Int1 enabled (pp, active low), data read on DRDY_G:
	if err = l.connection.I2cWrite(LSM9DS0_ADDRESS_G, []byte{LSM9DS0_CTRL_REG3_G, 0x88}); err != nil {
		return
	}

	if err = l.setGyroScale(); err != nil {
		return
	}

	// if err = l.connection.I2cWrite(LSM9DS0_ADDRESS_G, []byte{LSM9DS0_CTRL_REG5_G, 0x00}); err != nil {
	// 	return
	// }

	return
}
