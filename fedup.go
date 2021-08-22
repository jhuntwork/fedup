package fedup

import (
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/zeebo/blake3"
)

func computeB3Sum(f io.Reader) (string, error) {
	var buf []byte
	hash := blake3.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", fmt.Errorf("%w", err)
	}
	sum := hash.Sum(buf)
	return hex.EncodeToString(sum), nil
}

func computeB3SumFromFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	defer f.Close()
	return computeB3Sum(f)
}

// Dedup will de-duplicate all files in a given directory by creating hard-links
// of files with identical b3sums.
func Dedup(dirname string, dryrun bool, out io.Writer) (int, error) {
	sums := make(map[string]string, 1)
	var count int
	err := filepath.WalkDir(dirname, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.Type().IsRegular() {
			sum, err := computeB3SumFromFile(path)
			if err != nil {
				return err
			}
			// Check if this sum has already been found
			if firstFile, ok := sums[sum]; ok {
				count++
				fmt.Fprintf(out, "Linking %s to %s\n", path, firstFile)
				if dryrun {
					return nil
				}
				// Remove the duplicate first
				if err = os.Remove(path); err != nil {
					return err
				}
				if err = os.Link(firstFile, path); err != nil {
					return err
				}
				return nil
			}
			sums[sum] = path
		}
		return nil
	})
	return count, err
}
