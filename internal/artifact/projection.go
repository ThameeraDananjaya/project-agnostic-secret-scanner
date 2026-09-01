package artifact

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const projectionHeader = "PSCAN_ARTIFACT_PROJECTION_V1"

func (n *normalizer) project(source string, class Class, depth int) (string, error) {
	info, err := regularNoLink(source)
	if err != nil {
		return "", err
	}
	if info.Size() > n.limits.MaxFileBytes {
		return "", Rejection{Code: RejectResource}
	}
	digest, err := fileDigestContext(n.ctx, source)
	if err != nil {
		if ctxErr := n.ctx.Err(); ctxErr != nil {
			return "", contextRejection(ctxErr)
		}
		return "", Rejection{Code: RejectUnsafe}
	}
	index := len(n.result.Entries)
	relative := filepath.ToSlash(filepath.Join("raw", fmt.Sprintf("%08d.bin", index)))
	target := filepath.Join(n.result.RawRoot, fmt.Sprintf("%08d.bin", index))
	if err := copyExclusive(n.ctx, source, target); err != nil {
		return "", err
	}
	targetDigest, err := fileDigestContext(n.ctx, target)
	if err != nil || targetDigest != digest {
		return "", Rejection{Code: RejectUnsafe}
	}
	chunks, err := writeChunks(n.ctx, target, n.result.ProbeRoot, index, digest, info.Size())
	if err != nil {
		return "", err
	}
	n.result.Entries = append(n.result.Entries, Entry{Class: class, Depth: depth,
		Size: info.Size(), Digest: digest, Projection: relative, Chunks: chunks})
	for _, chunk := range chunks {
		n.result.ExpectedProbeFiles = append(n.result.ExpectedProbeFiles, chunk.RelativePath)
	}
	return target, nil
}

func writeChunks(ctx context.Context, source, root string, entryIndex int, digest string, size int64) ([]Chunk, error) {
	in, err := os.Open(source)
	if err != nil {
		return nil, Rejection{Code: RejectInvariant}
	}
	defer in.Close()
	chunks := []Chunk{}
	step := DetectorPayloadSize - DetectorOverlap
	for start, index := int64(0), 0; start < size || size == 0 && index == 0; start, index = start+step, index+1 {
		if err := ctx.Err(); err != nil {
			return nil, contextRejection(err)
		}
		end := start + DetectorPayloadSize
		if end > size {
			end = size
		}
		relative := filepath.ToSlash(filepath.Join("chunks", fmt.Sprintf("%08d", entryIndex), fmt.Sprintf("%08d.dat", index)))
		marker := coverageMarker(digest, relative, start, end)
		header := fmt.Sprintf("%s\n%s\nentry=%08d chunk=%08d start=%d end=%d\n", projectionHeader, marker, entryIndex, index, start, end)
		if int64(len(header))+end-start >= MaximumDetectorFile {
			return nil, Rejection{Code: RejectInvariant}
		}
		target := filepath.Join(root, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
			return nil, Rejection{Code: RejectInvariant}
		}
		out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
		if err != nil {
			return nil, Rejection{Code: RejectInvariant}
		}
		_, writeErr := io.WriteString(out, header)
		if writeErr == nil {
			_, writeErr = in.Seek(start, io.SeekStart)
		}
		if writeErr == nil {
			_, writeErr = io.CopyN(out, contextReader{ctx: ctx, r: in}, end-start)
		}
		if closeErr := out.Close(); writeErr == nil {
			writeErr = closeErr
		}
		if writeErr != nil {
			_ = os.Remove(target)
			if ctxErr := ctx.Err(); ctxErr != nil {
				return nil, contextRejection(ctxErr)
			}
			return nil, Rejection{Code: RejectInvariant}
		}
		chunkDigest, digestErr := fileDigestContext(ctx, target)
		chunkInfo, statErr := regularNoLink(target)
		if digestErr != nil || statErr != nil || chunkInfo.Size() >= MaximumDetectorFile {
			return nil, Rejection{Code: RejectInvariant}
		}
		chunks = append(chunks, Chunk{RelativePath: relative, Start: start, End: end,
			Size: chunkInfo.Size(), Digest: chunkDigest, Marker: marker})
		if end == size {
			break
		}
	}
	return chunks, nil
}

func Verify(result Result) error { return VerifyContext(context.Background(), result) }

func VerifyContext(ctx context.Context, result Result) error {
	if ctx == nil {
		return errors.New("invalid artifact projection")
	}
	if len(result.Entries) == 0 || len(result.ExpectedProbeFiles) == 0 ||
		!safeAbsolute(result.RawRoot) || !safeAbsolute(result.ProbeRoot) || result.RawRoot == result.ProbeRoot ||
		!validClass(result.Class) || !result.Profile.internallyValid() || !lowerHex(result.LedgerDigest) || !lowerHex(result.ProjectionDigest) {
		return errors.New("invalid artifact projection")
	}
	ledger, projection := sha256.New(), sha256.New()
	seenRaw, seenProbe := map[string]fileBinding{}, map[string]fileBinding{}
	var expanded int64
	expectedIndex := 0
	for index, entry := range result.Entries {
		if err := ctx.Err(); err != nil {
			return contextRejection(err)
		}
		if entry.Depth < 0 || entry.Depth > result.Profile.MaxDepth || entry.Size < 0 ||
			!validClass(entry.Class) || !lowerHex(entry.Digest) || len(entry.Chunks) == 0 ||
			entry.Projection != filepath.ToSlash(filepath.Join("raw", fmt.Sprintf("%08d.bin", index))) {
			return errors.New("invalid artifact entry")
		}
		raw := filepath.Join(result.RawRoot, fmt.Sprintf("%08d.bin", index))
		info, err := regularNoLink(raw)
		digest, digestErr := fileDigestContext(ctx, raw)
		rawRelative := filepath.ToSlash(filepath.Base(raw))
		if err != nil || digestErr != nil || info.Size() != entry.Size || digest != entry.Digest {
			return errors.New("artifact ledger mismatch")
		}
		if _, exists := seenRaw[rawRelative]; exists {
			return errors.New("artifact ledger mismatch")
		}
		seenRaw[rawRelative] = fileBinding{size: entry.Size, digest: entry.Digest}
		expanded += entry.Size
		previousEnd := int64(0)
		for chunkIndex, chunk := range entry.Chunks {
			if err := ctx.Err(); err != nil {
				return contextRejection(err)
			}
			if chunk.Start < 0 || chunk.End < chunk.Start || chunk.End > entry.Size ||
				chunkIndex == 0 && chunk.Start != 0 || chunkIndex > 0 && chunk.Start != previousEnd-DetectorOverlap ||
				!lowerHex(chunk.Digest) || chunk.Marker != coverageMarker(entry.Digest, chunk.RelativePath, chunk.Start, chunk.End) {
				return errors.New("artifact coverage discontinuity")
			}
			if expectedIndex >= len(result.ExpectedProbeFiles) || result.ExpectedProbeFiles[expectedIndex] != chunk.RelativePath {
				return errors.New("artifact coverage list mismatch")
			}
			expectedIndex++
			path := filepath.Join(result.ProbeRoot, filepath.FromSlash(chunk.RelativePath))
			chunkInfo, statErr := regularNoLink(path)
			chunkDigest, hashErr := fileDigestContext(ctx, path)
			if statErr != nil || hashErr != nil || chunkInfo.Size() != chunk.Size || chunkDigest != chunk.Digest {
				return errors.New("artifact projection mismatch")
			}
			if err := verifyChunkPayload(ctx, raw, path, index, chunkIndex, entry, chunk); err != nil {
				return err
			}
			if _, exists := seenProbe[chunk.RelativePath]; exists {
				return errors.New("artifact coverage duplicate")
			}
			seenProbe[chunk.RelativePath] = fileBinding{size: chunk.Size, digest: chunk.Digest}
			writeRecord(projection, chunk.RelativePath, strconv.FormatInt(chunk.Start, 10), strconv.FormatInt(chunk.End, 10), chunk.Digest, chunk.Marker)
			previousEnd = chunk.End
		}
		if previousEnd != entry.Size {
			return errors.New("artifact coverage incomplete")
		}
		writeRecord(ledger, string(entry.Class), strconv.Itoa(entry.Depth), strconv.FormatInt(entry.Size, 10), entry.Digest, entry.Projection)
	}
	if expectedIndex != len(result.ExpectedProbeFiles) || expanded != result.ExpandedBytes || len(seenProbe) != len(result.ExpectedProbeFiles) ||
		hex.EncodeToString(ledger.Sum(nil)) != result.LedgerDigest || hex.EncodeToString(projection.Sum(nil)) != result.ProjectionDigest {
		return errors.New("artifact proof digest mismatch")
	}
	if err := verifyExactTree(ctx, result.RawRoot, seenRaw); err != nil {
		return err
	}
	if err := verifyExactTree(ctx, result.ProbeRoot, seenProbe); err != nil {
		return err
	}
	return nil
}

func verifyChunkPayload(ctx context.Context, rawPath, chunkPath string, entryIndex, chunkIndex int, entry Entry, chunk Chunk) error {
	raw, err := os.Open(rawPath)
	if err != nil {
		return errors.New("artifact projection mismatch")
	}
	defer raw.Close()
	projected, err := os.Open(chunkPath)
	if err != nil {
		return errors.New("artifact projection mismatch")
	}
	defer projected.Close()
	header := fmt.Sprintf("%s\n%s\nentry=%08d chunk=%08d start=%d end=%d\n", projectionHeader, chunk.Marker, entryIndex, chunkIndex, chunk.Start, chunk.End)
	if chunk.Size != int64(len(header))+chunk.End-chunk.Start {
		return errors.New("artifact projection mismatch")
	}
	actualHeader := make([]byte, len(header))
	if _, err := io.ReadFull(projected, actualHeader); err != nil || string(actualHeader) != header {
		return errors.New("artifact projection mismatch")
	}
	if _, err := raw.Seek(chunk.Start, io.SeekStart); err != nil {
		return errors.New("artifact projection mismatch")
	}
	want := io.LimitReader(raw, chunk.End-chunk.Start)
	left, right := make([]byte, 64<<10), make([]byte, 64<<10)
	remaining := chunk.End - chunk.Start
	for remaining > 0 {
		if err := ctx.Err(); err != nil {
			return contextRejection(err)
		}
		step := int64(len(left))
		if step > remaining {
			step = remaining
		}
		if _, err := io.ReadFull(want, left[:step]); err != nil {
			return errors.New("artifact projection mismatch")
		}
		if _, err := io.ReadFull(projected, right[:step]); err != nil || !bytes.Equal(left[:step], right[:step]) {
			return errors.New("artifact projection mismatch")
		}
		remaining -= step
	}
	one := make([]byte, 1)
	if n, err := projected.Read(one); n != 0 || err != io.EOF {
		return errors.New("artifact projection mismatch")
	}
	return nil
}

type fileBinding struct {
	size   int64
	digest string
}

func verifyExactTree(ctx context.Context, root string, expected map[string]fileBinding) error {
	seen := map[string]bool{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if err := ctx.Err(); err != nil {
			return contextRejection(err)
		}
		if walkErr != nil {
			return errors.New("artifact proof tree unavailable")
		}
		if path == root {
			return nil
		}
		info, err := os.Lstat(path)
		if err != nil || info.Mode()&os.ModeSymlink != 0 || isReparse(info) {
			return errors.New("artifact proof tree unsafe")
		}
		if info.IsDir() {
			return nil
		}
		if !info.Mode().IsRegular() || hasMultipleLinks(path, info) {
			return errors.New("artifact proof tree unsafe")
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return errors.New("artifact proof tree mismatch")
		}
		relative = filepath.ToSlash(relative)
		binding, ok := expected[relative]
		if !ok || seen[relative] || info.Size() != binding.size {
			return errors.New("artifact proof tree mismatch")
		}
		digest, err := fileDigestContext(ctx, path)
		if err != nil || digest != binding.digest {
			return errors.New("artifact proof tree mismatch")
		}
		seen[relative] = true
		return nil
	})
	if err != nil || len(seen) != len(expected) {
		return errors.New("artifact proof tree mismatch")
	}
	return nil
}

func validClass(class Class) bool {
	switch class {
	case ClassDirectory, ClassFile, ClassZIP, ClassTAR, ClassGZIP, ClassTGZ, ClassOCILayout, ClassDockerSave:
		return true
	default:
		return false
	}
}

func seal(result *Result) {
	ledger, projection := sha256.New(), sha256.New()
	for _, entry := range result.Entries {
		writeRecord(ledger, string(entry.Class), strconv.Itoa(entry.Depth), strconv.FormatInt(entry.Size, 10), entry.Digest, entry.Projection)
		for _, chunk := range entry.Chunks {
			writeRecord(projection, chunk.RelativePath, strconv.FormatInt(chunk.Start, 10), strconv.FormatInt(chunk.End, 10), chunk.Digest, chunk.Marker)
		}
	}
	result.LedgerDigest = hex.EncodeToString(ledger.Sum(nil))
	result.ProjectionDigest = hex.EncodeToString(projection.Sum(nil))
}

func coverageMarker(digest, path string, start, end int64) string {
	h := sha256.New()
	writeRecord(h, "pscan.artifact.coverage.v1", digest, path, strconv.FormatInt(start, 10), strconv.FormatInt(end, 10))
	return "PSCAN_COVERAGE_MARKER_" + hex.EncodeToString(h.Sum(nil))
}

func copyExclusive(ctx context.Context, source, target string) error {
	in, err := os.Open(source)
	if err != nil {
		return Rejection{Code: RejectUnsafe}
	}
	defer in.Close()
	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return Rejection{Code: RejectInvariant}
	}
	_, copyErr := io.Copy(out, contextReader{ctx: ctx, r: in})
	if closeErr := out.Close(); copyErr == nil {
		copyErr = closeErr
	}
	if copyErr != nil {
		_ = os.Remove(target)
		if ctxErr := ctx.Err(); ctxErr != nil {
			return contextRejection(ctxErr)
		}
		return Rejection{Code: RejectInvariant}
	}
	return nil
}

func fileDigestContext(ctx context.Context, path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, contextReader{ctx: ctx, r: f}); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func fileDigest(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func writeRecord(w io.Writer, values ...string) {
	var length [8]byte
	for _, value := range values {
		binary.BigEndian.PutUint64(length[:], uint64(len(value)))
		_, _ = w.Write(length[:])
		_, _ = io.WriteString(w, value)
	}
}

func lowerHex(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, r := range value {
		if !strings.ContainsRune("0123456789abcdef", r) {
			return false
		}
	}
	return true
}
