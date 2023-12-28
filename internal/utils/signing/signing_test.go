package signing

import (
	"reflect"
	"testing"
)

func TestCheckSigning(t *testing.T) {
	type args struct {
		data []byte
		sign []byte
		key  []byte
	}
	tests := []struct {
		name  string
		args  args
		equal bool
	}{
		// TODO: Add test cases.
		{
			name: "Test CheckSigning",
			args: args{
				data: []byte("Create new signing for this line"),
				key:  []byte("newsupersecretkey"),
				sign: []byte("4xaGT5JdgMMDIELeoiWNux0jHi8RNP2ozeDR9PahOi4="),
			},
			equal: true,
		},
		{
			name: "Test CheckSigning",
			args: args{
				data: []byte("Create new signing for this line"),
				key:  []byte("newsupersecretkey"),
				sign: nil,
			},
			equal: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compare := CheckSigning(tt.args.data, tt.args.sign, tt.args.key)
			if compare != tt.equal {
				t.Errorf("CheckSigning() = %v, want %v", compare, tt.equal)
			}
		})
	}
}

func TestCreateSigning(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		key  []byte
		want []byte
	}{
		// TODO: Add test cases.
		{
			name: "Test CreateSigning 1",
			data: []byte("Create new signing for this line"),
			key:  []byte("newsupersecretkey"),
			want: []byte{52, 120, 97, 71, 84, 53, 74, 100, 103, 77, 77, 68, 73, 69, 76, 101, 111, 105, 87, 78, 117, 120, 48, 106, 72, 105, 56, 82, 78, 80, 50, 111, 122, 101, 68, 82, 57, 80, 97, 104, 79, 105, 52, 61},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateSigning(tt.data, tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSigning() = %v, want %v", got, tt.want)
			}
		})
	}
}
