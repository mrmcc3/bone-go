package bone

type StackItem struct {
	v *Value
	i int
}

func Encode(values []*Value) []byte {
	res := []byte{}
	stack := []*StackItem{}
	l := 0
	for _, top := range values {
		stack = append(stack, &StackItem{v: top})
		l++
		for l > 0 {
			s := stack[l-1]
			if s.i == 0 {
				for range s.v.Level {
					res = append(res, 0xFF)
				}
				res = append(res, s.v.Code)
			}
			if s.v.String() {
				for _, b := range s.v.Bytes {
					res = append(res, b)
					if b == 0x00 {
						res = append(res, 0x01)
					}
				}
			} else if s.v.Block() {
				res = append(res, s.v.Bytes...)
			}
			if s.i < len(s.v.Values) {
				stack = append(stack, &StackItem{v: s.v.Values[s.i]})
				l++
				s.i++
			} else {
				if s.v.List() || s.v.String() {
					res = append(res, 0x00)
				}
				stack = stack[:l-1]
				l--
			}
		}
	}
	return res
}
