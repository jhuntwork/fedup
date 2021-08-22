package fedup_test

import (
	"bytes"
	"testing"

	"github.com/jhuntwork/fedup"
	"github.com/stretchr/testify/assert"
)

func TestDedup(t *testing.T) {
	t.Parallel()
	tests := []struct {
		description string
		directory   string
		errMsg      string
	}{
		{
			description: "Should fail when provided directory doesn't exist",
			directory:   "testdata/no_such_file",
			errMsg:      "testdata/no_such_file: no such file or directory",
		},
		{
			description: "Should not fail when given a valid directory",
			directory:   "testdata",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)
			_, err := fedup.Dedup(tc.directory, true, &bytes.Buffer{})
			if tc.errMsg != "" {
				assert.Contains(err.Error(), tc.errMsg)
			} else {
				assert.NoError(err)
			}
		})
	}
}
