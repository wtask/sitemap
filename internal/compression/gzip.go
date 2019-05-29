package compression

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Gzip - compress source data into gzip and write to target.
// Argument `origin` is a data source reader, `gz` - target gzip data writer,
// `header` - optional header for gzip data.
func Gzip(origin io.Reader, gz io.Writer, header *gzip.Header) (err error) {
	gw := gzip.NewWriter(gz)
	defer func() {
		// flush compressed data and recheck error
		err = gw.Close()
	}()
	if header != nil {
		gw.Header = *header
	}
	_, err = io.Copy(gw, origin)
	if err != nil {
		return fmt.Errorf("compression.Gzip failed: %s", err)
	}
	return nil
}

// GzipFile - compress source data from origin file with gzip and save into new file.
// Argument `origin` is a file name with source data,
// `gz` is a whished file name of compression result.
func GzipFile(origin, gz string) error {
	if origin == gz {
		return fmt.Errorf("compression.GzipFile: origin is the as gzip target")
	}
	source, err := os.OpenFile(origin, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("compression.GzipFile, cannot open source: %s", err)
	}
	defer source.Close()

	target, err := os.OpenFile(gz, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("compression.GzipFile, cannot open target: %s", err)
	}
	defer target.Close()

	return Gzip(source, target, &gzip.Header{Name: filepath.Base(origin)})
}

// Ungzip - decompress gzip data from `gz` reader and write result into `origin` writer.
func Ungzip(gz io.Reader, origin io.Writer) error {
	gr, err := gzip.NewReader(gz)
	if err != nil {
		return fmt.Errorf("compression.Ungzip, cannot prepare gzip reader: %s", err)
	}
	defer gr.Close()
	_, err = io.Copy(origin, gr)
	if err != nil {
		return fmt.Errorf("compression.Ungzip: cannot read gzip data: %s", err)
	}
	return nil
}

// UngzipFile - decompress gzip data from `gz` file and writes result into `origin` writer.
func UngzipFile(gz string, origin io.Writer) error {
	source, err := os.OpenFile(gz, os.O_RDONLY, 0644)
	if err != nil {
		return fmt.Errorf("compression.UngzipFile, cannot open source: %s", err)
	}
	defer source.Close()

	return Ungzip(source, origin)
}
