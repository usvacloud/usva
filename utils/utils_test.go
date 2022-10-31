package utils

import (
	"reflect"
	"testing"
)

func TestStringOr(t *testing.T) {
	type args struct {
		str1 string
		str2 string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{
				str1: "a",
				str2: "b",
			},
			want: "a",
		},
		{
			name: "test2",
			args: args{
				str1: "a",
				str2: "",
			},
			want: "a",
		},
		{
			name: "test3",
			args: args{
				str1: "",
				str2: "b",
			},
			want: "b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringOr(tt.args.str1, tt.args.str2); got != tt.want {
				t.Errorf("StringOr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntOr(t *testing.T) {
	type args struct {
		int1 int
		int2 int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test1",
			args: args{
				int1: 0,
				int2: 1,
			},
			want: 1,
		},
		{
			name: "test2",
			args: args{
				int1: -50,
				int2: 1,
			},
			want: -50,
		},
		{
			name: "test3",
			args: args{
				int1: 0,
				int2: -10,
			},
			want: -10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntOr(tt.args.int1, tt.args.int2); got != tt.want {
				t.Errorf("IntOr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVarOr(t *testing.T) {
	type args struct {
		var1 any
		var2 any
	}
	tests := []struct {
		name string
		args args
		want any
	}{

		{
			name: "test1",
			args: args{
				var1: nil,
				var2: nil,
			},
			want: nil,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VarOr(tt.args.var1, tt.args.var2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VarOr() = %v, want %v", got, tt.want)
			}
		})
	}
}
