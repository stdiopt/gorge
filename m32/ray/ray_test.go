package ray_test

import (
	"testing"

	"github.com/stdiopt/gorge/m32"
	"github.com/stdiopt/gorge/m32/ray"
)

func TestCalcNormal(t *testing.T) {
	tests := []struct {
		a, b, c m32.Vec3
		want    m32.Vec3
	}{
		{
			a:    m32.Vec3{0, 0, 0},
			b:    m32.Vec3{1, 0, 0},
			c:    m32.Vec3{0, 1, 0},
			want: m32.Vec3{0, 0, 1},
		},
		{
			a:    m32.Vec3{0, 0, 0},
			b:    m32.Vec3{-1, 0, 0},
			c:    m32.Vec3{0, 0, 1},
			want: m32.Vec3{0, 1, 0},
		},
	}

	for _, tt := range tests {
		n := ray.CalcNormal(tt.a, tt.b, tt.c)
		if n != tt.want {
			t.Errorf("\nwant: %v\n got: %v\n", tt.want, n)
		}
	}
}
