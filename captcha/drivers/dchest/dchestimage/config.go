package dchestimage

import (
	"bytes"
	"image/color"
)

type Config struct {
	// Maximum absolute skew factor of a single digit.
	MaxSkew float64
	// Number of background circles.
	CircleCount     int
	Font            *Font
	BackgroundColor *color.RGBA
}

func (c *Config) SavePNG(id string, digits []byte, width, height int) ([]byte, error) {
	img := NewImage(c, id, digits, width, height)
	buf := bytes.NewBuffer(nil)
	_, err := img.WriteTo(buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *Config) MustSavePNG(id string, digits []byte, width, height int) []byte {
	data, err := c.SavePNG(id, digits, width, height)
	if err != nil {
		panic(err)
	}
	return data
}

var DefaultConfig = &Config{
	MaxSkew:         0.7,
	CircleCount:     20,
	Font:            DefaultFont,
	BackgroundColor: &color.RGBA{0xFF, 0xFF, 0xFF, 0xFF},
}
