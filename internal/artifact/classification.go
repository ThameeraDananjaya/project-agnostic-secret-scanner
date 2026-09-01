package artifact

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/h2non/filetype"
	"github.com/mholt/archives"
)

type broadClass uint8

const (
	broadText broadClass = iota
	broadBinary
	broadArchive
	broadCompression
	broadContainer
	broadAmbiguous
)

type broadMatcher struct {
	format archives.Format
	class  broadClass
}

var broadMatchers = []broadMatcher{
	{archives.SevenZip{}, broadArchive}, {archives.Rar{}, broadArchive},
	{archives.Tar{}, broadArchive}, {archives.Zip{}, broadArchive},
	{archives.Brotli{}, broadCompression}, {archives.Bz2{}, broadCompression},
	{archives.Gz{}, broadCompression}, {archives.Lz4{}, broadCompression},
	{archives.MinLZ{}, broadCompression}, {archives.Lzip{}, broadCompression},
	{archives.Sz{}, broadCompression}, {archives.Xz{}, broadCompression},
	{archives.Zlib{}, broadCompression}, {archives.Zstd{}, broadCompression},
}

var broadContainerHeaders = [][]byte{
	{0x7f, 'E', 'L', 'F'}, {'M', 'Z'}, {'%', 'P', 'D', 'F', '-'},
	{'S', 'Q', 'L', 'i', 't', 'e', ' ', 'f', 'o', 'r', 'm', 'a', 't', ' ', '3', 0},
}

type broadMatch struct {
	class     broadClass
	extension string
	byStream  bool
}

func broadClassify(ctx context.Context, path, logicalPath string) (broadClass, error) {
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return broadAmbiguous, Rejection{Code: RejectUnsupported}
	}
	matches := make([]broadMatch, 0, 2)
	for _, matcher := range broadMatchers {
		if err := ctx.Err(); err != nil {
			return broadAmbiguous, contextRejection(err)
		}
		f, err := os.Open(path)
		if err != nil {
			return broadAmbiguous, Rejection{Code: RejectUnsupported}
		}
		result, matchErr := matcher.format.Match(ctx, logicalPath, f)
		closeErr := f.Close()
		if matchErr != nil || closeErr != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return broadAmbiguous, contextRejection(ctxErr)
			}
			return broadAmbiguous, Rejection{Code: RejectUnsupported}
		}
		if result.Matched() {
			matches = append(matches, broadMatch{class: matcher.class, extension: strings.TrimPrefix(matcher.format.Extension(), "."), byStream: result.ByStream})
		}
	}
	f, err := os.Open(path)
	if err != nil {
		return broadAmbiguous, Rejection{Code: RejectUnsupported}
	}
	defer f.Close()
	reader := bufio.NewReader(contextReader{ctx: ctx, r: f})
	peek, peekErr := reader.Peek(261)
	if peekErr != nil && !errors.Is(peekErr, bufio.ErrBufferFull) && !errors.Is(peekErr, io.EOF) {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return broadAmbiguous, contextRejection(ctxErr)
		}
		return broadAmbiguous, Rejection{Code: RejectUnsupported}
	}
	containerMatch, containerExtension := false, ""
	if len(peek) > 0 {
		kind, matchErr := filetype.Match(peek)
		if matchErr != nil {
			return broadAmbiguous, Rejection{Code: RejectUnsupported}
		}
		containerMatch = kind.MIME.Type == "application"
		containerExtension = kind.Extension
	}
	for _, header := range broadContainerHeaders {
		if bytes.HasPrefix(peek, header) {
			containerMatch = true
		}
	}
	if len(matches) > 1 || len(matches) == 1 && containerMatch && (!matches[0].byStream || matches[0].extension != containerExtension) {
		return broadAmbiguous, Rejection{Code: RejectUnsupported}
	}
	if len(matches) == 1 {
		return matches[0].class, Rejection{Code: RejectUnsupported}
	}
	if containerMatch {
		return broadContainer, Rejection{Code: RejectUnsupported}
	}
	buffer := make([]byte, 64<<10)
	carry := make([]byte, 0, utf8.UTFMax)
	for {
		n, readErr := reader.Read(buffer)
		if n > 0 {
			data := append(append(make([]byte, 0, len(carry)+n), carry...), buffer[:n]...)
			carry = carry[:0]
			for len(data) > 0 {
				if data[0] == 0 {
					return broadBinary, nil
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
					return broadBinary, nil
				}
				data = data[size:]
			}
		}
		if errors.Is(readErr, io.EOF) {
			if len(carry) != 0 {
				return broadBinary, nil
			}
			return broadText, nil
		}
		if readErr != nil {
			if ctxErr := ctx.Err(); ctxErr != nil {
				return broadAmbiguous, contextRejection(ctxErr)
			}
			return broadAmbiguous, Rejection{Code: RejectUnsupported}
		}
	}
}
