package m32_test

import (
	"testing"

	"github.com/stdiopt/gorge/m32"
)

func TestQuat(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want any
	}{

		{
			name: "Add",
			in:   m32.Quat{1, 2, 3, 4}.Add(m32.Quat{1, 1, 1, 1}),
			want: m32.Quat{2, 3, 4, 5},
		},
		{
			name: "Normalize",
			in:   m32.Quat{5, 4, 3, 2}.Normalize(),
			want: m32.Quat{0.68041384, 0.5443311, 0.4082483, 0.27216554},
		},
		{
			name: "Mul",
			in:   m32.Quat{1, 1, 1, 1}.Mul(m32.Quat{2, 0, 0, 1}),
			want: m32.Quat{3, 3, -1, -1},
		},
		{
			name: "Mat4",
			in:   m32.Quat{1, 1, 1, 1}.Mat4(),
			want: m32.Mat4{
				-3, 4, 0, 0,
				0, -3, 4, 0,
				4, 0, -3, 0,
				0, 0, 0, 1,
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
