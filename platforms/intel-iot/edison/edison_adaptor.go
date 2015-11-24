package edison

import (
	"errors"
	"os"
	"strconv"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/i2c"
	"github.com/hybridgroup/gobot/sysfs"
)

var _ gobot.Adaptor = (*EdisonAdaptor)(nil)

var _ gpio.DigitalReader = (*EdisonAdaptor)(nil)
var _ gpio.DigitalWriter = (*EdisonAdaptor)(nil)
var _ gpio.AnalogReader = (*EdisonAdaptor)(nil)
var _ gpio.PwmWriter = (*EdisonAdaptor)(nil)

var _ i2c.I2c = (*EdisonAdaptor)(nil)

func writeFile(path string, data []byte) (i int, err error) {
	file, err := sysfs.OpenFile(path, os.O_WRONLY, 0644)
	defer file.Close()
	if err != nil {
		return
	}

	return file.Write(data)
}

func readFile(path string) ([]byte, error) {
	file, err := sysfs.OpenFile(path, os.O_RDONLY, 0644)
	defer file.Close()
	if err != nil {
		return make([]byte, 0), err
	}

	buf := make([]byte, 200)
	var i = 0
	i, err = file.Read(buf)
	if i == 0 {
		return buf, err
	}
	return buf[:i], err
}

// EdisonAdaptor represents an Intel Edison
type EdisonAdaptor struct {
	name        string
	tristate    sysfs.DigitalPin
	digitalPins map[int]sysfs.DigitalPin
	sysfsPinMap map[string]sysfsPin
	pwmPins     map[int]*pwmPin
	i2cDevice   sysfs.I2cDevice
	connect     func(e *EdisonAdaptor) (err error)
	miniboard   bool
}

// changePinMode writes pin mode to current_pinmux file
func changePinMode(gpioPin, mode string) (err error) {
	_, err = writeFile(
		"/sys/kernel/debug/gpio_debug/gpio"+gpioPin+"/current_pinmux",
		[]byte("mode"+mode),
	)
	return
}

// NewEdisonAdaptor returns a new EdisonAdaptor with specified name
func NewEdisonAdaptor(name string) *EdisonAdaptor {
	return &EdisonAdaptor{
		name: name,
		//i2cDevices: make(map[int]io.ReadWriteCloser),
		//i2cDevices: make(map[int]io.ReadWriteCloser),
		connect: func(e *EdisonAdaptor) (err error) {
			e.tristate = sysfs.NewDigitalPin(214)
			if err = e.tristate.Export(); err != nil {
				// edison: Failed to initialize Arduino board TriState
				// assuming miniboard:
				e.miniboard = true
				e.sysfsPinMap = sysfsPinsMiniboard
				return e.connectMiniboard()
			} else {
				e.sysfsPinMap = sysfsPinsArduino
				return e.connectArduinoBoard()
			}

		},
	}
}

func (e *EdisonAdaptor) connectMiniboard() (err error) {
	// TODO: replace the sysfsPinMap with the miniboard pins..
	// TODO: do any initialization of the miniboard, following mraa

	/*
	   b->adv_func->gpio_init_post = &mraa_intel_edison_gpio_init_post;
	   b->adv_func->pwm_init_pre = &mraa_intel_edison_pwm_init_pre;
	   b->adv_func->i2c_init_pre = &mraa_intel_edison_i2c_init_pre;
	   b->adv_func->i2c_set_frequency_replace = &mraa_intel_edison_i2c_freq;
	   b->adv_func->spi_init_pre = &mraa_intel_edison_spi_init_pre;
	   b->adv_func->gpio_mode_replace = &mraa_intel_edsion_mb_gpio_mode;
	   b->adv_func->uart_init_pre = &mraa_intel_edison_uart_init_pre;
	   b->adv_func->gpio_mmap_setup = &mraa_intel_edison_mmap_setup;
	*/

	// Miniboard J17-10, J17-11, J17-12, J18-11,
	for _, i := range []string{"111", "109", "115", "114"} {
		if err = changePinMode(i, "1"); err != nil {
			return err
		}
	}

	return nil
}
func (e *EdisonAdaptor) connectArduinoBoard() (err error) {
	if err = e.tristate.Direction(sysfs.OUT); err != nil {
		return err
	}
	if err = e.tristate.Write(sysfs.LOW); err != nil {
		return err
	}

	// Setup mux HIGH
	for _, i := range []int{263, 262} {
		io := sysfs.NewDigitalPin(i)
		if err = io.Export(); err != nil {
			return err
		}
		if err = io.Direction(sysfs.OUT); err != nil {
			return err
		}
		if err = io.Write(sysfs.HIGH); err != nil {
			return err
		}
		if err = io.Unexport(); err != nil {
			return err
		}
	}

	// Setup mux LOW
	for _, i := range []int{240, 241, 242, 243} {
		io := sysfs.NewDigitalPin(i)
		if err = io.Export(); err != nil {
			return err
		}
		if err = io.Direction(sysfs.OUT); err != nil {
			return err
		}
		if err = io.Write(sysfs.LOW); err != nil {
			return err
		}
		if err = io.Unexport(); err != nil {
			return err
		}

	}

	// Arduino IO1, IO4 and IO13
	for _, i := range []string{"131", "129", "40"} {
		if err = changePinMode(i, "0"); err != nil {
			return err
		}
	}

	err = e.tristate.Write(sysfs.HIGH)
	return
}

// Name returns the EdisonAdaptors name
func (e *EdisonAdaptor) Name() string { return e.name }

// Connect initializes the Edison for use with the Arduino beakout board
func (e *EdisonAdaptor) Connect() (errs []error) {
	e.digitalPins = make(map[int]sysfs.DigitalPin)
	e.pwmPins = make(map[int]*pwmPin)
	if err := e.connect(e); err != nil {
		return []error{err}
	}
	return
}

// Finalize releases all i2c devices and exported analog, digital, pwm pins.
func (e *EdisonAdaptor) Finalize() (errs []error) {
	if err := e.tristate.Unexport(); err != nil {
		errs = append(errs, err)
	}
	for _, pin := range e.digitalPins {
		if pin != nil {
			if err := pin.Unexport(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	for _, pin := range e.pwmPins {
		if pin != nil {
			if err := pin.enable("0"); err != nil {
				errs = append(errs, err)
			}
			if err := pin.unexport(); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if e.i2cDevice != nil {
		if err := e.i2cDevice.Close(); errs != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// digitalPin returns matched digitalPin for specified values
func (e *EdisonAdaptor) digitalPin(pin string, dir string) (sysfsPin sysfs.DigitalPin, err error) {
	i := e.sysfsPinMap[pin]
	if e.digitalPins[i.gpioPin] == nil {
		e.digitalPins[i.gpioPin] = sysfs.NewDigitalPin(i.gpioPin)
		if err = e.digitalPins[i.gpioPin].Export(); err != nil {
			return
		}

		if i.resistor != 0 {
			e.digitalPins[i.resistor] = sysfs.NewDigitalPin(i.resistor)
			if err = e.digitalPins[i.resistor].Export(); err != nil {
				return
			}
		}

		if i.levelShifter != 0 {
			e.digitalPins[i.levelShifter] = sysfs.NewDigitalPin(i.levelShifter)
			if err = e.digitalPins[i.levelShifter].Export(); err != nil {
				return
			}
		}

		for _, mux := range i.mux {
			e.digitalPins[mux.pin] = sysfs.NewDigitalPin(mux.pin)
			if err = e.digitalPins[mux.pin].Export(); err != nil {
				return
			}

			if err = e.digitalPins[mux.pin].Direction(sysfs.OUT); err != nil {
				return
			}

			if err = e.digitalPins[mux.pin].Write(mux.value); err != nil {
				return
			}

		}
	}

	if dir == "in" {
		if err = e.digitalPins[i.gpioPin].Direction(sysfs.IN); err != nil {
			return
		}

		if i.resistor != 0 {
			if err = e.digitalPins[i.resistor].Direction(sysfs.OUT); err != nil {
				return
			}

			if err = e.digitalPins[i.resistor].Write(sysfs.LOW); err != nil {
				return
			}
		}

		if i.levelShifter != 0 {
			if err = e.digitalPins[i.levelShifter].Direction(sysfs.OUT); err != nil {
				return
			}

			if err = e.digitalPins[i.levelShifter].Write(sysfs.LOW); err != nil {
				return
			}
		}

	} else if dir == "out" {
		if err = e.digitalPins[i.gpioPin].Direction(sysfs.OUT); err != nil {
			return
		}

		if i.resistor != 0 {
			if err = e.digitalPins[i.resistor].Direction(sysfs.IN); err != nil {
				return
			}
		}

		if i.levelShifter != 0 {
			if err = e.digitalPins[i.levelShifter].Direction(sysfs.OUT); err != nil {
				return
			}

			if err = e.digitalPins[i.levelShifter].Write(sysfs.HIGH); err != nil {
				return
			}
		}

	}
	return e.digitalPins[i.gpioPin], nil
}

// DigitalRead reads digital value from pin
func (e *EdisonAdaptor) DigitalRead(gpioPin string) (i int, err error) {
	sysfsPin, err := e.digitalPin(gpioPin, "in")
	if err != nil {
		return
	}
	return sysfsPin.Read()
}

// DigitalWrite writes a value to the pin. Acceptable values are 1 or 0.
func (e *EdisonAdaptor) DigitalWrite(gpioPin string, val byte) (err error) {
	sysfsPin, err := e.digitalPin(gpioPin, "out")
	if err != nil {
		return
	}
	return sysfsPin.Write(int(val))
}

// PwmWrite writes the 0-254 value to the specified pin
func (e *EdisonAdaptor) PwmWrite(pin string, val byte) (err error) {
	sysPin := e.sysfsPinMap[pin]
	if sysPin.pwmPin != -1 {
		if e.pwmPins[sysPin.pwmPin] == nil {
			if err = e.DigitalWrite(pin, 1); err != nil {
				return
			}
			// FIXME: is this really the GPIO pin or the pwmPin we need to tweak ??
			// seems like an oversight, doesn't it?
			if err = changePinMode(strconv.Itoa(int(sysPin.gpioPin)), "1"); err != nil {
				return
			}
			e.pwmPins[sysPin.pwmPin] = newPwmPin(sysPin.pwmPin)
			if err = e.pwmPins[sysPin.pwmPin].export(); err != nil {
				return
			}
			if err = e.pwmPins[sysPin.pwmPin].enable("1"); err != nil {
				return
			}
		}
		p, err := e.pwmPins[sysPin.pwmPin].period()
		if err != nil {
			return err
		}
		period, err := strconv.Atoi(p)
		if err != nil {
			return err
		}
		duty := gobot.FromScale(float64(val), 0, 255.0)
		return e.pwmPins[sysPin.pwmPin].writeDuty(strconv.Itoa(int(float64(period) * duty)))
	}
	return errors.New("Not a PWM pin")
}

// AnalogRead returns value from analog reading of specified pin
func (e *EdisonAdaptor) AnalogRead(pin string) (val int, err error) {
	buf, err := readFile(
		"/sys/bus/iio/devices/iio:device1/in_voltage" + pin + "_raw",
	)
	if err != nil {
		return
	}

	val, err = strconv.Atoi(string(buf[0 : len(buf)-1]))

	return val / 4, err
}

// I2cStart initializes i2c device for addresss
func (e *EdisonAdaptor) I2cStart(address int) (err error) {
	if e.i2cDevice != nil {
		return
	}

	// FIXME: the I2cDevice interface doesn,t support multiple buses,
	// like the miniboard.. We select bus-6 by default here.

	if e.miniboard {
		for _, i := range []string{"20", "19"} /* bus-1 */ {
			//for _, i := range []string{"28", "27"} /* bus-6 */ {
			if err = changePinMode(i, "1"); err != nil {
				return
			}
		}

		e.i2cDevice, err = sysfs.NewI2cDevice("/dev/i2c-1", address)

	} else {

		if err = e.tristate.Write(sysfs.LOW); err != nil {
			return
		}

		// Confiure IO18, IO19
		// (ref: mraa's intel_edison_fab_c.c, func mraa_intel_edison_i2c_init_pre)
		for _, i := range []int{14, 165, 212, 213} {
			io := sysfs.NewDigitalPin(i)
			if err = io.Export(); err != nil {
				return
			}
			if err = io.Direction(sysfs.IN); err != nil {
				return
			}
			if err = io.Unexport(); err != nil {
				return
			}
		}

		// Continue on..
		for _, i := range []int{236, 237, 204, 205} {
			io := sysfs.NewDigitalPin(i)
			if err = io.Export(); err != nil {
				return
			}
			if err = io.Direction(sysfs.OUT); err != nil {
				return
			}
			if err = io.Write(sysfs.LOW); err != nil {
				return
			}
			if err = io.Unexport(); err != nil {
				return
			}
		}

		// Activate the I2c bus-6
		for _, i := range []string{"28", "27"} {
			if err = changePinMode(i, "1"); err != nil {
				return
			}
		}

		if err = e.tristate.Write(sysfs.HIGH); err != nil {
			return
		}

		e.i2cDevice, err = sysfs.NewI2cDevice("/dev/i2c-6", address)
	}

	return
}

// I2cReadBytesData uses I2C_RDWR to read bytes of data from a register.
func (e *EdisonAdaptor) I2cReadBytesData(addr uint16, reg byte, b []byte) error {
	i2cCtx := e.i2cDevice.(i2c.I2cReaderBytesData)
	return i2cCtx.I2cReadBytesData(addr, reg, b)
}

// I2cWrite writes data to i2c device
func (e *EdisonAdaptor) I2cWrite(address int, data []byte) (err error) {
	if err = e.i2cDevice.SetAddress(address); err != nil {
		return err
	}
	_, err = e.i2cDevice.Write(data)
	return
}

// I2cRead returns size bytes from the i2c device
func (e *EdisonAdaptor) I2cRead(address int, size int) (data []byte, err error) {
	data = make([]byte, size)
	if err = e.i2cDevice.SetAddress(address); err != nil {
		return
	}
	_, err = e.i2cDevice.Read(data)
	return
}
