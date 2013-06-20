package main

import (
	"math"
)

type RGB struct {
	R, G, B float64
}

func RGBFromBytes(r, g, b byte) RGB {
	return RGB{float64(r) / 255, float64(g) / 255, float64(b) / 255}
}

func RGBFromHex(x uint32) RGB {
	return RGBFromBytes(
		byte((x>>16)&0xFF),
		byte((x>>8)&0xFF),
		byte(x&0xFF))
}

func (c RGB) ToBytes() (r, g, b byte) {
	return byte(255 * c.R), byte(255 * c.G), byte(255 * c.B)
}

func (c RGB) Times(ratio float64) RGB {
	if ratio > 1.0 {
		ratio = 1.0
	}
	return RGB{c.R * ratio, c.G * ratio, c.B * ratio}
}

func (c RGB) Plus(o RGB) RGB {
	return RGB{Min(c.R+o.R, 1.0), Min(c.G+o.G, 1.0), Min(c.B+o.B, 1.0)}
}

func (c RGB) ToHSV() HSV {
	r, g, b := c.R, c.G, c.B
	h, s, v := 0.0, 0.0, 0.0 // v

	min := Min(r, g, b)
	max := Max(r, g, b)
	delta := max - min

	v = max
	if max != 0 {
		s = delta / max // s
	} else {
		// r = g = b = 0        // s = 0, v is undefined
		s = 0
		h = -1
		return HSV{h, s, v}
	}

	if r == max {
		h = (g - b) / delta // between yellow & magenta
	} else if g == max {
		h = 2 + (b-r)/delta // between cyan & yellow
	} else {
		h = 4 + (r-g)/delta // between magenta & cyan
	}
	h *= 60 // degrees
	if h < 0 {
		h += 360
	}
	return HSV{h, s, v}
}

type HSV struct {
	H, S, V float64
}

func (c HSV) Times(ratio float64) HSV {
	if ratio > 1.0 {
		ratio = 1.0
	}
	return HSV{c.H * ratio, c.S * ratio, c.V * ratio}
}

func (c HSV) Plus(o HSV) HSV {
	h := c.H + o.H
	if h > 360 {
		h -= 360
	}
	return HSV{h, Min(c.S+o.S, 1.0), Min(c.V+o.V, 1.0)}
}

func (c HSV) ToRGB() RGB {
	h, s, v := c.H, c.S, c.V
	r, g, b := 0.0, 0.0, 0.0

	if s == 0 {
		// achromatic (grey)
		return RGB{v, v, v}
	}

	h /= 60 // sector 0 to 5
	i := math.Floor(h)
	f := h - i // factorial part of h
	p := v * (1 - s)
	q := v * (1 - s*f)
	t := v * (1 - s*(1-f))

	switch i {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	default: // case 5:
		r, g, b = v, p, q
	}
	return RGB{r, g, b}
}

func Min(fs ...float64) float64 {
	if len(fs) == 0 {
		return 0
	}
	min := math.MaxFloat64
	for _, f := range fs {
		if f < min {
			min = f
		}
	}
	return min
}

func Max(fs ...float64) float64 {
	if len(fs) == 0 {
		return 0
	}
	max := -math.MaxFloat64
	for _, f := range fs {
		if f > max {
			max = f
		}
	}
	return max
}
