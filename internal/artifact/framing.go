package artifact

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"os"
	"strconv"
	"strings"
)

type zipEnd struct {
	entries       int
	centralOffset int64
	centralSize   int64
	endOffset     int64
}

func readZIPEnd(path string) (zipEnd, error) {
	f, err := os.Open(path)
	if err != nil {
		return zipEnd{}, Rejection{Code: RejectUnsafe}
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || info.Size() < 22 {
		return zipEnd{}, Rejection{Code: RejectMalformed}
	}
	size := info.Size()
	tailSize := int64(65_557)
	if size < tailSize {
		tailSize = size
	}
	tail := make([]byte, tailSize)
	if _, err := f.ReadAt(tail, size-tailSize); err != nil && err != io.EOF {
		return zipEnd{}, Rejection{Code: RejectMalformed}
	}
	offset := bytes.LastIndex(tail, []byte{'P', 'K', 5, 6})
	if offset < 0 || offset+22 > len(tail) {
		return zipEnd{}, Rejection{Code: RejectMalformed}
	}
	commentLength := int(binary.LittleEndian.Uint16(tail[offset+20 : offset+22]))
	if offset+22+commentLength != len(tail) {
		return zipEnd{}, Rejection{Code: RejectMalformed}
	}
	eocd := tail[offset:]
	disk := binary.LittleEndian.Uint16(eocd[4:6])
	centralDisk := binary.LittleEndian.Uint16(eocd[6:8])
	diskEntries := binary.LittleEndian.Uint16(eocd[8:10])
	totalEntries := binary.LittleEndian.Uint16(eocd[10:12])
	centralSize := binary.LittleEndian.Uint32(eocd[12:16])
	centralOffset := binary.LittleEndian.Uint32(eocd[16:20])
	endOffset := size - tailSize + int64(offset)
	if disk != 0 || centralDisk != 0 || diskEntries != totalEntries || totalEntries == 0 ||
		totalEntries == 0xffff || centralSize == 0xffffffff || centralOffset == 0xffffffff ||
		int64(centralOffset)+int64(centralSize) != endOffset {
		return zipEnd{}, Rejection{Code: RejectUnsupported}
	}
	return zipEnd{entries: int(totalEntries), centralOffset: int64(centralOffset), centralSize: int64(centralSize), endOffset: endOffset}, nil
}

func validateZIPFraming(path string) error {
	_, err := readZIPEnd(path)
	return err
}

func preflightZIP(ctx context.Context, path string, remainingEntries int, remainingBytes int64, limits Limits) error {
	end, err := readZIPEnd(path)
	if err != nil {
		return err
	}
	if end.entries > remainingEntries {
		return Rejection{Code: RejectResource}
	}
	f, err := os.Open(path)
	if err != nil {
		return Rejection{Code: RejectUnsafe}
	}
	defer f.Close()
	if _, err := f.Seek(end.centralOffset, io.SeekStart); err != nil {
		return Rejection{Code: RejectMalformed}
	}
	var total int64
	seen := map[string]bool{}
	position := end.centralOffset
	fixed := make([]byte, 46)
	for index := 0; index < end.entries; index++ {
		if err := ctx.Err(); err != nil {
			return contextRejection(err)
		}
		if _, err := io.ReadFull(f, fixed); err != nil || binary.LittleEndian.Uint32(fixed[:4]) != 0x02014b50 {
			return Rejection{Code: RejectMalformed}
		}
		flags := binary.LittleEndian.Uint16(fixed[8:10])
		method := binary.LittleEndian.Uint16(fixed[10:12])
		compressed := binary.LittleEndian.Uint32(fixed[20:24])
		expanded := binary.LittleEndian.Uint32(fixed[24:28])
		nameLength := int(binary.LittleEndian.Uint16(fixed[28:30]))
		extraLength := int(binary.LittleEndian.Uint16(fixed[30:32]))
		commentLength := int(binary.LittleEndian.Uint16(fixed[32:34]))
		disk := binary.LittleEndian.Uint16(fixed[34:36])
		if flags&1 != 0 {
			return Rejection{Code: RejectUnsafe}
		}
		if method != 0 && method != 8 || compressed == 0xffffffff || expanded == 0xffffffff || disk != 0 || nameLength == 0 {
			return Rejection{Code: RejectUnsupported}
		}
		name := make([]byte, nameLength)
		if _, err := io.ReadFull(f, name); err != nil {
			return Rejection{Code: RejectMalformed}
		}
		logical := strings.TrimSuffix(string(name), "/")
		if !safeArchivePath(logical) || seen[strings.ToLower(logical)] {
			return Rejection{Code: RejectUnsafe}
		}
		seen[strings.ToLower(logical)] = true
		if int64(expanded) > limits.MaxFileBytes || exceedsRatio(int64(expanded), int64(compressed), limits.MaxCompressionRatio) {
			return Rejection{Code: RejectResource}
		}
		if int64(expanded) > remainingBytes-total {
			return Rejection{Code: RejectResource}
		}
		total += int64(expanded)
		skip := int64(extraLength + commentLength)
		if _, err := f.Seek(skip, io.SeekCurrent); err != nil {
			return Rejection{Code: RejectMalformed}
		}
		position += 46 + int64(nameLength) + skip
	}
	if position != end.centralOffset+end.centralSize {
		return Rejection{Code: RejectMalformed}
	}
	return nil
}

func validateTARFile(ctx context.Context, path string) error {
	_, err := preflightTAR(ctx, path, int(^uint(0)>>1), int64(^uint64(0)>>1), int64(^uint64(0)>>1))
	return err
}

type tarPreflight struct {
	entries  int
	expanded int64
}

func preflightTAR(ctx context.Context, path string, remainingEntries int, remainingBytes, maxFile int64) (tarPreflight, error) {
	f, err := os.Open(path)
	if err != nil {
		return tarPreflight{}, Rejection{Code: RejectUnsafe}
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil || info.Size() < 1536 || info.Size()%512 != 0 {
		return tarPreflight{}, Rejection{Code: RejectMalformed}
	}
	header := make([]byte, 512)
	zeroBlocks := 0
	entries := 0
	var expanded int64
	for offset := int64(0); offset+512 <= info.Size(); {
		if err := ctx.Err(); err != nil {
			return tarPreflight{}, contextRejection(err)
		}
		if _, err := io.ReadFull(f, header); err != nil {
			return tarPreflight{}, Rejection{Code: RejectMalformed}
		}
		offset += 512
		if allZero(header) {
			zeroBlocks++
			if zeroBlocks >= 2 {
				for offset < info.Size() {
					if err := ctx.Err(); err != nil {
						return tarPreflight{}, contextRejection(err)
					}
					if _, err := io.ReadFull(f, header); err != nil || !allZero(header) {
						return tarPreflight{}, Rejection{Code: RejectMalformed}
					}
					offset += 512
				}
				return tarPreflight{entries: entries, expanded: expanded}, nil
			}
			continue
		}
		if zeroBlocks != 0 {
			return tarPreflight{}, Rejection{Code: RejectMalformed}
		}
		if !looksLikeTARHeader(header) {
			return tarPreflight{}, Rejection{Code: RejectMalformed}
		}
		size, ok := parseTAROctal(header[124:136])
		if !ok || size < 0 {
			return tarPreflight{}, Rejection{Code: RejectMalformed}
		}
		typeflag := header[156]
		if typeflag == 0 {
			typeflag = '0'
		}
		if typeflag != '0' && typeflag != '5' {
			return tarPreflight{}, Rejection{Code: RejectUnsafe}
		}
		entries++
		if entries > remainingEntries {
			return tarPreflight{}, Rejection{Code: RejectResource}
		}
		if typeflag == '0' {
			if size > maxFile || size > remainingBytes-expanded {
				return tarPreflight{}, Rejection{Code: RejectResource}
			}
			expanded += size
		}
		blocks := (size + 511) / 512
		next := offset + blocks*512
		if next > info.Size() {
			return tarPreflight{}, Rejection{Code: RejectMalformed}
		}
		if blocks > 0 {
			if _, err := f.Seek(blocks*512, io.SeekCurrent); err != nil {
				return tarPreflight{}, Rejection{Code: RejectMalformed}
			}
		}
		offset = next
	}
	return tarPreflight{}, Rejection{Code: RejectMalformed}
}

func looksLikeTARHeader(header []byte) bool {
	if len(header) < 512 || allZero(header[:512]) {
		return false
	}
	want, ok := parseTAROctal(header[148:156])
	if !ok {
		return false
	}
	var sum int64
	for index, value := range header[:512] {
		if index >= 148 && index < 156 {
			sum += int64(' ')
		} else {
			sum += int64(value)
		}
	}
	return sum == want
}

func parseTAROctal(raw []byte) (int64, bool) {
	trimmed := bytes.Trim(raw, " \x00")
	if len(trimmed) == 0 || trimmed[0]&0x80 != 0 {
		return 0, false
	}
	value, err := strconv.ParseInt(string(trimmed), 8, 64)
	return value, err == nil
}

func allZero(raw []byte) bool {
	for _, value := range raw {
		if value != 0 {
			return false
		}
	}
	return true
}
