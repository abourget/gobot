package sysfs

import (
	"fmt"
	"io"
	"os"
	"syscall"
	"unsafe"
)

const (
	// ioctl signals
	I2C_SLAVE = 0x0703
	I2C_SMBUS = 0x0720
	I2C_RDWR  = 0x0707
	// I2C Message flags
	I2C_M_RD = 0x01
	// Read/write markers
	I2C_SMBUS_READ  = 1
	I2C_SMBUS_WRITE = 0
	// Transaction types
	I2C_SMBUS_BYTE                = 1
	I2C_SMBUS_BYTE_DATA           = 2
	I2C_SMBUS_WORD_DATA           = 3
	I2C_SMBUS_PROC_CALL           = 4
	I2C_SMBUS_BLOCK_DATA          = 5
	I2C_SMBUS_I2C_BLOCK_DATA      = 6
	I2C_SMBUS_BLOCK_PROC_CALL     = 7  /* SMBus 2.0 */
	I2C_SMBUS_BLOCK_DATA_PEC      = 8  /* SMBus 2.0 */
	I2C_SMBUS_PROC_CALL_PEC       = 9  /* SMBus 2.0 */
	I2C_SMBUS_BLOCK_PROC_CALL_PEC = 10 /* SMBus 2.0 */
	I2C_SMBUS_WORD_DATA_PEC       = 11 /* SMBus 2.0 */
)

type i2cSmbusIoctlData struct {
	readWrite byte
	command   byte
	size      uint32
	data      uintptr
}

type I2cDevice interface {
	io.ReadWriteCloser
	SetAddress(int) error
}

type i2cDevice struct {
	file        File
	lastAddress int
}

// NewI2cDevice returns an io.ReadWriteCloser with the proper ioctrl given
// an i2c bus location and device address
func NewI2cDevice(location string, address int) (d *i2cDevice, err error) {
	d = &i2cDevice{}

	if d.file, err = OpenFile(location, os.O_RDWR, os.ModeExclusive); err != nil {
		return
	}

	err = d.SetAddress(address)

	return
}

func (d *i2cDevice) SetAddress(address int) (err error) {
	d.lastAddress = address

	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_SLAVE,
		uintptr(byte(address)),
	)

	if errno != 0 {
		err = fmt.Errorf("Failed with syscall.Errno %v", errno)
	}

	return
}

func (d *i2cDevice) Close() (err error) {
	return d.file.Close()
}

func (d *i2cDevice) Read(b []byte) (n int, err error) {
	return d.SmbusReadI2cBlockData(b)
	//return d.SmbusReadBlockDataPec(b)
}

func (d *i2cDevice) SmbusReadI2cBlockData(b []byte) (n int, err error) {
	data := make([]byte, len(b)+1)
	data[0] = byte(len(b))

	err = d.smbusAccess(I2C_SMBUS_READ, 0, I2C_SMBUS_I2C_BLOCK_DATA, uintptr(unsafe.Pointer(&data[0])))
	if err != nil {
		return
	}

	copy(b, data[1:])

	return int(data[0]), nil
}

func (d *i2cDevice) SmbusReadBlockDataPec(b []byte) (n int, err error) {
	data := make([]byte, len(b)+1)
	data[0] = byte(len(b))

	fmt.Println("DATA:", data)
	err = d.smbusAccess(I2C_SMBUS_READ, 0, I2C_SMBUS_BLOCK_DATA_PEC, uintptr(unsafe.Pointer(&data[0])))
	if err != nil {
		return
	}

	copy(b, data[1:])

	return int(data[0]), nil
}

func (d *i2cDevice) SmbusReadByte() (b byte, err error) {
	err = d.smbusAccess(I2C_SMBUS_READ, 0, I2C_SMBUS_BYTE, uintptr(unsafe.Pointer(&b)))
	if err != nil {
		return
	}

	return b, nil
}

type i2cMessage struct {
	addr  uint16
	flags uint16
	len   uint16
	buf   uintptr
}
type i2cRdWrIoctlData struct {
	msgs  uintptr
	nmsgs uint32
}

func (d *i2cDevice) I2cReadBytesData(addr uint16, reg byte, b []byte) (err error) {
	msgs := []i2cMessage{
		{
			addr:  addr,
			flags: 0x00,
			len:   1,
			buf:   uintptr(unsafe.Pointer(&reg)),
		},
		{
			addr:  addr,
			flags: I2C_M_RD,
			len:   uint16(len(b)),
			buf:   uintptr(unsafe.Pointer(&b[0])),
		},
	}
	data := i2cRdWrIoctlData{
		msgs:  uintptr(unsafe.Pointer(&msgs[0])),
		nmsgs: 2,
	}
	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_RDWR,
		uintptr(unsafe.Pointer(&data)),
	)

	if errno != 0 {
		return fmt.Errorf("Failed with syscall.Errno %v", errno)
	}

	return nil
}

func (d *i2cDevice) smbusAccess(readWrite byte, command byte, size uint32, data uintptr) error {
	smbus := &i2cSmbusIoctlData{
		readWrite: readWrite,
		command:   0,
		size:      size,
		data:      data,
	}

	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_SMBUS,
		uintptr(unsafe.Pointer(smbus)),
	)

	if errno != 0 {
		return fmt.Errorf("Failed with syscall.Errno %v", errno)
	}

	return nil
}

func (d *i2cDevice) Write(b []byte) (n int, err error) {
	if len(b) <= 2 {
		return d.file.Write(b)
	}

	command := byte(b[0])
	buf := b[1:]

	data := make([]byte, len(buf)+1)
	data[0] = byte(len(buf))

	copy(data[1:], buf)

	smbus := &i2cSmbusIoctlData{
		readWrite: I2C_SMBUS_WRITE,
		command:   command,
		size:      I2C_SMBUS_BLOCK_DATA_PEC,
		data:      uintptr(unsafe.Pointer(&data[0])),
	}

	_, _, errno := Syscall(
		syscall.SYS_IOCTL,
		d.file.Fd(),
		I2C_SMBUS,
		uintptr(unsafe.Pointer(smbus)),
	)

	if errno != 0 {
		err = fmt.Errorf("Failed with syscall.Errno %v", errno)
	}

	return len(b), err
}
