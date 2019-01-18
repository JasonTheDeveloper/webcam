package frame

import (
	"fmt"
	"image"
	"image/color"
)

type fRGB3 struct {
	model   color.Model
	b       image.Rectangle
	width   int
	frame   []byte
	release func()
}

// Register this framer for this format.
func init() {
	RegisterFramer("RGB3", newFramerRGB3)
}

// Return a function that is used as a framer for this format.
func newFramerRGB3(w, h int) func([]byte, func()) (Frame, error) {
	var size, bw int
	if *padded {
		bw = (w + 31) &^ 31
		size = 3 * bw * ((h + 15) &^ 15)
	} else {
		size = 3 * h * w
	}
	return func(b []byte, rel func()) (Frame, error) {
		return frameRGB3(size, bw, w, h, b, rel)
	}
}

// Wrap a raw webcam frame in a Frame so that it can be used as an image.
func frameRGB3(size, bw, w, h int, b []byte, rel func()) (Frame, error) {
	if len(b) != size {
		if rel != nil {
			defer rel()
		}
		return nil, fmt.Errorf("Wrong frame length (exp: %d, read %d)", size, len(b))
	}
	return &fRGB3{model: color.RGBAModel, b: image.Rect(0, 0, w, h), width: bw, frame: b, release: rel}, nil
}

func (f *fRGB3) ColorModel() color.Model {
	return f.model
}

func (f *fRGB3) Bounds() image.Rectangle {
	return f.b
}

func (f *fRGB3) At(x, y int) color.Color {
	i := f.width*y*3 + x*3
	return color.RGBA{f.frame[i], f.frame[i+1], f.frame[i+2], 0xFF}
}

// Done with frame, release back to camera (if required).
func (f *fRGB3) Release() {
	if f.release != nil {
		f.release()
	}
}
