package gm_test

import (
	"testing"

	"github.com/stdiopt/gorge/math/gm"
)

func TestVec3(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want any
	}{
		{
			name: "Add",
			in:   gm.Vec3{1, 1, 1}.Add(gm.Vec3{1, 1, 1}),
			want: gm.Vec3{2, 2, 2},
		},
		{
			name: "Sub",
			in:   gm.Vec3{1, 1, 1}.Sub(gm.Vec3{2, 2, 2}),
			want: gm.Vec3{-1, -1, -1},
		},
		{
			name: "Mul",
			in:   gm.Vec3{2, -2, 2}.Mul(2),
			want: gm.Vec3{4, -4, 4},
		},
		{
			name: "Len",
			in:   gm.Vec3{1, 1, 1}.Len(),
			want: gm.Float(1.7320508),
		},
		{
			name: "Normalize",
			in:   gm.Vec3{5, 4, 3}.Normalize(),
			want: gm.Vec3{0.7071068, 0.56568545, 0.42426407},
		},
		{
			name: "Cross",
			in:   gm.Vec3{1, 0, 0}.Cross(gm.Vec3{0, 0, 1}),
			want: gm.Vec3{0, -1, 0},
		},
		{
			name: "Dot",
			in:   gm.Vec3{1, 2, 3}.Dot(gm.Vec3{3, 2, 1}),
			want: gm.Float(10),
		},
		{
			name: "Vec4",
			in:   gm.Vec3{1, 0, 0}.Vec4(1),
			want: gm.Vec4{1, 0, 0, 1},
		},
		{
			name: "V3Lerp",
			in:   gm.Vec3{0, 0, 0}.Lerp(gm.Vec3{1, 1, 1}, .5),
			want: gm.Vec3{.5, .5, .5},
		},
		{
			name: "V3Clamp",
			in:   gm.Vec3{-2, 2, 1}.Clamp(gm.Vec3{}, gm.Vec3{1, 1, 1}),
			want: gm.Vec3{0, 1, 1},
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
