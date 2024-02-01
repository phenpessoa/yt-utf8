package main

import "errors"

func main() {}

// Char. number range  |        UTF-8 octet sequence
//    (hexadecimal)    |              (binary)
// --------------------+---------------------------------------------
// 0000 0000-0000 007F | 0xxxxxxx
// 0000 0080-0000 07FF | 110xxxxx 10xxxxxx
// 0000 0800-0000 FFFF | 1110xxxx 10xxxxxx 10xxxxxx
// 0001 0000-0010 FFFF | 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx

func decodeRune(b []byte) (r rune, s int, err error) {
	if len(b) == 0 {
		return 0, 0, errors.New("empty input")
	}
	b0 := b[0]

	switch {
	case b0 < 0x80: // ASCII
		if len(b) > 1 {
			if b[1]&0xC0 == 0x80 {
				return 0, 0, errors.New("invalid length")
			}
		}
		r = rune(b0)
		s = 1
	case b0&0xE0 == 0xC0: // 2 byte character
		if len(b) < 2 {
			return 0, 0, errors.New("invalid length")
		}
		if len(b) > 2 {
			if b[2]&0xC0 == 0x80 {
				return 0, 0, errors.New("invalid length")
			}
		}

		b1 := b[1]

		if b1&0xC0 != 0x80 {
			return 0, 0, errors.New("invalid continuation byte")
		}

		s = 2
		r = ((rune(b0) & 0x1F) << 6) |
			(rune(b1) & 0x3F)

		if r < 0x80 {
			return 0, 0, errors.New("overlong")
		}
	case b0&0xF0 == 0xE0: // 3 byte character
		if len(b) < 3 {
			return 0, 0, errors.New("invalid length")
		}
		if len(b) > 3 {
			if b[3]&0xC0 == 0x80 {
				return 0, 0, errors.New("invalid length")
			}
		}
		b1 := b[1]
		b2 := b[2]

		if b0 == 0xE0 && b1 < 0xA0 {
			return 0, 0, errors.New("overlong")
		}

		if b1&0xC0 != 0x80 || b2&0xC0 != 0x80 {
			return 0, 0, errors.New("invalid continuation byte")
		}
		s = 3
		r = ((rune(b0) & 0x0F) << 12) |
			((rune(b1) & 0x3F) << 6) |
			(rune(b2) & 0x3F)
	case b0&0xF8 == 0xF0: // 4 byte character
		if len(b) < 4 {
			return 0, 0, errors.New("invalid length")
		}
		if len(b) > 4 {
			if b[4]&0xC0 == 0x80 {
				return 0, 0, errors.New("invalid length")
			}
		}
		b1 := b[1]
		b2 := b[2]
		b3 := b[3]

		if b0 == 0xF0 && b1 < 0x90 {
			return 0, 0, errors.New("overlong")
		}

		if b1&0xC0 != 0x80 || b2&0xC0 != 0x80 || b3&0xC0 != 0x80 {
			return 0, 0, errors.New("invalid continuation byte")
		}

		s = 4
		r = ((rune(b0) & 0x07) << 18) |
			((rune(b1) & 0x3F) << 12) |
			((rune(b2) & 0x3F) << 6) |
			(rune(b3) & 0x3F)
	default:
		return 0, 0, errors.New("invalid utf8")
	}

	if r >= 0xD800 && r <= 0xDFFF {
		return 0, 0, errors.New("surrogate halfs")
	}

	if r > 0x10FFFF {
		return 0, 0, errors.New("too big")
	}

	return r, s, nil
}
