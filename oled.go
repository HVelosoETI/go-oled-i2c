package goled

import (
	i2c "github.com/talkkonnect/go-i2c"
	"log"
)

const (
	//	Registers
	oledRegisterData    = 0x40
	oledRegisterCommand = 0x80

	//	Commands
	oledCommandSetLowColumn  = 0x00
	oledCommandSetHighColumn = 0x10
	oledCommandDisplayOff    = 0xAE
	oledCommandDisplayOn     = 0xAF
)

var (
	OLEDDefaultI2cAddress uint8 = 0
	OLEDDefaultI2cBus = 1
        OLEDScreenWidth = 130
	OLEDScreenHeight = 64
	OLEDDisplayRows = 8
	OLEDDisplayColumns byte = 21
	OLEDCharLength = 6
	OLEDCommandColumnAddressing = 0x21
	OLEDAddressBasePageStart = 0
	OLEDStartColumn = 1
)

var oledASCIITable = [...][6]byte{
	{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, // SPACE

	{0x00, 0x00, 0x4F, 0x00, 0x00, 0x00}, // !
	{0x00, 0x07, 0x00, 0x07, 0x00, 0x00}, // "
	{0x14, 0x7F, 0x14, 0x7F, 0x14, 0x00}, // #
	{0x24, 0x2A, 0x7F, 0x2A, 0x12, 0x00}, // $
	{0x23, 0x13, 0x08, 0x64, 0x62, 0x00}, // %
	{0x36, 0x49, 0x55, 0x22, 0x50, 0x00}, // &
	{0x00, 0x05, 0x03, 0x00, 0x00, 0x00}, // '
	{0x00, 0x1C, 0x22, 0x41, 0x00, 0x00}, // (
	{0x00, 0x41, 0x22, 0x1C, 0x00, 0x00}, // )
	{0x14, 0x08, 0x3E, 0x08, 0x14, 0x00}, // *
	{0x08, 0x08, 0x3E, 0x08, 0x08, 0x00}, // +
	{0x00, 0x50, 0x30, 0x00, 0x00, 0x00}, // ,
	{0x08, 0x08, 0x08, 0x08, 0x08, 0x00}, // -
	{0x00, 0x60, 0x60, 0x00, 0x00, 0x00}, // .
	{0x20, 0x10, 0x08, 0x04, 0x02, 0x00}, // /

	{0x3E, 0x51, 0x49, 0x45, 0x3E, 0x00}, // 0
	{0x00, 0x42, 0x7F, 0x40, 0x00, 0x00}, // 1
	{0x42, 0x61, 0x51, 0x49, 0x46, 0x00}, // 2
	{0x21, 0x41, 0x45, 0x4B, 0x31, 0x00}, // 3
	{0x18, 0x14, 0x12, 0x7F, 0x10, 0x00}, // 4
	{0x27, 0x45, 0x45, 0x45, 0x39, 0x00}, // 5
	{0x3C, 0x4A, 0x49, 0x49, 0x30, 0x00}, // 6
	{0x01, 0x71, 0x09, 0x05, 0x03, 0x00}, // 7
	{0x36, 0x49, 0x49, 0x49, 0x36, 0x00}, // 8
	{0x06, 0x49, 0x49, 0x29, 0x1E, 0x00}, // 9

	{0x36, 0x36, 0x00, 0x00, 0x00, 0x00}, // :
	{0x56, 0x36, 0x00, 0x00, 0x00, 0x00}, // ;
	{0x08, 0x14, 0x22, 0x41, 0x00, 0x00}, // <
	{0x14, 0x14, 0x14, 0x14, 0x14, 0x00}, // =
	{0x00, 0x41, 0x22, 0x14, 0x08, 0x00}, // >
	{0x02, 0x01, 0x51, 0x09, 0x06, 0x00}, // ?
	{0x30, 0x49, 0x79, 0x41, 0x3E, 0x00}, // @

	{0x7E, 0x11, 0x11, 0x11, 0x7E, 0x00}, // A
	{0x7F, 0x49, 0x49, 0x49, 0x36, 0x00}, // B
	{0x3E, 0x41, 0x41, 0x41, 0x22, 0x00}, // C
	{0x7F, 0x41, 0x41, 0x22, 0x1C, 0x00}, // D
	{0x7F, 0x49, 0x49, 0x49, 0x41, 0x00}, // E
	{0x7F, 0x09, 0x09, 0x09, 0x01, 0x00}, // F
	{0x3E, 0x41, 0x49, 0x49, 0x7A, 0x00}, // G
	{0x7F, 0x08, 0x08, 0x08, 0x7F, 0x00}, // H
	{0x00, 0x41, 0x7F, 0x41, 0x00, 0x00}, // I
	{0x20, 0x40, 0x41, 0x3F, 0x01, 0x00}, // J
	{0x7F, 0x08, 0x14, 0x22, 0x41, 0x00}, // K
	{0x7F, 0x40, 0x40, 0x40, 0x40, 0x00}, // L
	{0x7F, 0x02, 0x0C, 0x02, 0x7F, 0x00}, // M
	{0x7F, 0x04, 0x08, 0x10, 0x7F, 0x00}, // N
	{0x3E, 0x41, 0x41, 0x41, 0x3E, 0x00}, // O
	{0x7F, 0x09, 0x09, 0x09, 0x06, 0x00}, // P
	{0x3E, 0x41, 0x51, 0x21, 0x5E, 0x00}, // Q
	{0x7F, 0x09, 0x19, 0x29, 0x46, 0x00}, // R
	{0x46, 0x49, 0x49, 0x49, 0x31, 0x00}, // S
	{0x01, 0x01, 0x7F, 0x01, 0x01, 0x00}, // T
	{0x3F, 0x40, 0x40, 0x40, 0x3F, 0x00}, // U
	{0x1F, 0x20, 0x40, 0x20, 0x1F, 0x00}, // V
	{0x3F, 0x40, 0x30, 0x40, 0x3F, 0x00}, // W
	{0x63, 0x14, 0x08, 0x14, 0x63, 0x00}, // X
	{0x07, 0x08, 0x70, 0x08, 0x07, 0x00}, // Y
	{0x61, 0x51, 0x49, 0x45, 0x43, 0x00}, // Z

	{0x00, 0x7F, 0x41, 0x41, 0x00, 0x00}, // [
	{0x02, 0x04, 0x08, 0x10, 0x20, 0x00}, // backslash
	{0x00, 0x41, 0x41, 0x7F, 0x00, 0x00}, // ]
	{0x04, 0x02, 0x01, 0x02, 0x04, 0x00}, // ^
	{0x40, 0x40, 0x40, 0x40, 0x40, 0x00}, // _
	{0x00, 0x01, 0x02, 0x04, 0x00, 0x00}, // `

	{0x20, 0x54, 0x54, 0x54, 0x78, 0x00}, // a
	{0x7F, 0x50, 0x48, 0x48, 0x30, 0x00}, // b
	{0x38, 0x44, 0x44, 0x44, 0x20, 0x00}, // c
	{0x38, 0x44, 0x44, 0x48, 0x7F, 0x00}, // d
	{0x38, 0x54, 0x54, 0x54, 0x18, 0x00}, // e
	{0x08, 0x7E, 0x09, 0x01, 0x02, 0x00}, // f
	{0x0C, 0x52, 0x52, 0x52, 0x3E, 0x00}, // g
	{0x7F, 0x08, 0x04, 0x04, 0x78, 0x00}, // h
	{0x00, 0x44, 0x7D, 0x40, 0x00, 0x00}, // i
	{0x20, 0x40, 0x44, 0x3D, 0x00, 0x00}, // j
	{0x7F, 0x10, 0x28, 0x44, 0x00, 0x00}, // k
	{0x00, 0x41, 0x7F, 0x40, 0x00, 0x00}, // l
	{0x78, 0x04, 0x78, 0x04, 0x78, 0x00}, // m
	{0x7C, 0x08, 0x04, 0x04, 0x78, 0x00}, // n
	{0x38, 0x44, 0x44, 0x44, 0x38, 0x00}, // o
	{0x7C, 0x14, 0x14, 0x14, 0x08, 0x00}, // p
	{0x08, 0x14, 0x14, 0x18, 0x7C, 0x00}, // q
	{0x7C, 0x08, 0x04, 0x04, 0x08, 0x00}, // r
	{0x48, 0x54, 0x54, 0x54, 0x20, 0x00}, // s
	{0x04, 0x3F, 0x44, 0x40, 0x20, 0x00}, // t
	{0x3C, 0x40, 0x40, 0x20, 0x7C, 0x00}, // u
	{0x1C, 0x20, 0x40, 0x20, 0x1C, 0x00}, // v
	{0x3C, 0x40, 0x30, 0x40, 0x3C, 0x00}, // w
	{0x44, 0x28, 0x10, 0x28, 0x44, 0x00}, // x
	{0x0C, 0x50, 0x50, 0x50, 0x3C, 0x00}, // y
	{0x44, 0x64, 0x54, 0x4C, 0x44, 0x00}, // z

	{0x00, 0x08, 0x36, 0x41, 0x00, 0x00}, // {
	{0x00, 0x00, 0x7F, 0x00, 0x00, 0x00}, // |
	{0x00, 0x41, 0x36, 0x08, 0x00, 0x00}, // }
	{0x0C, 0x02, 0x0C, 0x10, 0x0C, 0x00}, // ~
	{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}}

// Oled Represents oled display
type Oled struct {
	_i2c          *i2c.I2C
	currentRow    byte
	currentColumn byte
}

// BeginOled Creates a new Oled reference
func BeginOled(mOLEDDefaultI2cAddress uint8, mOLEDDefaultI2cBus int, mOLEDScreenWidth int, mOLEDScreenHeight int, mOLEDDisplayRows int, mOLEDDisplayColumns uint8, mOLEDStartColumn int, mOLEDCharLength int, mOLEDCommandColumnAddressing int, mOLEDAddressBasePageStart int)  (*Oled, error) {

	OLEDDefaultI2cAddress = mOLEDDefaultI2cAddress
	OLEDDefaultI2cBus = mOLEDDefaultI2cBus
        OLEDScreenWidth = mOLEDScreenWidth
	OLEDScreenHeight = mOLEDScreenHeight
	OLEDDisplayRows = mOLEDDisplayRows
	OLEDDisplayColumns = mOLEDDisplayColumns
	OLEDStartColumn = mOLEDStartColumn
	OLEDCharLength = mOLEDCharLength
	OLEDCommandColumnAddressing = mOLEDCommandColumnAddressing
	OLEDAddressBasePageStart = mOLEDAddressBasePageStart

	res := &Oled{}
	err := res.init()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Init intiializes oled display
func (v *Oled) init() error {
	_i2c, err := i2c.NewI2C(OLEDDefaultI2cAddress, OLEDDefaultI2cBus)
	if err != nil {
		return err
	}
	v._i2c = _i2c
	// for 128 x 64 Init Commands
	res, err := v.sendOledCommand(0xAE)
	res, err =  v.sendOledCommand(0xA8)
	res, err =  v.sendOledCommand(0x3F)
	res, err =  v.sendOledCommand(0xD3)
	res, err =  v.sendOledCommand(0x00)
	res, err =  v.sendOledCommand(0x40)
	res, err =  v.sendOledCommand(0xA1)
	res, err =  v.sendOledCommand(0xC8)

	res, err = v.sendOledCommand(0xA6)
	res, err = v.sendOledCommand(0xD5)
	res, err = v.sendOledCommand(0x80)
	res, err = v.sendOledCommand(0xDA)
	res, err = v.sendOledCommand(0x12)
	res, err = v.sendOledCommand(0x81)
	res, err = v.sendOledCommand(0xFF)

	res, err = v.sendOledCommand(0xA4)
	res, err = v.sendOledCommand(0xDB)
	res, err = v.sendOledCommand(0x40)
	res, err = v.sendOledCommand(0x20)
	res, err = v.sendOledCommand(0x00)
	res, err = v.sendOledCommand(0x00)
	res, err = v.sendOledCommand(0x10)
	res, err = v.sendOledCommand(0x8D)

	res, err = v.sendOledCommand(0x14)
	res, err = v.sendOledCommand(0x2E)
	res, err = v.sendOledCommand(0xA6)
	res, err = v.sendOledCommand(0xAF)


	if res != 2 {
		log.Println("warn Some I2C Error Result should be = 2 actual result is ",res)
		log.Printf("warn OLED I2C Address %v bus %v ",OLEDDefaultI2cAddress,OLEDDefaultI2cBus)
	}
	if err != nil {
		return err
	}

	return nil
}

// sendOledCommand sends the specified command to oled
func (v *Oled) sendOledCommand(command int) (int, error) {
	return v._i2c.WriteBytes([]byte{oledRegisterCommand, byte(command)})
}

func (v *Oled) sendOledData(data int) (int, error) {
	return v._i2c.WriteBytes([]byte{oledRegisterData, byte(data)})
}

// SetColumnAddressing sets the oled viewport
func (v *Oled) SetColumnAddressing(startPixel int, endPixel int) (int, error) {
	res, err := v.sendOledCommand(OLEDCommandColumnAddressing)
	if err != nil {
		return res, err
	}
	res, err = v.sendOledCommand(startPixel)
	if err != nil {
		return res, err
	}
	res, err = v.sendOledCommand(endPixel)
	if err != nil {
		return res, err
	}
	v.currentRow = 0
	return res, nil
}

// DisplayOff switches off the oled display
func (v *Oled) DisplayOff() error {
	_, err := v.sendOledCommand(oledCommandDisplayOff)
	if err != nil {
		return err
	}
	return nil
}

// DisplayOn switches on the oled display
func (v *Oled) DisplayOn() error {
	_, err := v.sendOledCommand(oledCommandDisplayOn)
	if err != nil {
		return err
	}
	return nil
}

// SetCursor sets the cursor at specified row and column
func (v *Oled) SetCursor(row int, column int) error {
	_, err := v.sendOledCommand(OLEDAddressBasePageStart + row)
	if err != nil {
		return err
	}

	_, err = v.sendOledCommand(oledCommandSetLowColumn + (OLEDCharLength * column & 0x0F))
	if err != nil {
		return err
	}

	_, err = v.sendOledCommand(oledCommandSetHighColumn + ((OLEDCharLength * column >> 4) & 0x0F))
	if err != nil {
		return err
	}
	v.currentRow = byte(row)
	return nil
}

func (v *Oled) writeBlankChars() error {
	for row := 0; row < OLEDDisplayRows; row++ {
		err := v.SetCursor(row, 0)
		if err != nil {
			return err
		}
		for pixel := 0; pixel <OLEDScreenWidth; pixel++ {
			_, err = v.sendOledData(0x00)
			if err != nil {
				return err
			}
		}
	}
	v.currentRow = 0
	v.currentColumn = 0
	return nil
}

// Clear clears oled screen
func (v *Oled) Clear() error {
	if _, err := v.SetColumnAddressing(0,OLEDScreenWidth-1); err != nil {
		return err
	}
	if err := v.writeBlankChars(); err != nil {
		return err
	}
	return nil
}

// WriteCharUnchecked write specified character as it is on oled
func (v *Oled) WriteCharUnchecked(c int) error {
	if _, err := v.sendOledData(c); err != nil {
		return err
	}
	return nil
}

// WriteChar writes specified character on oled display
func (v *Oled) WriteChar(c int) error {
	index := c - 32
	for ctr := 0; ctr < len(oledASCIITable[index]); ctr++ {
		if _, err := v.sendOledData(int(oledASCIITable[index][ctr])); err != nil {
			return err
		}
	}
	v.currentColumn++
	if v.currentColumn > OLEDDisplayColumns {
		v.currentRow++
		v.currentColumn = 0
	}
	return nil
}

// Write writes the specified string on oled
func (v *Oled) Write(message string) (int, error) {
	var res int
	messageLen := len(message)
	for ctr := 0; ctr < messageLen; ctr++ {
		if message[ctr] == '\n' {
			v.currentRow++
			if err := v.SetCursor(int(v.currentRow), 0); err != nil {
				return res, err
			}
			continue
		}
		if err := v.WriteChar(int(message[ctr])); err != nil {
			return res, err
		}
		res++
	}
	return res, nil
}

// Close closes the oled i2c bus
func (v *Oled) Close() {
	v._i2c.Close()
}
