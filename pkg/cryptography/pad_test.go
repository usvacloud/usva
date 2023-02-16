package cryptography

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestPadding(t *testing.T) {
	type args struct {
		in        []byte
		blocksize int
	}
	tests := []struct {
		args args
		want []byte
	}{
		{
			args: args{
				in:        []byte{},
				blocksize: 16,
			},
			want: bytes.Repeat([]byte{16}, 16),
		},
		{
			args: args{
				in:        []byte{1, 2, 3, 4},
				blocksize: 16,
			},
			want: append([]byte{1, 2, 3, 4}, bytes.Repeat([]byte{12}, 12)...),
		},
		{
			args: args{
				in:        []byte{1, 2, 4, 1, 2, 3, 4, 1, 2, 3, 4},
				blocksize: 16,
			},
			want: append([]byte{1, 2, 4, 1, 2, 3, 4, 1, 2, 3, 4}, bytes.Repeat([]byte{byte(5)}, 5)...),
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("run %d", i), func(t *testing.T) {
			got := pad(tt.args.in, tt.args.blocksize)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pad() = %v, want %v", got, tt.want)
			} else if f := safeUnpad(got, tt.args.blocksize); !reflect.DeepEqual(f, tt.args.in) {
				t.Errorf("Unpad() = %v, want %v", f, tt.args.in)
			}
		})
	}
}

func TestUnpad(t *testing.T) {
	type args struct {
		in        []byte
		blocksize int
	}
	tests := []struct {
		args args
		want []byte
	}{
		{
			args: args{
				in:        bytes.Repeat([]byte{byte(16)}, 16),
				blocksize: 16,
			},
			want: []byte{},
		},
		{
			args: args{
				in:        append([]byte{1, 2, 3}, bytes.Repeat([]byte{byte(13)}, 13)...),
				blocksize: 16,
			},
			want: []byte{1, 2, 3},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("run %d", i), func(t *testing.T) {
			if got := safeUnpad(tt.args.in, tt.args.blocksize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unpad() = %v, want %v", got, tt.want)
			}
		})
	}
}
