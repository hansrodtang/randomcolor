package randomcolor

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type Options struct {
	Hue        int
	Luminosity string
	Color      *Color
	saturation int
}

func (o *Options) HueRange() Range {
	if o.Color != nil {
		return o.Color.HueRange
	} else {
		if o.Hue < 360 && o.Hue > 0 {
			return Range{o.Hue, o.Hue}
		}
	}

	return Range{0, 360}
}

func (o *Options) SaturationRange() Range {
	if o.Color != nil {
		return o.Color.SaturationRange()
	}
	return ColorInfo(o.Hue).SaturationRange()
}

func (o *Options) MinimumBrightness() int {
	var lowerBounds []Range
	if o.Color != nil {
		lowerBounds = o.Color.LowerBounds
	} else {
		lowerBounds = ColorInfo(o.Hue).LowerBounds
	}
	for i := 0; i < (len(lowerBounds) - 1); i++ {
		s1 := lowerBounds[i][0]
		v1 := lowerBounds[i][1]

		s2 := lowerBounds[i+1][0]
		v2 := lowerBounds[i+1][1]

		if o.saturation >= s1 && o.saturation <= s2 {
			m := (v2 - v1) / (s2 - s1)
			b := v1 - m*s1
			return m*o.saturation + b
		}
	}
	return 0
}

type Range [2]int

var Monochrome = Color{
	HueRange:    Range{0, 0},
	LowerBounds: []Range{{0, 0}, {100, 0}},
}

var Red = Color{
	HueRange:    Range{-26, 18},
	LowerBounds: []Range{{20, 100}, {30, 92}, {40, 89}, {50, 85}, {60, 78}, {70, 70}, {80, 60}, {90, 55}, {100, 50}},
}

var Orange = Color{
	HueRange:    Range{19, 46},
	LowerBounds: []Range{{20, 100}, {30, 93}, {40, 88}, {50, 86}, {60, 85}, {70, 70}, {100, 70}},
}

var Yellow = Color{
	HueRange:    Range{47, 62},
	LowerBounds: []Range{{25, 100}, {40, 94}, {50, 89}, {60, 86}, {70, 84}, {80, 82}, {90, 80}, {100, 75}},
}

var Green = Color{
	HueRange:    Range{63, 178},
	LowerBounds: []Range{{30, 100}, {40, 90}, {50, 85}, {60, 81}, {70, 74}, {80, 64}, {90, 50}, {100, 40}},
}

var Blue = Color{
	HueRange:    Range{179, 257},
	LowerBounds: []Range{{20, 100}, {30, 86}, {40, 80}, {50, 74}, {60, 60}, {70, 52}, {80, 44}, {90, 39}, {100, 35}},
}

var Purple = Color{
	HueRange:    Range{258, 282},
	LowerBounds: []Range{{20, 100}, {30, 87}, {40, 79}, {50, 70}, {60, 65}, {70, 59}, {80, 52}, {90, 45}, {100, 42}},
}

var Pink = Color{
	HueRange:    Range{283, 334},
	LowerBounds: []Range{{20, 100}, {30, 90}, {40, 86}, {60, 84}, {80, 80}, {90, 75}, {100, 73}},
}

var colors = []Color{Monochrome, Red, Orange, Yellow, Green, Blue, Purple, Pink}

type Color struct {
	H           float64
	S           float64
	B           float64
	HueRange    Range
	LowerBounds []Range
}

func (c Color) SaturationRange() Range {
	sMin := c.LowerBounds[0][0]
	sMax := c.LowerBounds[len(c.LowerBounds)-1][0]

	return Range{sMin, sMax}
}

func (c Color) BrightnessRange() Range {
	bMin := c.LowerBounds[len(c.LowerBounds)-1][1]
	bMax := c.LowerBounds[0][1]

	return Range{bMin, bMax}
}

func ColorInfo(hue int) Color {
	if hue >= 334 && hue <= 360 {
		hue = hue - 360
	}

	for _, color := range colors {
		if hue >= color.HueRange[0] && hue <= color.HueRange[1] {
			return color
		}
	}
	return Monochrome // Quick hack, add error handling instead
}

func (c *Color) RGBA() (r, g, b, a uint32) {
	h, s, v := c.H, c.S, c.B

	if h == 0 {
		h = 1
	}
	if h == 360 {
		h = 359
	}

	h = h / 360
	s = s / 100
	v = v / 100

	h_i := math.Floor(h * 6)
	f := h*6 - h_i
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)
	r2, g2, b2 := 256.0, 256.0, 256.0

	switch h_i {
	case 0:
		r2, g2, b2 = v, t, p
	case 1:
		r2, g2, b2 = q, v, p
	case 2:
		r2, g2, b2 = p, v, t
	case 3:
		r2, g2, b2 = p, q, v
	case 4:
		r2, g2, b2 = t, p, v
	case 5:
		r2, g2, b2 = v, p, q
	}

	r = uint32(math.Floor(r2 * 255))
	g = uint32(math.Floor(g2 * 255))
	b = uint32(math.Floor(b2 * 255))
	a = 1
	return
}

func NewColor(options Options) Color {
	c := Color{}
	options.Hue = setHue(options)
	c.H = float64(options.Hue)
	options.saturation = setSaturation(options)
	c.S = float64(options.saturation)
	c.B = float64(setBrightness(options))
	//fmt.Println(c.H, c.S, c.B)
	return c
}

func setHue(options Options) int {
	hueRange := options.HueRange()
	hue := randWithin(hueRange)

	if hue < 0 {
		hue = 360 + hue
	}
	return hue
}

func setSaturation(options Options) int {

	if options.Color == &Monochrome {
		return 0
	}

	saturationRange := options.SaturationRange()

	var sMin = saturationRange[0]
	var sMax = saturationRange[1]

	switch options.Luminosity {
	case "bright":
		sMin = 55
	case "dark":
		sMin = sMax - 10
	case "light":
		sMax = 55
	case "random":
		return randWithin(Range{0, 100})
	}
	return randWithin(Range{sMin, sMax})
}

func setBrightness(options Options) int {
	bMin := options.MinimumBrightness()
	bMax := 100

	switch options.Luminosity {
	case "dark":
		bMax = bMin + 20
	case "light":
		bMin = (bMax + bMin) / 2
	case "bright":
		//
	default:
		bMin = 0
		bMax = 100
	}
	return randWithin(Range{bMin, bMax})
}

func randWithin(r Range) int {
	return r[0] + rand.Intn(r[1]+1-r[0])
}

/*func main() {
	options := Options{
		Color: &Purple,
	}

	for i := 0; i < 1000; i++ {
		color := NewColor(options)
		fmt.Println(color.RGBA())
	}
}*/
