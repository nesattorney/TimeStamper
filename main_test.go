package main

import "testing"

func Test_convertSeconds(t *testing.T) {
	type args struct {
		seconds int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			"test1",
			args{
				150,
			},
			"00:02:30",
		},
		{
			"test2",
			args{
				10,
			},
			"00:00:10",
		},
		{
			"test3",
			args{
				0,
			},
			"00:00:00",
		},
		{
			"test1",
			args{
				-10,
			},
			"00:00:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertSeconds(tt.args.seconds); got != tt.want {
				t.Errorf("convertSeconds() = %v, want %v", got, tt.want)
			}
		})
	}
}
