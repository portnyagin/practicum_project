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
			name: "OrderService. checkOrderNum. Case 1. bad symbol",
			args: args{orderNum: "1213sdf45678"},
			want: false,
		},
		{
			name: "OrderService. checkOrderNum. Case 2. Correct. Even number of characters",
			args: args{orderNum: "4561261212345467"},
			want: true,
		},
		{
			name: "OrderService. checkOrderNum. Case 3. Correct. Odd number of characters",
			args: args{orderNum: "8841524506523"},
			want: true,
		},
		{
			name: "OrderService. checkOrderNum. Case 4. Incorrect. Even number of characters",
			args: args{orderNum: "777777"},
			want: false,
		},
		{
			name: "OrderService. checkOrderNum. Case 5. Inorrect. Odd number of characters",
			args: args{orderNum: "55555"},
			want: false,
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
