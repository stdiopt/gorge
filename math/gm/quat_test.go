package gm_test

import (
	"testing"

	"github.com/stdiopt/gorge/math/gm"
)

func TestQuat(t *testing.T) {
	tests := []struct {
		name string
		in   any
		want any
	}{

		{
			name: "Add",
			in:   gm.Quat{1, 2, 3, 4}.Add(gm.Quat{1, 1, 1, 1}),
			want: gm.Quat{2, 3, 4, 5},
		},
		{
			name: "Normalize",
			in:   gm.Quat{5, 4, 3, 2}.Normalize(),
			want: gm.Quat{0.68041384, 0.5443311, 0.4082483, 0.27216554},
		},
		{
			name: "Mul",
			in:   gm.Quat{1, 1, 1, 1}.Mul(gm.Quat{2, 0, 0, 1}),
			want: gm.Quat{3, 3, -1, -1},
		},
		{
			name: "Mat4",
			in:   gm.Quat{1, 1, 1, 1}.Mat4(),
			want: gm.Mat4{
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
