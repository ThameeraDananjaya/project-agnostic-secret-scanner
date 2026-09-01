package gitinput

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"unicode/utf8"

	"github.com/mholt/archives"
)

// RawClass is authoritative only for exact, pre-preparation bytes.
type RawClass string

const (
	RawText        RawClass = "text"
	RawBinary      RawClass = "binary"
	RawArchive     RawClass = "archive"
	RawCompression RawClass = "compression"
	RawContainer   RawClass = "container"
	RawAmbiguous   RawClass = "ambiguous"
)

var ErrUnsupportedRawClass = errors.New("unsupported raw input class")

type rawMatcher struct {
	format archives.Format
	class  RawClass
}

// These are every independently registered archive/compression family in the
// exact github.com/mholt/archives v0.1.2 stack used by Gitleaks v8.30.1.
var rawMatchers = []rawMatcher{
	{archives.SevenZip{}, RawArchive}, {archives.Rar{}, RawArchive},
	{archives.Tar{}, RawArchive}, {archives.Zip{}, RawArchive},
	{archives.Brotli{}, RawCompression}, {archives.Bz2{}, RawCompression},
	{archives.Gz{}, RawCompression}, {archives.Lz4{}, RawCompression},
	{archives.MinLZ{}, RawCompression}, {archives.Lzip{}, RawCompression},
	{archives.Sz{}, RawCompression}, {archives.Xz{}, RawCompression},
	{archives.Zlib{}, RawCompression}, {archives.Zstd{}, RawCompression},
}

var containerHeaders = [][]byte{
	{0x7f, 'E', 'L', 'F'}, {'M', 'Z'}, {'%', 'P', 'D', 'F', '-'},
	{'S', 'Q', 'L', 'i', 't', 'e', ' ', 'f', 'o', 'r', 'm', 'a', 't', ' ', '3', 0},
}

// ClassifyRawFile runs before any framing. Any family the pinned detector
// stack recognizes, any conflicting family, or a known executable/document
// container is explicit non-pass. Ordinary text and binary remain admissible.
func ClassifyRawFile(path, logicalPath string) (RawClass, error) {
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return RawAmbiguous, ErrUnsupportedRawClass
	}
	matches := make([]RawClass, 0, 2)
	for _, matcher := range rawMatchers {
		f, openErr := os.Open(path)
		if openErr != nil {
			return RawAmbiguous, ErrUnsupportedRawClass
		}
		result, matchErr := matcher.format.Match(context.Background(), logicalPath, f)
		closeErr := f.Close()
		if matchErr != nil || closeErr != nil {
			return RawAmbiguous, ErrUnsupportedRawClass
		}
		if result.Matched() {
			matches = append(matches, matcher.class)
		}
	}
	if len(matches) > 1 {
		return RawAmbiguous, ErrUnsupportedRawClass
	}
	if len(matches) == 1 {
		return matches[0], ErrUnsupportedRawClass
	}
	f, err := os.Open(path)
	if err != nil {
		return RawAmbiguous, ErrUnsupportedRawClass
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	peek, _ := reader.Peek(16)
	for _, header := range containerHeaders {
		if bytes.HasPrefix(peek, header) {
			return RawContainer, ErrUnsupportedRawClass
		}
	}
	validText := true
	buffer := make([]byte, 64<<10)
	carry := make([]byte, 0, utf8.UTFMax)
	for {
		n, readErr := reader.Read(buffer)
		if n > 0 {
			data := append(append(make([]byte, 0, len(carry)+n), carry...), buffer[:n]...)
			carry = carry[:0]
			for len(data) > 0 {
				if data[0] == 0 {
					validText = false
					break
				}
				if data[0] < utf8.RuneSelf {
					data = data[1:]
					continue
				}
				if !utf8.FullRune(data) && readErr == nil {
					carry = append(carry, data...)
					break
				}
				_, size := utf8.DecodeRune(data)
				if size == 1 {
					validText = false
					break
				}
				data = data[size:]
			}
		}
		if !validText {
			_, _ = io.Copy(io.Discard, reader)
			return RawBinary, nil
		}
		if errors.Is(readErr, io.EOF) {
			if len(carry) != 0 {
				return RawBinary, nil
			}
			return RawText, nil
		}
		if readErr != nil {
			return RawAmbiguous, ErrUnsupportedRawClass
		}
	}
}
