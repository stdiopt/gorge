package gm_test

import (
	"testing"

	"github.com/stdiopt/gorge/math/gm"
)

func TestVec4(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want any
	}{
		{
			name: "Add",
			in:   gm.Vec4{1, 1, 1, 1}.Add(gm.Vec4{1, 1, 1, 1}),
			want: gm.Vec4{2, 2, 2, 2},
		},
		{
			name: "Sub",
			in:   gm.Vec4{1, 1, 1, 1}.Sub(gm.Vec4{2, 2, 2, 2}),
			want: gm.Vec4{-1, -1, -1, -1},
		},
		{
			name: "Mul",
			in:   gm.Vec4{2, -2, 2, -2}.Mul(2),
			want: gm.Vec4{4, -4, 4, -4},
		},
		{
			name: "Len",
			in:   gm.Vec4{1, 1, 1, 1}.Len(),
			want: gm.Float(2),
		},
		{
			name: "Normalize",
			in:   gm.Vec4{5, 4, 3, 2}.Normalize(),
			want: gm.Vec4{0.68041384, 0.5443311, 0.4082483, 0.27216554},
		},
		{
			name: "Dot",
			in:   gm.Vec4{1, 2, 3, 4}.Dot(gm.Vec4{4, 3, 2, 1}),
			want: gm.Float(20),
		},
		{
			name: "Vec3",
			in:   gm.Vec4{1, 2, 3, 4}.Vec3(),
			want: gm.Vec3{1, 2, 3},
		},
		{
			name: "Vec2",
			in:   gm.Vec4{1, 2, 3, 4}.Vec2(),
			want: gm.Vec2{1, 2},
		},

		{
			name: "V4Lerp",
			in:   gm.Vec4{0, 0, 0, 0}.Lerp(gm.Vec4{1, 1, 1, 1}, .5),
			want: gm.Vec4{.5, .5, .5, .5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.in
			if want := tt.want; got != want {
				t.Errorf("\nwant: %v\n got: %v\n", want, got)
			}
		})
	}
}
