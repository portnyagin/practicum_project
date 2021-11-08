package service

import (
	"testing"
)

func TestCheckOrderNum(t *testing.T) {
	type args struct {
		orderNum string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "OrderService. checkOrderNum. Case 1.",
			args: args{orderNum: "1213sdf45678"},
			want: false,
		},
		{
			name: "OrderService. checkOrderNum. Case 1.",
			args: args{orderNum: "4561261212345467"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckOrderNum(tt.args.orderNum); got != tt.want {
				t.Errorf("checkOrderNum() = %v, want %v", got, tt.want)
			}
		})
	}
}
