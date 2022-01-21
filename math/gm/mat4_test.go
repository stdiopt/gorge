package gm_test

import (
	"testing"

	"github.com/stdiopt/gorge/math/gm"
)

func TestMat4(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want any
	}{
		{
			name: "M3Ident",
			in:   gm.M4Ident(),
			want: gm.Mat4{
				1, 0, 0, 0,
				0, 1, 0, 0,
				0, 0, 1, 0,
				0, 0, 0, 1,
			},
		},
		{
			name: "Mat3",
			in:   mat4Test().Mat3(),
			want: gm.Mat3{
				0, 1, 2,
				4, 5, 6,
				8, 9, 10,
			},
		},
		{
			name: "Diag",
			in:   mat4Test().Diag(),
			want: gm.Vec4{0, 5, 10, 15},
		},
		{
			name: "Add",
			in: mat4Test().Add(gm.Mat4{
				1, 1, 1, 1,
				1, 1, 1, 1,
				1, 1, 1, 1,
				1, 1, 1, 1,
			}),
			want: gm.Mat4{
				1, 2, 3, 4,
				5, 6, 7, 8,
				9, 10, 11, 12,
				13, 14, 15, 16,
			},
		},
		{
			name: "Sub",
			in: gm.Mat4{
				0, 1, 2, 3,
				4, 5, 6, 7,
				8, 9, 10, 11,
				12, 13, 14, 15,
			}.Sub(gm.Mat4{
				1, 1, 1, 1,
				1, 1, 1, 1,
				1, 1, 1, 1,
				1, 1, 1, 1,
			}),
			want: gm.Mat4{
				-1, 0, 1, 2,
				3, 4, 5, 6,
				7, 8, 9, 10,
				11, 12, 13, 14,
			},
		},
		{
			name: "MulS",
			in:   mat4Test().MulS(2),
			want: gm.Mat4{
				0, 2, 4, 6,
				8, 10, 12, 14,
				16, 18, 20, 22,
				24, 26, 28, 30,
			},
		},
		{
			name: "Mul",
			in:   mat4Test().Mul(mat4Test()),
			want: gm.Mat4{
				56, 62, 68, 74,
				152, 174, 196, 218,
				248, 286, 324, 362,
				344, 398, 452, 506,
			},
		},
		{
			name: "Transpose",
			in:   mat4Test().Transpose(),
			want: gm.Mat4{
				0, 4, 8, 12,
				1, 5, 9, 13,
				2, 6, 10, 14,
				3, 7, 11, 15,
			},
		},
		{
			name: "Det",
			in: gm.Mat4{
				1, 0, 0, 0,
				0, 1, 0, 0,
				0, 0, 1, 0,
				0, 0, 0, 1,
			}.Det(),
			want: gm.Float(1),
		},
		{
			name: "Inv",
			in: gm.Mat4{
				2, 0, 0, 0,
				0, 2, 0, 0,
				0, 0, 2, 0,
				0, 0, 0, 2,
			}.Inv(),
			want: gm.Mat4{
				.5, 0, 0, 0,
				0, .5, 0, 0,
				0, 0, .5, 0,
				0, 0, 0, .5,
			},
		},
		{
			name: "ApproxEqual",
			in: gm.Mat4{1, 0, 0, 0, 1, 0, 0, 0, 1}.
				ApproxEqual(gm.Mat4{1 + gm.Epsilon, 0, 0, 0, gm.Epsilon + 1, 0, 0, 0, 1}),
			want: true,
		},
		{
			name: "ApproxEqualThreshold",
			in: gm.Mat4{0, 1, 2, 3, 4, 5, 6, 7, 8}.
				ApproxEqualThreshold(gm.Mat4{1, 2, 3, 4, 5, 6, 7, 8, 9}, 2),
			want: true,
		},
		{
			name: "Abs",
			in: gm.Mat4{0, -1, 2, -3, -4, -5, 6, -7, 8}.
				Abs(),
			want: gm.Mat4{0, 1, 2, 3, 4, 5, 6, 7, 8},
		},
		{
			name: "MulV4",
			in:   mat4Test().MulV4(gm.Vec4{0, 0, 0, 1}),
			want: gm.Vec4{12, 13, 14, 15},
		},
		{
			name: "Translate3D",
			in:   gm.Translate3D(5, 5, 5),
			want: gm.Mat4{
				1, 0, 0, 0,
				0, 1, 0, 0,
				0, 0, 1, 0,
				5, 5, 5, 1,
			},
		},
		{
			name: "Scale3D",
			in:   gm.Scale3D(5, 5, 5),
			want: gm.Mat4{
				5, 0, 0, 0,
				0, 5, 0, 0,
				0, 0, 5, 0,
				0, 0, 0, 1,
			},
		},
		{
			name: "LookAt",
			in:   gm.LookAt(gm.Vec3{0, 0, -1}, gm.Vec3{}, gm.Up()),
			want: gm.Mat4{
				-1, 0, 0, 0,
				0, 1, 0, 0,
				0, 0, -1, 0,
				0, 0, -1, 1,
			},
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

func mat4Test() gm.Mat4 {
	return gm.Mat4{
		0, 1, 2, 3,
		4, 5, 6, 7,
		8, 9, 10, 11,
		12, 13, 14, 15,
	}
}
