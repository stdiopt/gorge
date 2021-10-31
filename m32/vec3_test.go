package m32_test

import (
	"testing"

	"github.com/stdiopt/gorge/m32"
)

func TestVec3(t *testing.T) {
	tests := []struct {
		name string
		in   interface{}
		want interface{}
	}{
		{
			name: "Add",
			in:   m32.Vec3{1, 1, 1}.Add(m32.Vec3{1, 1, 1}),
			want: m32.Vec3{2, 2, 2},
		},
		{
			name: "Sub",
			in:   m32.Vec3{1, 1, 1}.Sub(m32.Vec3{2, 2, 2}),
			want: m32.Vec3{-1, -1, -1},
		},
		{
			name: "Mul",
			in:   m32.Vec3{2, -2, 2}.Mul(2),
			want: m32.Vec3{4, -4, 4},
		},
		{
			name: "Len",
			in:   m32.Vec3{1, 1, 1}.Len(),
			want: float32(1.7320508),
		},
		{
			name: "Normalize",
			in:   m32.Vec3{5, 4, 3}.Normalize(),
			want: m32.Vec3{0.7071068, 0.56568545, 0.42426407},
		},
		{
			name: "Cross",
			in:   m32.Vec3{1, 0, 0}.Cross(m32.Vec3{0, 0, 1}),
			want: m32.Vec3{0, -1, 0},
		},
		{
			name: "Dot",
			in:   m32.Vec3{1, 2, 3}.Dot(m32.Vec3{3, 2, 1}),
			want: float32(10),
		},
		{
			name: "Vec4",
			in:   m32.Vec3{1, 0, 0}.Vec4(1),
			want: m32.Vec4{1, 0, 0, 1},
		},
		{
			name: "V3Lerp",
			in:   m32.V3Lerp(m32.Vec3{0, 0, 0}, m32.Vec3{1, 1, 1}, .5),
			want: m32.Vec3{.5, .5, .5},
		},
		{
			name: "V3Clamp",
			in:   m32.Vec3{-2, 2, 1}.Clamp(0, 1),
			want: m32.Vec3{0, 1, 1},
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
