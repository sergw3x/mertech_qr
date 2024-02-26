package mertech_qr

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/tarm/serial"
	"time"
)

var (
	_Head             = []byte{0x02, 0xf2, 0x02}
	_Tail             = []byte{0x02, 0xf2, 0x03}
	_Clear            = []byte{0x02, 0xf0, 0x03, 0x43, 0x4c, 0x53, 0x03}
	_ScreenOn         = []byte{0x02, 0xf0, 0x03, 0x42, 0x54, 0x4f, 0x4e, 0x03}
	_ScreenOff        = []byte{0x02, 0xf0, 0x03, 0x42, 0x54, 0x4f, 0x46, 0x46, 0x03}
	_CheckFirmware    = []byte{0x02, 0xF0, 0x03, 0x30, 0x44, 0x31, 0x33, 0x30, 0x32, 0x3F, 0x03}
	_PicOk            = []byte{0x02, 0xF0, 0x03, 0x63, 0x6F, 0x72, 0x72, 0x65, 0x63, 0x74, 0x03}
	_PicFalse         = []byte{0x02, 0xF0, 0x03, 0x6D, 0x69, 0x73, 0x74, 0x61, 0x6B, 0x65, 0x03}
	_EnableBluetooth  = []byte{0x02, 0xF0, 0x03, 0x42, 0x54, 0x4F, 0x4E, 0x03}
	_DisableBluetooth = []byte{0x02, 0x02, 0xF0, 0x03, 0x42, 0x54, 0x4F, 0x46, 0x46, 0x03}
)

const (
	DataBits  byte = 8
	SpeedBaud      = 115200
)

type MertechQr struct {
	conn        *serial.Port
	isConnected bool
	config      *serial.Config
}

type Config struct {
	Name        string
	Baud        int
	ReadTimeout time.Duration // Total timeout

	// Size is the number of data bits. If 0, DefaultSize is used.
	Size byte
}

func NewMertechQr(conf *Config) *MertechQr {

	return &MertechQr{config: &serial.Config{
		Name:        conf.Name,
		Baud:        conf.Baud,
		ReadTimeout: conf.ReadTimeout,
		Size:        conf.Size,
		Parity:      serial.ParityNone,
		StopBits:    serial.Stop1,
	}}
}

func (m *MertechQr) Connect() error {

	if m.config.Name == "" {
		return errors.New("empty serialNo")
	}

	conn, err := serial.OpenPort(m.config)
	if err != nil {
		return err
	}

	m.isConnected = true
	m.conn = conn

	return nil

}

func (m *MertechQr) Disconnect() error {
	return m.conn.Close()
}

func (m *MertechQr) CheckFirmware() ([]byte, error) {
	_, err := m.conn.Write(_CheckFirmware)
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Second)

	buf := make([]byte, 128)
	n, err := m.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], err
}

func (m *MertechQr) ShowPicOk() (int, error) {
	return m.conn.Write(_PicOk)
}

func (m *MertechQr) ShowPicFalse() (int, error) {
	return m.conn.Write(_PicFalse)
}

func (m *MertechQr) EnableBluetooth() (int, error) {
	return m.conn.Write(_EnableBluetooth)
}

func (m *MertechQr) DisableBluetooth() (int, error) {
	return m.conn.Write(_DisableBluetooth)
}

func (m *MertechQr) ScreenClear() (int, error) {
	return m.conn.Write(_Clear)
}

func (m *MertechQr) ScreenOn() (int, error) {
	return m.conn.Write(_ScreenOn)
}

// ScreenOff In my case, when working with windows, after disabling the display, the device changes the com port
func (m *MertechQr) ScreenOff() (int, error) {
	return m.conn.Write(_ScreenOff)
}

func (m *MertechQr) ShowQr(qrText string) (writeCnt int, readCnt int, result []byte, err error) {

	if qrText == "" {
		return writeCnt, readCnt, result, errors.New("qr string is empty")
	} else if len(qrText) > 1000 {
		return writeCnt, readCnt, result, errors.New("qr string is too large")
	}

	getBinaryLen(qrText)

	buff := concat([][]byte{
		_Head, {0x00}, getBinaryLen(qrText), []byte(qrText), _Tail,
	})

	readCnt, err = m.conn.Write(buff)
	if err != nil {
		return writeCnt, readCnt, result, err
	}

	time.Sleep(time.Second)

	buf := make([]byte, 128)
	writeCnt, err = m.conn.Read(buf)
	if err != nil {
		return writeCnt, readCnt, result, err
	}

	return writeCnt, readCnt, buf, nil

}

func getBinaryLen(text string) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, uint8(len(text)))
	if err != nil {
		return nil
	}
	return buf.Bytes()
}

func concat(slices [][]byte) []byte {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]byte, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}
