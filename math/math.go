package math

type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

func Sum[num Number](arr []num) num {
	var s num

	for _, v := range arr {
		s += num(v)
	}

	return s
}

func Max[num Number](arr []num) num {
	var m num

	for _, v := range arr {
		if m < v {
			m = v
		}
	}

	return m
}
