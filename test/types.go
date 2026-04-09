//go:generate go run ./../cmd/enum/ -types Color,Mode -case-insensitive
package test

import (
	"github.com/marvinpeter95/enum"
)

func init() {
	assertEnum[Color]()
	assertEnum[Mode]()
	assertEnumPointer[Color]()
	assertEnumPointer[Mode]()
}

func assertEnum[E enum.Enum[E]]() {}

func assertEnumPointer[E enum.Enum[E], EP enum.EnumPointer[E]]() {}

type ABC int

type Color string

const (
	ColorRed   Color = "red"
	ColorGreen Color = "green"
	ColorBlue  Color = "blue"
	XXX        int   = 1
	DD
)

type Pill string

const (
	Placebo       Pill = "placebo"
	Aspirin       Pill = "aspirin"
	Ibuprofen     Pill = "ibuprofen"
	Paracetamol   Pill = "paracetamol"
	Acetaminophen      = Paracetamol
)

type Mode int

const (
	ModeLight Mode = iota
	ModeDark
	ModeAuto
	ModeSystem = ModeAuto
)

func ParseMode(value string) (Mode, error) {
	switch value {
	case "0", "light", "l":
		return ModeLight, nil
	case "1", "dark", "d":
		return ModeDark, nil
	case "2", "system", "auto":
		return ModeAuto, nil
	default:
		var v Mode
		return v, v.Validate()
	}
}
