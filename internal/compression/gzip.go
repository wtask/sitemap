package compression

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Gzip - compress source data into gzip and write to target.
// Returns number of read bytes from origin and error if it occurred.
// Argument `origin` is a data source reader, `gz` - target gzip data writer,
// `header` - optional header for gzip data.
func Gzip(origin io.Reader, gz io.Writer, header *gzip.Header) (read int64, err error) {
	gw := gzip.NewWriter(gz)
	defer func() {
		// flush compressed data and recheck error
		err = gw.Close()
	}()
	if header != nil {
		gw.Header = *header
	}
	read, err = io.Copy(gw, origin)
	if err != nil {
		return read, fmt.Errorf("compression.Gzip failed: %s", err)
	}
	return read, nil
}

// GzipFile - compress source data from origin file with gzip and save into new file.
// Returns number of read bytes from source and error if it occurred.
// Argument `origin` is a file name with source data,
// `gz` is a whished file name of compression result.
func GzipFile(origin, gz string) (read int64, err error) {
	if origin == gz {
		return 0, fmt.Errorf("compression.GzipFile: origin is the as gzip target")
	}
	source, err := os.OpenFile(origin, os.O_RDONLY, 0644)
	if err != nil {
		return 0, fmt.Errorf("compression.GzipFile, cannot open source: %s", err)
	}
	defer source.Close()

	target, err := os.OpenFile(gz, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return 0, fmt.Errorf("compression.GzipFile, cannot open target: %s", err)
	}
	defer target.Close()

	return Gzip(source, target, &gzip.Header{Name: filepath.Base(origin)})
}

// Ungzip - decompress gzip data from `gz` reader and write result into `origin` writer.
// Returns number of read bytes from 'gz' and error if it was occurred.
func Ungzip(gz io.Reader, origin io.Writer) (read int64, err error) {
	gr, err := gzip.NewReader(gz)
	if err != nil {
		return 0, fmt.Errorf("compression.Ungzip, cannot prepare gzip reader: %s", err)
	}
	defer gr.Close()
	read, err = io.Copy(origin, gr)
	if err != nil {
		return read, fmt.Errorf("compression.Ungzip: cannot read gzip data: %s", err)
	}
	return read, nil
}

// UngzipFile - decompress gzip data from `gz` file and writes result into `origin` writer.
// Returns number of read bytes from 'gz' and error if it was occurred.
func UngzipFile(gz string, origin io.Writer) (read int64, err error) {
	source, err := os.OpenFile(gz, os.O_RDONLY, 0644)
	if err != nil {
		return 0, fmt.Errorf("compression.UngzipFile, cannot open source: %s", err)
	}
	defer source.Close()

	return Ungzip(source, origin)
}
