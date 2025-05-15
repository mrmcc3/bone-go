package bone

type Value struct {
	Code   byte
	Level  int
	Bytes  []byte
	Values []*Value
}

func (v *Value) Block() bool {
	return v.Code >= 0x08 && v.Code < 0x90
}

func (v *Value) String() bool {
	return v.Code >= 0x90 && v.Code < 0xA0
}

func (v *Value) List() bool {
	return v.Code >= 0xF0 && v.Code < 0xFF
}

func (v *Value) Complete() bool {
	if v.Code < 0x08 {
		panic("illegal")
	}
	switch v.Code {
	case 0x0F, 0x18:
		return len(v.Bytes) == 1
	case 0x0E, 0x19:
		return len(v.Bytes) == 2
	case 0x0D, 0x1A:
		return len(v.Bytes) == 3
	case 0x0C, 0x1B:
		return len(v.Bytes) == 4
	case 0x0B, 0x1C:
		return len(v.Bytes) == 5
	case 0x0A, 0x1D:
		return len(v.Bytes) == 6
	case 0x09, 0x1E:
		return len(v.Bytes) == 7
	case 0x08, 0x1F:
		return len(v.Bytes) == 8
	}
	switch {
	case v.Code < 0x30:
		return true
	case v.Code < 0x40:
		return len(v.Bytes) == 1
	case v.Code < 0x50:
		return len(v.Bytes) == 2
	case v.Code < 0x60:
		return len(v.Bytes) == 3
	case v.Code < 0x70:
		return len(v.Bytes) == 4
	case v.Code < 0x80:
		return len(v.Bytes) == 8
	case v.Code < 0x90:
		return len(v.Bytes) == 16
	case v.Code < 0xA0:
		return false
	case v.Code < 0xB0:
		return len(v.Values) == 1
	case v.Code < 0xC0:
		return len(v.Values) == 2
	case v.Code < 0xD0:
		return len(v.Values) == 3
	case v.Code < 0xE0:
		return len(v.Values) == 4
	case v.Code < 0xFF:
		return false
	}
	panic("illegal")
}
