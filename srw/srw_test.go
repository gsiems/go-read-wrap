package srw

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestMulti(t *testing.T) {
	for _, readSz := range []int{500, 922, 1000, 2000} {
		label := fmt.Sprintf("TestMulti[%d]", readSz)
		fmt.Printf("%s...\n", label)
		readers := loadFiles(label, t)
		x := MultiReader(readers...)
		testReaders(label, readSz, x, t)
	}
}

func TestSingleBuff(t *testing.T) {
	for _, readSz := range []int{500, 922, 1000, 2000} {
		for _, buffSz := range []int{0, 90, 500, 922, 1000, 1024, 6000} {
			label := fmt.Sprintf("TestSingleBuff[%d, %d]", readSz, buffSz)
			fmt.Printf("%s...\n", label)

			filename := catDir([]string{"testdata", "t1_expected"})

			r, err := os.Open(filename)
			if err != nil {
				t.Errorf("%s failed opening %q", label, filename)
			}

			x := BuffReader(buffSz, r)
			testReaders(label, readSz, x, t)
		}
	}
}

func TestMultiBuff(t *testing.T) {
	for _, readSz := range []int{500, 922, 1000, 2000} {
		for _, buffSz := range []int{0, 90, 500, 922, 1000, 1024, 6000} {
			label := fmt.Sprintf("TestMultiBuff[%d, %d]", readSz, buffSz)
			fmt.Printf("%s...\n", label)
			readers := loadFiles(label, t)
			x := BuffMultiReader(buffSz, readers...)
			testReaders(label, readSz, x, t)
		}
	}
}

////////////////////////////////////////////////////////////////////////
func testReaders(label string, readSz int, x io.Reader, t *testing.T) {
	expected := getExpected(label, t)

	buff := make([]byte, readSz)
	var inDat []byte

	for {
		n, err := x.Read(buff)

		if err == io.EOF {
			if n > 0 {
				inDat = append(inDat, buff[:n]...)
			}
			break
		} else if err != nil {
			t.Errorf("%s failed reading: %q", label, err)
			if n > 0 {
				inDat = append(inDat, buff[:n]...)
			}
			break
		} else if n != readSz {
			t.Errorf("%s failed reading: expected %d bytes, received %d bytes", label, readSz, n)
			inDat = append(inDat, buff[:n]...)
		} else if readSz == 0 && n == 0 {
			break
		} else {
			inDat = append(inDat, buff...)
		}
	}
	checkExpected(label, expected, inDat, t)
}

func loadFiles(label string, t *testing.T) (readers []io.Reader) {

	inFiles := []string{"t1_01", "t1_02", "t1_03", "t1_04", "t1_05", "t1_06"}

	for _, f := range inFiles {
		filename := catDir([]string{"testdata", f})
		r, err := os.Open(filename)
		if err != nil {
			t.Errorf("%s failed opening %q", label, filename)
			continue
		}
		readers = append(readers, r)
	}
	return readers
}

func getExpected(label string, t *testing.T) (expected []byte) {

	exFile := catDir([]string{"testdata", "t1_expected"})

	expected, err := ioutil.ReadFile(exFile)
	if err != nil {
		t.Errorf("%s failed reading %q", label, exFile)
	}
	return expected
}

func checkExpected(label string, e, a []byte, t *testing.T) {

	if string(e) != string(a) {

		fmt.Println(string(e))
		fmt.Println()
		fmt.Println(string(a))
		fmt.Println()
		fmt.Println()

		t.Errorf("%s expected does not match read", label)
	}
}

func catDir(t []string) (dir string) {
	dir = strings.Join(t, string(os.PathSeparator))
	return dir
}
