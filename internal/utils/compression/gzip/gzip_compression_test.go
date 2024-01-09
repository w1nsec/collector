package gzip

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestCompressDecompress(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantErr bool
	}{
		{
			name:    "Test Compression/Decompression 1",
			text:    strings.Repeat("compression test", 10),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compress([]byte(tt.text))
			if (err != nil) != tt.wantErr {
				t.Errorf("Compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			result, err := Decompress(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decompress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.text != string(result) {
				t.Errorf("Compress/Decompress error got = %v, want %v", got, tt.text)
			}
		})
	}
}

func TestCompress(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Test Compression",
			data:    strings.Repeat("compression test", 10),
			want:    []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 74, 206, 207, 45, 40, 74, 45, 46, 206, 204, 207, 83, 40, 73, 45, 46, 25, 108, 124, 64, 0, 0, 0, 255, 255, 171, 84, 54, 21, 160, 0, 0, 0},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compress([]byte(tt.data))
			if (err != nil) != tt.wantErr {
				t.Errorf("Compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecompress(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Test Decompression",
			data:    []byte{31, 139, 8, 0, 0, 0, 0, 0, 0, 255, 74, 73, 77, 206, 207, 45, 40, 74, 45, 46, 206, 204, 207, 83, 40, 73, 45, 46, 25, 252, 34, 128, 0, 0, 0, 255, 255, 171, 212, 72, 153, 180, 0, 0, 0},
			want:    strings.Repeat("decompression test", 10),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decompress(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(got) != tt.want {
				t.Errorf("Compress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleCompress() {
	data := strings.Repeat("test data", 10)
	fmt.Println("Text len:          ", len(data))

	dataComp, err := Compress([]byte(data))
	if err != nil {
		fmt.Printf("can't compress data: %v\n", err)
		return
	}
	fmt.Println("Compressed len:    ", len(dataComp))

	dataDecomp, err := Decompress(dataComp)
	if err != nil {
		fmt.Printf("can't decompress data: %v\n", err)
		return
	}
	fmt.Println("Decompressed len:  ", len(dataDecomp))
	fmt.Println("Equal:             ", data == string(dataDecomp))

	// Output:
	// Text len:           90
	// Compressed len:     35
	// Decompressed len:   90
	// Equal:              true
}
