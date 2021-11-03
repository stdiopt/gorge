package m32_test

import (
	"testing"

	"github.com/stdiopt/gorge/m32"
)

func TestVec4(t *testing.T) {
	tests := []struct {
		name string
		in   interface{}
		want interface{}
	}{
		{
			name: "Add",
			in:   m32.Vec4{1, 1, 1, 1}.Add(m32.Vec4{1, 1, 1, 1}),
			want: m32.Vec4{2, 2, 2, 2},
		},
		{
			name: "Sub",
			in:   m32.Vec4{1, 1, 1, 1}.Sub(m32.Vec4{2, 2, 2, 2}),
			want: m32.Vec4{-1, -1, -1, -1},
		},
		{
			name: "Mul",
			in:   m32.Vec4{2, -2, 2, -2}.Mul(2),
			want: m32.Vec4{4, -4, 4, -4},
		},
		{
			name: "Len",
			in:   m32.Vec4{1, 1, 1, 1}.Len(),
			want: float32(2),
		},
		{
			name: "Normalize",
			in:   m32.Vec4{5, 4, 3, 2}.Normalize(),
			want: m32.Vec4{0.68041384, 0.5443311, 0.4082483, 0.27216554},
		},
		{
			name: "Dot",
			in:   m32.Vec4{1, 2, 3, 4}.Dot(m32.Vec4{4, 3, 2, 1}),
			want: float32(20),
		},
		{
			name: "Vec3",
			in:   m32.Vec4{1, 2, 3, 4}.Vec3(),
			want: m32.Vec3{1, 2, 3},
		},
		{
			name: "Vec2",
			in:   m32.Vec4{1, 2, 3, 4}.Vec2(),
			want: m32.Vec2{1, 2},
		},

		{
			name: "V4Lerp",
			in:   m32.Vec4{0, 0, 0, 0}.Lerp(m32.Vec4{1, 1, 1, 1}, .5),
			want: m32.Vec4{.5, .5, .5, .5},
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
