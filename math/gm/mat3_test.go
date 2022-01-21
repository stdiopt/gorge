package gm_test

import (
	"testing"

	"github.com/stdiopt/gorge/math/gm"
)

func TestMat3(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want any
	}{
		{
			name: "M3Ident",
			in:   gm.M3Ident(),
			want: gm.Mat3{1, 0, 0, 0, 1, 0, 0, 0, 1},
		},
		{
			name: "Mat4",
			in: gm.Mat3{0, 1, 2, 3, 4, 5, 6, 7, 8}.
				Mat4(),
			want: gm.Mat4{
				0, 1, 2, 0,
				3, 4, 5, 0,
				6, 7, 8, 0,
				0, 0, 0, 1,
			},
		},
		{
			name: "Diag",
			in: gm.Mat3{0, 1, 2, 3, 4, 5, 6, 7, 8}.
				Diag(),
			want: gm.Vec3{0, 4, 8},
		},
		{
			name: "Add",
			in: gm.Mat3{0, 1, 2, 3, 4, 5, 6, 7, 8}.
				Add(gm.Mat3{1, 1, 1, 1, 1, 1, 1, 1, 1}),
			want: gm.Mat3{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name: "Sub",
			in: gm.Mat3{0, 1, 2, 3, 4, 5, 6, 7, 8}.
				Sub(gm.Mat3{1, 1, 1, 1, 1, 1, 1, 1, 1}),
			want: gm.Mat3{-1, 0, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "MulS",
			in: gm.Mat3{0, 1, 2, 3, 4, 5, 6, 7, 8}.
				MulS(2),
			want: gm.Mat3{0, 2, 4, 6, 8, 10, 12, 14, 16},
		},
		{
			name: "Mul",
			in: gm.Mat3{0, 1, 2, 3, 4, 5, 6, 7, 8}.
				Mul(gm.Mat3{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			want: gm.Mat3{24, 30, 36, 51, 66, 81, 78, 102, 126},
		},
		{
			name: "Transpose",
			in: gm.Mat3{0, 1, 2, 3, 4, 5, 6, 7, 8}.
				Transpose(),
			want: gm.Mat3{0, 3, 6, 1, 4, 7, 2, 5, 8},
		},
		{
			name: "Det",
			in: gm.Mat3{1, 0, 0, 0, 1, 0, 0, 0, 1}.
				Det(),
			want: gm.Float(1),
		},
		{
			name: "Inv",
			in: gm.Mat3{2, 0, 0, 0, 2, 0, 0, 0, 2}.
				Inv(),
			want: gm.Mat3{.5, 0, 0, 0, .5, 0, 0, 0, .5},
		},
		{
			name: "ApproxEqual",
			in: gm.Mat3{1, 0, 0, 0, 1, 0, 0, 0, 1}.
				ApproxEqual(gm.Mat3{1 + gm.Epsilon, 0, 0, 0, gm.Epsilon + 1, 0, 0, 0, 1}),
			want: true,
		},
		{
			name: "ApproxEqualThreshold",
			in: gm.Mat3{0, 1, 2, 3, 4, 5, 6, 7, 8}.
				ApproxEqualThreshold(gm.Mat3{1, 2, 3, 4, 5, 6, 7, 8, 9}, 2),
			want: true,
		},
		{
			name: "Abs",
			in: gm.Mat3{0, -1, 2, -3, -4, -5, 6, -7, 8}.
				Abs(),
			want: gm.Mat3{0, 1, 2, 3, 4, 5, 6, 7, 8},
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
