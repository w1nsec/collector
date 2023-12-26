package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
)

func Compress(data []byte) ([]byte, error) {
	var writer bytes.Buffer
	//gzipWr, err := gzip.NewWriterLevel(&writer, gzip.BestSpeed)
	gzipWr, err := gzip.NewWriterLevel(&writer, gzip.DefaultCompression)
	if err != nil {
		return nil, err
	}

	_, err = gzipWr.Write(data)
	if err != nil {
		return nil, err
	}
	err = gzipWr.Close()
	if err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}

func Decompress(data []byte) ([]byte, error) {
	reader := bytes.NewBuffer(data)
	gz, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	defer gz.Close()
	data, err = io.ReadAll(gz)
	if err != nil {
		return nil, err
	}
	return data, nil
}
