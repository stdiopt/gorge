package m32_test

import (
	"testing"

	"github.com/stdiopt/gorge/m32"
)

func TestMat3(t *testing.T) {
	tests := []struct {
		name string
		in   interface{}
		want interface{}
	}{
		{
			name: "M3Ident",
			in:   m32.M3Ident(),
			want: m32.Mat3{1, 0, 0, 0, 1, 0, 0, 0, 1},
		},
		{
			name: "Mat4",
			in: m32.Mat3{0, 1, 2, 3, 4, 5, 6, 7, 8}.
				Mat4(),
			want: m32.Mat4{
				0, 1, 2, 0,
				3, 4, 5, 0,
				6, 7, 8, 0,
				0, 0, 0, 1,
			},
		},
		{
			name: "Diag",
			in: m32.Mat3{0, 1, 2, 3, 4, 5, 6, 7, 8}.
				Diag(),
			want: m32.Vec3{0, 4, 8},
		},
		{
			name: "Add",
			in: m32.Mat3{0, 1, 2, 3, 4, 5, 6, 7, 8}.
				Add(m32.Mat3{1, 1, 1, 1, 1, 1, 1, 1, 1}),
			want: m32.Mat3{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name: "Sub",
			in: m32.Mat3{0, 1, 2, 3, 4, 5, 6, 7, 8}.
				Sub(m32.Mat3{1, 1, 1, 1, 1, 1, 1, 1, 1}),
			want: m32.Mat3{-1, 0, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "MulS",
			in: m32.Mat3{0, 1, 2, 3, 4, 5, 6, 7, 8}.
				MulS(2),
			want: m32.Mat3{0, 2, 4, 6, 8, 10, 12, 14, 16},
		},
		{
			name: "Mul",
			in: m32.Mat3{0, 1, 2, 3, 4, 5, 6, 7, 8}.
				Mul(m32.Mat3{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			want: m32.Mat3{24, 30, 36, 51, 66, 81, 78, 102, 126},
		},
		{
			name: "Transpose",
			in: m32.Mat3{0, 1, 2, 3, 4, 5, 6, 7, 8}.
				Transpose(),
			want: m32.Mat3{0, 3, 6, 1, 4, 7, 2, 5, 8},
		},
		{
			name: "Det",
			in: m32.Mat3{1, 0, 0, 0, 1, 0, 0, 0, 1}.
				Det(),
			want: float32(1),
		},
		{
			name: "Inv",
			in: m32.Mat3{2, 0, 0, 0, 2, 0, 0, 0, 2}.
				Inv(),
			want: m32.Mat3{.5, 0, 0, 0, .5, 0, 0, 0, .5},
		},
		{
			name: "ApproxEqual",
			in: m32.Mat3{1, 0, 0, 0, 1, 0, 0, 0, 1}.
				ApproxEqual(m32.Mat3{1 + m32.Epsilon, 0, 0, 0, m32.Epsilon + 1, 0, 0, 0, 1}),
			want: true,
		},
		{
			name: "ApproxEqualThreshold",
			in: m32.Mat3{0, 1, 2, 3, 4, 5, 6, 7, 8}.
				ApproxEqualThreshold(m32.Mat3{1, 2, 3, 4, 5, 6, 7, 8, 9}, 2),
			want: true,
		},
		{
			name: "Abs",
			in: m32.Mat3{0, -1, 2, -3, -4, -5, 6, -7, 8}.
				Abs(),
			want: m32.Mat3{0, 1, 2, 3, 4, 5, 6, 7, 8},
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
