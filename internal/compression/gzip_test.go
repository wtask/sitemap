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
	err = GzipFile(sourceFile, gzipFile)
	defer os.Remove(gzipFile) // ignore error

	if err != nil {
		test.Error(err)
	}
	stat, err := os.Stat(gzipFile)
	if err != nil {
		test.Error(err)
	} else if stat.Size() == 0 {
		test.Error("Target gzip has zero size for non-empty source without error")
	}	

	actual := bytes.Buffer{}
	err = UngzipFile(gzipFile, &actual)
	if err != nil {
		test.Error(err)
	} else if actual.Len() == 0 {
		test.Error("Buffer has zero size after for non-empty gzip without error")
	}

	if !reflect.DeepEqual(expected, actual.Bytes()) {
		test.Error("Unexpected Gzip/Ungzip results")
	}
}
