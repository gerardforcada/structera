package helpers

import "testing"

func TestSub(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"1 - 1 = 0", args{1, 1}, 0},
		{"2 - 1 = 1", args{2, 1}, 1},
		{"3 - 1 = 2", args{3, 1}, 2},
		{"-4 - 1 = 5", args{-4, 1}, -5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Sub(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Sub() = %v, want %v", got, tt.want)
			}
		})
	}
}
