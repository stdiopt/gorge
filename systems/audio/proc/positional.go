package proc

import (
	"io"
	"math"
	"unsafe"

	"github.com/stdiopt/gorge/m32"
)

// Positional controls the audio levels based on a 3D position
// TODO: implement panning
type Positional struct {
	player io.Writer
	// Should have a format but lets assume 44100/2/16
	Position m32.Vec3
	out      []int16
}

// NewPositional returns a new position audio writer.
func NewPositional(target io.Writer) *Positional {
	return &Positional{player: target}
}

func (p *Positional) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	if p.out == nil {
		p.out = make([]int16, 10)
	}
	data := (*[math.MaxUint16]int16)(unsafe.Pointer(&b[0]))[0 : len(b)/2 : len(b)/2]
	out := p.out
	totalWrite := 0
	for len(data) > 0 {
		i := 0
		for ; i < len(data) && i < len(out); i++ {
			maxDist := float32(20) // if we reached max volume should be 0
			dist := p.Position.Len()
			distVol := (maxDist - dist) / maxDist
			distVol = m32.Max(distVol, 0)
			out[i] = int16(float32(data[i]) * distVol)
		}
		data = data[i:]
		totalWrite += i

		bout := (*[math.MaxUint16]byte)(unsafe.Pointer(&out[0]))[0 : len(out)*2 : len(out)*2]
		_, err := p.player.Write(bout)
		if err != nil {
			return totalWrite * 2, err
		}
	}
	return len(b), nil
	// return int(wrote), err
}
