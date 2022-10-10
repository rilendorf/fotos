package neopixel

func Smooth(length int, progress float32) (f []uint8) {
	if length <= 0 {
		return []uint8{}
	}

	if progress < 0 || progress > 1 {
		return []uint8{}
	}

	f = make([]uint8, length)

	rp := progress * float32(255*length)

	for k := range f {
		if rp >= 255 {
			rp -= 255

			f[k] = 255
		} else {
			f[k] = uint8(rp)
			break
		}
	}

	return
}
