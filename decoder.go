package bone

import (
	"errors"
)

type Decoder struct {
	Values []*Value
	Stack  []*Value
	Level  int
}

func (d *Decoder) Collapse() {
	for {
		l := len(d.Stack)
		if l == 0 {
			return
		}
		v := d.Stack[l-1]
		if !v.Complete() {
			return
		}
		d.Stack = d.Stack[:l-1]
		l--
		if l == 0 {
			d.Values = append(d.Values, v)
			return
		}
		d.Stack[l-1].Values = append(d.Stack[l-1].Values, v)
	}
}

func (d *Decoder) TerminateString(b byte) {
	l := len(d.Stack)
	if l > 0 {
		v := d.Stack[l-1]
		if b == 0x01 || !v.String() || len(v.Values) == 0 {
			return
		}
		v.Values = v.Values[:0]
		d.Stack = d.Stack[:l-1]
		l--
		if l == 0 {
			d.Values = append(d.Values, v)
		} else {
			d.Stack[l-1].Values = append(d.Stack[l-1].Values, v)
		}
		d.Collapse()
	}
}

func (d *Decoder) Accept(b byte) error {
	d.TerminateString(b)
	l := len(d.Stack)
	if l > 0 {
		v := d.Stack[l-1]
		if v.String() {
			if b == 0x00 {
				v.Values = append(v.Values, nil)
				return nil
			}
			if b == 0x01 && len(v.Values) == 1 {
				v.Values = v.Values[:0]
				v.Bytes = append(v.Bytes, 0x00)
				return nil
			}
			v.Bytes = append(v.Bytes, b)
			return nil
		}
		if v.Block() {
			v.Bytes = append(v.Bytes, b)
			d.Collapse()
			return nil
		}
		if b == 0x00 && v.List() {
			if d.Level != 0 {
				return errors.New("list terminated with non-zero level")
			}
			d.Stack = d.Stack[:l-1]
			l--
			if l == 0 {
				d.Values = append(d.Values, v)
			} else {
				d.Stack[l-1].Values = append(d.Stack[l-1].Values, v)
			}
			d.Collapse()
			return nil
		}
	}
	if b == 0xFF {
		d.Level++
		return nil
	}
	if b < 0x08 {
		return errors.New("illegal type code")
	}
	if d.Level > 0 && b < 0x20 {
		return errors.New("illegal level extension")
	}
	d.Stack = append(d.Stack, &Value{Code: b, Level: d.Level})
	d.Level = 0
	d.Collapse()
	return nil
}

func Decode(bytes []byte) ([]*Value, error) {
	decoder := Decoder{}
	for _, b := range bytes {
		if err := decoder.Accept(b); err != nil {
			return decoder.Values, err
		}
	}
	decoder.TerminateString(0xFF)
	if len(decoder.Stack) != 0 {
		return nil, errors.New("partial value left on stack")
	}
	if decoder.Level != 0 {
		return nil, errors.New("non zero level")
	}
	return decoder.Values, nil
}
