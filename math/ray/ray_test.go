package ray_test

import (
	"testing"

	"github.com/stdiopt/gorge/math/gm"
	"github.com/stdiopt/gorge/math/ray"
)

func TestCalcNormal(t *testing.T) {
	tests := []struct {
		a, b, c gm.Vec3
		want    gm.Vec3
	}{
		{
			a:    gm.Vec3{0, 0, 0},
			b:    gm.Vec3{1, 0, 0},
			c:    gm.Vec3{0, 1, 0},
			want: gm.Vec3{0, 0, 1},
		},
		{
			a:    gm.Vec3{0, 0, 0},
			b:    gm.Vec3{-1, 0, 0},
			c:    gm.Vec3{0, 0, 1},
			want: gm.Vec3{0, 1, 0},
		},
	}

	for _, tt := range tests {
		n := ray.CalcNormal(tt.a, tt.b, tt.c)
		if n != tt.want {
			t.Errorf("\nwant: %v\n got: %v\n", tt.want, n)
		}
	}
}
