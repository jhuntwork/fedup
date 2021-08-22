package fedup

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// b3sum of testdata/sample.
	goodB3Sum = "05c9ecf478842d332ce3f4a860f4f92f2df040ca5af670b931b007cc2e537862"
)

var (
	errRead  = errors.New("this is a mock Read failure")
	errClose = errors.New("this is a mock Close failure")
)

type badReader struct{}

func (*badReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("%w", errRead)
}

func (*badReader) Close() error {
	return fmt.Errorf("%w", errClose)
}

func Test_computeB3Sum(t *testing.T) {
	t.Parallel()
	computeB3SumTests := []struct {
		description string
		shouldErr   bool
		filename    string
		reader      io.Reader
		expected    string
		errMsg      string
	}{
		{
			description: "should work on typical files",
			filename:    "testdata/sample",
			expected:    goodB3Sum,
		},
		{
			description: "should fail when cannot read from file",
			shouldErr:   true,
			errMsg:      "this is a mock Read failure",
		},
	}
	for _, test := range computeB3SumTests {
		test := test
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)
			var err error
			if test.filename != "" {
				test.reader, err = os.Open(test.filename)
				if err != nil {
					t.Errorf("Unable to open %s", test.filename)
				}
			} else {
				test.reader = &badReader{}
			}
			sum, err := computeB3Sum(test.reader)
			if test.shouldErr {
				assert.EqualError(err, test.errMsg)
			} else {
				assert.NoError(err)
				assert.Equal(test.expected, sum)
			}
		})
	}
}

func Test_computeB3SumFromFile(t *testing.T) {
	t.Parallel()
	t.Run("should not fail if given a valid file", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		output, err := computeB3SumFromFile("testdata/sample")
		assert.Equal(goodB3Sum, output)
		assert.NoError(err)
	})
	t.Run("should fail if given a bad file", func(t *testing.T) {
		t.Parallel()
		assert := assert.New(t)
		output, err := computeB3SumFromFile("testdata/no_such_file")
		assert.Equal("", output)
		assert.EqualError(err, "open testdata/no_such_file: no such file or directory")
	})
}
