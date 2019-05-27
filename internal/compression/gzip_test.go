package compression

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGzipUngzip(test *testing.T) {
	sourceFile := filepath.Join("testdata", "sitemap.xml")
	expected, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		test.Fatal("unable to load source file", err)
	}

	gzipFile := sourceFile + ".gzip"
	read, err := GzipFile(sourceFile, gzipFile)
	defer os.Remove(gzipFile) // ignore error

	if err != nil {
		test.Error(err)
	} else if read == 0 {
		test.Error("Read zero bytes from non-empty source without error")
	}

	actual := bytes.Buffer{}
	read, err = UngzipFile(gzipFile, &actual)
	if err != nil {
		test.Error(err)
	} else if read == 0 {
		test.Error("Read zero bytes from non-empty gzip without error")
	}

	if !reflect.DeepEqual(expected, actual.Bytes()) {
		test.Error("Unexpected Gzip/Ungzip results")
	}
}
