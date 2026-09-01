package artifact

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type normalizer struct {
	ctx             context.Context
	limits          Limits
	result          Result
	staging         string
	entries         int
	count           int
	total           int64
	semanticEntries int
	semanticBytes   int64
	semanticLayers  map[string]string
}

type extracted struct {
	name     string
	path     string
	size     int64
	admitted bool
}

func Normalize(parent context.Context, input, output string, limits Limits) (result Result, err error) {
	if parent == nil || !safeAbsolute(input) || !safeAbsolute(output) || !limits.internallyValid() || input == output || confined(input, output) || confined(output, input) {
		return Result{}, Rejection{Code: RejectInvariant}
	}
	info, statErr := os.Lstat(input)
	if statErr != nil || info.Mode()&os.ModeSymlink != 0 || isReparse(info) || (!info.IsDir() && !info.Mode().IsRegular()) {
		return Result{}, Rejection{Code: RejectUnsafe}
	}
	ctx, cancel := context.WithTimeout(parent, limits.Timeout)
	defer cancel()
	if err := createPrivateRoot(output); err != nil {
		return Result{}, err
	}
	complete := false
	defer func() {
		if recovered := recover(); recovered != nil {
			complete = false
			err = Rejection{Code: RejectInvariant}
		}
		if !complete {
			_ = os.RemoveAll(output)
			result = Result{}
		}
	}()
	for _, child := range []string{"raw", "probe", "staging"} {
		if makeErr := os.Mkdir(filepath.Join(output, child), 0o700); makeErr != nil {
			return Result{}, Rejection{Code: RejectInvariant}
		}
	}
	n := &normalizer{ctx: ctx, limits: limits, staging: filepath.Join(output, "staging"), semanticLayers: map[string]string{}}
	n.result = Result{RawRoot: filepath.Join(output, "raw"), ProbeRoot: filepath.Join(output, "probe"), Profile: limits}
	if info.IsDir() {
		snapshotRoot, paths, err := n.preflightDirectory(input)
		if err != nil {
			return Result{}, err
		}
		if isOCILayout(snapshotRoot) {
			n.result.Class = ClassOCILayout
			if err := n.validateOCI(snapshotRoot); err != nil {
				return Result{}, err
			}
		} else {
			n.result.Class = ClassDirectory
		}
		if err := n.directoryPaths(paths, 0); err != nil {
			return Result{}, err
		}
	} else {
		class, err := classifyNamedContext(ctx, input, filepath.Base(input))
		if err != nil {
			return Result{}, err
		}
		n.result.Class = class
		if err := n.file(input, filepath.Base(input), 0, class, false); err != nil {
			return Result{}, err
		}
	}
	if len(n.result.Entries) == 0 || len(n.result.Entries) != n.entries {
		return Result{}, Rejection{Code: RejectInvariant}
	}
	n.result.ExpandedBytes = n.total
	seal(&n.result)
	if err := os.RemoveAll(n.staging); err != nil {
		return Result{}, Rejection{Code: RejectInvariant}
	}
	if err := VerifyContext(n.ctx, n.result); err != nil {
		var rejection Rejection
		if errors.As(err, &rejection) {
			return Result{}, rejection
		}
		return Result{}, Rejection{Code: RejectInvariant}
	}
	complete = true
	return n.result, nil
}

func (n *normalizer) check() error {
	select {
	case <-n.ctx.Done():
		if errors.Is(n.ctx.Err(), context.DeadlineExceeded) {
			return Rejection{Code: RejectTimeout}
		}
		return Rejection{Code: RejectCancelled}
	default:
		return nil
	}
}

func contextRejection(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return Rejection{Code: RejectTimeout}
	}
	return Rejection{Code: RejectCancelled}
}

type contextReader struct {
	ctx context.Context
	r   io.Reader
}

func (r contextReader) Read(buffer []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.r.Read(buffer)
}

func (n *normalizer) admit(size int64) error {
	if err := n.check(); err != nil {
		return err
	}
	if size < 0 || size > n.limits.MaxFileBytes || n.count >= n.limits.MaxEntries || size > n.limits.MaxExpandedBytes-n.total {
		return Rejection{Code: RejectResource}
	}
	n.count++
	n.entries++
	n.total += size
	return nil
}

func (n *normalizer) admitDirectory() error {
	if err := n.check(); err != nil {
		return err
	}
	if n.count >= n.limits.MaxEntries {
		return Rejection{Code: RejectResource}
	}
	n.count++
	return nil
}

type sourceRecord struct {
	directory bool
	info      os.FileInfo
	size      int64
	digest    string
}

func (n *normalizer) preflightDirectory(root string) (string, []string, error) {
	paths := []string{}
	seenTree := map[string]string{}
	records := map[string]sourceRecord{}
	snapshotRoot := filepath.Join(n.staging, "source")
	if err := os.Mkdir(snapshotRoot, 0o700); err != nil {
		return "", nil, Rejection{Code: RejectInvariant}
	}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return Rejection{Code: RejectUnsafe}
		}
		if err := n.check(); err != nil {
			return err
		}
		if path == root {
			return nil
		}
		relative, relErr := filepath.Rel(root, path)
		if relErr != nil || !safeArchivePath(filepath.ToSlash(relative)) {
			return Rejection{Code: RejectUnsafe}
		}
		folded := strings.ToLower(filepath.ToSlash(relative))
		if previous, exists := seenTree[folded]; exists && previous != filepath.ToSlash(relative) {
			return Rejection{Code: RejectUnsafe}
		}
		seenTree[folded] = filepath.ToSlash(relative)
		target := filepath.Join(snapshotRoot, filepath.FromSlash(filepath.ToSlash(relative)))
		if !confined(snapshotRoot, target) {
			return Rejection{Code: RejectUnsafe}
		}
		if entry.IsDir() {
			if err := n.admitDirectory(); err != nil {
				return err
			}
			if err := os.Mkdir(target, 0o700); err != nil {
				return Rejection{Code: RejectInvariant}
			}
			info, err := os.Lstat(path)
			if err != nil {
				return Rejection{Code: RejectUnsafe}
			}
			records[filepath.ToSlash(relative)] = sourceRecord{directory: true, info: info}
			return nil
		}
		info, err := os.Lstat(path)
		if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || isReparse(info) || hasMultipleLinks(path, info) {
			return Rejection{Code: RejectUnsafe}
		}
		if err := n.admit(info.Size()); err != nil {
			return err
		}
		digest, err := fileDigestContext(n.ctx, path)
		if err != nil {
			return Rejection{Code: RejectUnsafe}
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
			return Rejection{Code: RejectInvariant}
		}
		if err := copyExclusive(n.ctx, path, target); err != nil {
			return err
		}
		targetDigest, err := fileDigestContext(n.ctx, target)
		if err != nil || targetDigest != digest {
			return Rejection{Code: RejectUnsafe}
		}
		records[filepath.ToSlash(relative)] = sourceRecord{info: info, size: info.Size(), digest: digest}
		paths = append(paths, target)
		return nil
	})
	if err != nil {
		return "", nil, err
	}
	if err := n.verifySourceTree(root, records); err != nil {
		return "", nil, err
	}
	sort.Strings(paths)
	return snapshotRoot, paths, nil
}

func (n *normalizer) verifySourceTree(root string, expected map[string]sourceRecord) error {
	seen := map[string]bool{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return Rejection{Code: RejectUnsafe}
		}
		if err := n.check(); err != nil {
			return err
		}
		if path == root {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return Rejection{Code: RejectUnsafe}
		}
		relative = filepath.ToSlash(relative)
		record, ok := expected[relative]
		info, statErr := os.Lstat(path)
		if !ok || seen[relative] || statErr != nil || record.directory != info.IsDir() || !os.SameFile(record.info, info) {
			return Rejection{Code: RejectUnsafe}
		}
		if !record.directory {
			digest, err := fileDigestContext(n.ctx, path)
			if err != nil || info.Size() != record.size || digest != record.digest {
				return Rejection{Code: RejectUnsafe}
			}
		}
		seen[relative] = true
		return nil
	})
	if err != nil || len(seen) != len(expected) {
		return Rejection{Code: RejectUnsafe}
	}
	return nil
}

func (n *normalizer) directoryPaths(paths []string, depth int) error {
	if depth > n.limits.MaxDepth {
		return Rejection{Code: RejectResource}
	}
	seen := map[string]bool{}
	for _, path := range paths {
		folded := strings.ToLower(filepath.ToSlash(path))
		if seen[folded] {
			return Rejection{Code: RejectUnsafe}
		}
		seen[folded] = true
		class, classErr := classifyNamedContext(n.ctx, path, filepath.Base(path))
		if classErr != nil {
			return classErr
		}
		if err := n.file(path, filepath.Base(path), depth, class, true); err != nil {
			return err
		}
	}
	return nil
}

func (n *normalizer) reserveSemantic(entries int, bytes int64) error {
	if entries < 0 || bytes < 0 || entries > n.limits.MaxEntries-n.count-n.semanticEntries || bytes > n.limits.MaxExpandedBytes-n.total-n.semanticBytes {
		return Rejection{Code: RejectResource}
	}
	n.semanticEntries += entries
	n.semanticBytes += bytes
	return nil
}

func (n *normalizer) file(path, logicalName string, depth int, class Class, admitted bool) error {
	if depth > n.limits.MaxDepth {
		return Rejection{Code: RejectResource}
	}
	info, err := regularNoLink(path)
	if err != nil {
		return err
	}
	if !admitted {
		if err := n.admit(info.Size()); err != nil {
			return err
		}
	}
	snapshot, err := n.project(path, class, depth)
	if err != nil {
		return err
	}
	snapshotClass, err := classifyNamedContext(n.ctx, snapshot, logicalName)
	if err != nil || snapshotClass != class {
		return Rejection{Code: RejectUnsafe}
	}
	switch class {
	case ClassFile:
		return nil
	case ClassZIP:
		members, err := n.extractZIP(snapshot, depth)
		if err != nil {
			return err
		}
		return n.members(members, depth+1)
	case ClassTAR:
		members, err := n.extractTAR(snapshot, depth)
		if err != nil {
			return err
		}
		if dockerSave(members) {
			if err := n.validateDockerSave(members); err != nil {
				return err
			}
			// The root class is upgraded only when the submitted object itself is a save.
			if depth == 0 && n.result.Class == ClassTAR {
				n.result.Class = ClassDockerSave
			}
		}
		return n.members(members, depth+1)
	case ClassGZIP:
		member, tgz, err := n.extractGZIP(snapshot, depth)
		if err != nil {
			return err
		}
		if tgz && depth == 0 && n.result.Class == ClassGZIP {
			n.result.Class = ClassTGZ
		}
		return n.members([]extracted{member}, depth+1)
	default:
		return Rejection{Code: RejectUnsupported}
	}
}

func (n *normalizer) members(members []extracted, depth int) error {
	if depth > n.limits.MaxDepth {
		return Rejection{Code: RejectResource}
	}
	if len(members) == 0 {
		return Rejection{Code: RejectMalformed}
	}
	for _, member := range members {
		class, err := classifyNamedContext(n.ctx, member.path, member.name)
		if err != nil {
			return err
		}
		if err := n.file(member.path, member.name, depth, class, member.admitted); err != nil {
			return err
		}
	}
	return nil
}

func classify(path string) (Class, error) {
	return classifyNamedContext(context.Background(), path, filepath.Base(path))
}

func classifyNamed(path, logicalName string) (Class, error) {
	return classifyNamedContext(context.Background(), path, logicalName)
}

func classifyNamedContext(ctx context.Context, path, logicalName string) (Class, error) {
	if err := ctx.Err(); err != nil {
		return "", contextRejection(err)
	}
	f, err := os.Open(path)
	if err != nil {
		return "", Rejection{Code: RejectUnsafe}
	}
	defer f.Close()
	header := make([]byte, 560)
	n, readErr := io.ReadFull(contextReader{ctx: ctx, r: f}, header)
	if readErr != nil && readErr != io.ErrUnexpectedEOF && readErr != io.EOF {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return "", contextRejection(ctxErr)
		}
		return "", Rejection{Code: RejectMalformed}
	}
	header = header[:n]
	zipMagic := len(header) >= 4 && (bytes.Equal(header[:4], []byte{'P', 'K', 3, 4}) || bytes.Equal(header[:4], []byte{'P', 'K', 5, 6}) || bytes.Equal(header[:4], []byte{'P', 'K', 7, 8}))
	gzipMagic := len(header) >= 2 && header[0] == 0x1f && header[1] == 0x8b
	tarMagic := looksLikeTARHeader(header)
	known := 0
	for _, present := range []bool{zipMagic, gzipMagic, tarMagic} {
		if present {
			known++
		}
	}
	if known > 1 {
		return "", Rejection{Code: RejectUnsupported}
	}
	broad, broadErr := broadClassify(ctx, path, logicalName)
	if IsCode(broadErr, RejectTimeout) || IsCode(broadErr, RejectCancelled) {
		return "", broadErr
	}
	if zipMagic {
		if broadErr == nil || broad != broadArchive {
			return "", Rejection{Code: RejectUnsupported}
		}
		return ClassZIP, nil
	}
	if gzipMagic {
		if broadErr == nil || broad != broadCompression {
			return "", Rejection{Code: RejectUnsupported}
		}
		return ClassGZIP, nil
	}
	if tarMagic {
		if broadErr == nil || broad != broadArchive {
			return "", Rejection{Code: RejectUnsupported}
		}
		return ClassTAR, nil
	}
	unsupported := [][]byte{{'7', 'z', 0xbc, 0xaf, 0x27, 0x1c}, {'R', 'a', 'r', '!'}, {'B', 'Z', 'h'}, {0xfd, '7', 'z', 'X', 'Z', 0}, {0x28, 0xb5, 0x2f, 0xfd}}
	for _, magic := range unsupported {
		if len(header) >= len(magic) && bytes.Equal(header[:len(magic)], magic) {
			return "", Rejection{Code: RejectUnsupported}
		}
	}
	if broadErr != nil || broad != broadText && broad != broadBinary {
		return "", Rejection{Code: RejectUnsupported}
	}
	return ClassFile, nil
}

func (n *normalizer) extractionPath() (string, error) {
	path := filepath.Join(n.staging, fmt.Sprintf("%08d", len(n.result.Entries)))
	if !confined(n.staging, path) || os.Mkdir(path, 0o700) != nil {
		return "", Rejection{Code: RejectInvariant}
	}
	return path, nil
}

func (n *normalizer) extractZIP(path string, depth int) ([]extracted, error) {
	if err := preflightZIP(n.ctx, path, n.limits.MaxEntries-n.count, n.limits.MaxExpandedBytes-n.total, n.limits); err != nil {
		return nil, err
	}
	reader, err := zip.OpenReader(path)
	if err != nil || len(reader.File) == 0 {
		return nil, Rejection{Code: RejectMalformed}
	}
	defer reader.Close()
	root, err := n.extractionPath()
	if err != nil {
		return nil, err
	}
	members := []extracted{}
	seen := map[string]bool{}
	for index, file := range reader.File {
		if err := n.check(); err != nil {
			return nil, err
		}
		name := strings.TrimSuffix(file.Name, "/")
		if !safeArchivePath(name) || seen[strings.ToLower(name)] || file.Flags&1 != 0 ||
			file.Method != zip.Store && file.Method != zip.Deflate {
			return nil, Rejection{Code: RejectUnsafe}
		}
		seen[strings.ToLower(name)] = true
		mode := file.Mode()
		if file.FileInfo().IsDir() {
			if err := n.admitDirectory(); err != nil {
				return nil, err
			}
			continue
		}
		if !mode.IsRegular() || mode&os.ModeSymlink != 0 || file.UncompressedSize64 > uint64(n.limits.MaxFileBytes) {
			return nil, Rejection{Code: RejectUnsafe}
		}
		if exceedsRatio(int64(file.UncompressedSize64), int64(file.CompressedSize64), n.limits.MaxCompressionRatio) {
			return nil, Rejection{Code: RejectResource}
		}
		if err := n.admit(int64(file.UncompressedSize64)); err != nil {
			return nil, err
		}
		target := filepath.Join(root, fmt.Sprintf("%08d.bin", index))
		in, openErr := file.Open()
		if openErr != nil {
			return nil, Rejection{Code: RejectMalformed}
		}
		writeErr := n.writeBounded(target, in, int64(file.UncompressedSize64))
		closeErr := in.Close()
		if writeErr != nil {
			return nil, writeErr
		}
		if closeErr != nil {
			return nil, Rejection{Code: RejectMalformed}
		}
		members = append(members, extracted{name: name, path: target, size: int64(file.UncompressedSize64), admitted: true})
	}
	return members, nil
}

func (n *normalizer) extractTAR(path string, depth int) ([]extracted, error) {
	if _, err := preflightTAR(n.ctx, path, n.limits.MaxEntries-n.count, n.limits.MaxExpandedBytes-n.total, n.limits.MaxFileBytes); err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, Rejection{Code: RejectUnsafe}
	}
	defer f.Close()
	return n.extractTARReader(tar.NewReader(contextReader{ctx: n.ctx, r: f}))
}

func (n *normalizer) extractTARReader(reader *tar.Reader) ([]extracted, error) {
	root, err := n.extractionPath()
	if err != nil {
		return nil, err
	}
	members := []extracted{}
	seen := map[string]bool{}
	for index := 0; ; index++ {
		if err := n.check(); err != nil {
			return nil, err
		}
		header, nextErr := reader.Next()
		if nextErr == io.EOF {
			break
		}
		if nextErr != nil {
			return nil, Rejection{Code: RejectMalformed}
		}
		name := strings.TrimSuffix(header.Name, "/")
		if !safeArchivePath(name) || seen[strings.ToLower(name)] {
			return nil, Rejection{Code: RejectUnsafe}
		}
		seen[strings.ToLower(name)] = true
		switch header.Typeflag {
		case tar.TypeDir:
			if err := n.admitDirectory(); err != nil {
				return nil, err
			}
			continue
		case tar.TypeReg, tar.TypeRegA:
		default:
			return nil, Rejection{Code: RejectUnsafe}
		}
		if header.Size < 0 || header.Size > n.limits.MaxFileBytes {
			return nil, Rejection{Code: RejectResource}
		}
		if err := n.admit(header.Size); err != nil {
			return nil, err
		}
		target := filepath.Join(root, fmt.Sprintf("%08d.bin", index))
		if err := n.writeBounded(target, reader, header.Size); err != nil {
			return nil, err
		}
		members = append(members, extracted{name: name, path: target, size: header.Size, admitted: true})
	}
	return members, nil
}

func (n *normalizer) extractGZIP(path string, depth int) (extracted, bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return extracted{}, false, Rejection{Code: RejectUnsafe}
	}
	defer f.Close()
	buffered := bufio.NewReader(f)
	reader, err := gzip.NewReader(buffered)
	if err != nil {
		return extracted{}, false, Rejection{Code: RejectMalformed}
	}
	reader.Multistream(false)
	root, err := n.extractionPath()
	if err != nil {
		return extracted{}, false, err
	}
	if n.count >= n.limits.MaxEntries {
		return extracted{}, false, Rejection{Code: RejectResource}
	}
	n.count++
	n.entries++
	target := filepath.Join(root, "00000000.bin")
	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return extracted{}, false, Rejection{Code: RejectInvariant}
	}
	writer := &expansionWriter{normalizer: n, writer: out}
	written, copyErr := io.Copy(writer, io.LimitReader(contextReader{ctx: n.ctx, r: reader}, n.limits.MaxFileBytes+1))
	closeOutErr := out.Close()
	closeReaderErr := reader.Close()
	if copyErr != nil {
		var rejection Rejection
		if errors.As(copyErr, &rejection) {
			return extracted{}, false, rejection
		}
		return extracted{}, false, Rejection{Code: RejectMalformed}
	}
	if closeOutErr != nil || closeReaderErr != nil || written > n.limits.MaxFileBytes {
		return extracted{}, false, Rejection{Code: RejectResource}
	}
	if _, peekErr := buffered.Peek(1); peekErr != io.EOF {
		return extracted{}, false, Rejection{Code: RejectMalformed}
	}
	compressed, _ := regularNoLink(path)
	if exceedsRatio(written, compressed.Size(), n.limits.MaxCompressionRatio) {
		return extracted{}, false, Rejection{Code: RejectResource}
	}
	class, classErr := classifyNamedContext(n.ctx, target, filepath.Base(target))
	if classErr != nil {
		return extracted{}, false, classErr
	}
	return extracted{name: "member", path: target, size: written, admitted: true}, class == ClassTAR, nil
}

type expansionWriter struct {
	normalizer *normalizer
	writer     io.Writer
	written    int64
}

func (w *expansionWriter) Write(value []byte) (int, error) {
	if err := w.normalizer.check(); err != nil {
		return 0, err
	}
	remainingFile := w.normalizer.limits.MaxFileBytes - w.written
	remainingTotal := w.normalizer.limits.MaxExpandedBytes - w.normalizer.total
	if int64(len(value)) > remainingFile || int64(len(value)) > remainingTotal {
		return 0, Rejection{Code: RejectResource}
	}
	n, err := w.writer.Write(value)
	w.written += int64(n)
	w.normalizer.total += int64(n)
	return n, err
}

func exceedsRatio(expanded, compressed, maximum int64) bool {
	if expanded <= 0 {
		return false
	}
	if compressed <= 0 || maximum <= 0 {
		return true
	}
	if compressed > (1<<63-1)/maximum {
		return false
	}
	return expanded > compressed*maximum
}

func (n *normalizer) writeBounded(target string, reader io.Reader, expected int64) error {
	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	written, copyErr := io.Copy(out, io.LimitReader(contextReader{ctx: n.ctx, r: reader}, expected+1))
	if closeErr := out.Close(); copyErr == nil {
		copyErr = closeErr
	}
	if copyErr != nil || written != expected {
		_ = os.Remove(target)
		var rejection Rejection
		if errors.As(copyErr, &rejection) {
			return rejection
		}
		if ctxErr := n.ctx.Err(); ctxErr != nil {
			return contextRejection(ctxErr)
		}
		return errors.New("bounded extraction failed")
	}
	return nil
}

type dockerManifest struct {
	Config   string   `json:"Config"`
	RepoTags []string `json:"RepoTags"`
	Layers   []string `json:"Layers"`
}

var dockerConfigName = regexp.MustCompile(`^([0-9a-f]{64})\.json$`)
var dockerBlobName = regexp.MustCompile(`^blobs/sha256/([0-9a-f]{64})$`)

func dockerSave(members []extracted) bool {
	for _, member := range members {
		if member.name == "manifest.json" {
			return true
		}
	}
	return false
}

func (n *normalizer) validateDockerSave(members []extracted) error {
	byName := map[string]extracted{}
	for _, member := range members {
		byName[member.name] = member
	}
	manifestMember, ok := byName["manifest.json"]
	if !ok || manifestMember.size > 16<<20 {
		return Rejection{Code: RejectMalformed}
	}
	raw, err := os.ReadFile(manifestMember.path)
	if err != nil {
		return Rejection{Code: RejectMalformed}
	}
	if err := rejectDuplicateJSON(raw); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	var manifest []dockerManifest
	if decoder.Decode(&manifest) != nil || len(manifest) == 0 || decoder.Decode(&struct{}{}) != io.EOF {
		return Rejection{Code: RejectMalformed}
	}
	referenced := map[string]bool{"manifest.json": true}
	modern := false
	for _, image := range manifest {
		if !safeArchivePath(image.Config) || len(image.Layers) == 0 || referenced[image.Config] {
			return Rejection{Code: RejectMalformed}
		}
		referenced[image.Config] = true
		if dockerBlobName.MatchString(image.Config) {
			modern = true
		}
		for _, layer := range image.Layers {
			if !safeArchivePath(layer) || referenced[layer] {
				return Rejection{Code: RejectMalformed}
			}
			referenced[layer] = true
		}
	}
	for name := range referenced {
		if _, exists := byName[name]; !exists {
			return Rejection{Code: RejectMalformed}
		}
	}
	if modern {
		layout, layoutOK := byName["oci-layout"]
		index, indexOK := byName["index.json"]
		if !layoutOK || !indexOK || layout.size > 4096 || index.size > 16<<20 {
			return Rejection{Code: RejectMalformed}
		}
		for name, member := range byName {
			if name == "manifest.json" || name == "oci-layout" || name == "index.json" {
				continue
			}
			match := dockerBlobName.FindStringSubmatch(name)
			digest, err := fileDigestContext(n.ctx, member.path)
			if err != nil || len(match) != 2 || digest != match[1] {
				return Rejection{Code: RejectMalformed}
			}
		}
		if err := n.validateDockerOCIRoot(layout.path, index.path, byName, manifest); err != nil {
			return err
		}
	} else if len(referenced) != len(byName) {
		return Rejection{Code: RejectMalformed}
	}
	for _, image := range manifest {
		config := byName[image.Config]
		if config.size > 16<<20 {
			return Rejection{Code: RejectResource}
		}
		configRaw, err := os.ReadFile(config.path)
		configNameMatch := dockerConfigName.FindStringSubmatch(image.Config)
		if len(configNameMatch) != 2 {
			configNameMatch = dockerBlobName.FindStringSubmatch(image.Config)
		}
		configSum := sha256.Sum256(configRaw)
		if err != nil || len(configNameMatch) != 2 || hex.EncodeToString(configSum[:]) != configNameMatch[1] || !jsonObject(configRaw) {
			return Rejection{Code: RejectMalformed}
		}
		configClass, err := classifyNamedContext(n.ctx, config.path, image.Config)
		if err != nil || configClass != ClassFile {
			return Rejection{Code: RejectMalformed}
		}
		diffIDs := make([]string, 0, len(image.Layers))
		for _, layerName := range image.Layers {
			layer := byName[layerName]
			layerClass, err := classifyNamedContext(n.ctx, layer.path, layerName)
			if err != nil || layerClass != ClassTAR && layerClass != ClassGZIP {
				return Rejection{Code: RejectMalformed}
			}
			mediaType := "application/vnd.oci.image.layer.v1.tar"
			if layerClass == ClassGZIP {
				mediaType = "application/vnd.docker.image.rootfs.diff.tar.gzip"
			}
			diffID, err := n.layerDiffID(layer.path, mediaType)
			if err != nil {
				return err
			}
			diffIDs = append(diffIDs, diffID)
		}
		if err := validateImageConfiguration(configRaw, diffIDs); err != nil {
			return err
		}
	}
	return nil
}

func (n *normalizer) validateDockerOCIRoot(layoutPath, indexPath string, byName map[string]extracted, selected []dockerManifest) error {
	layoutRaw, err := os.ReadFile(layoutPath)
	if err != nil || rejectDuplicateJSON(layoutRaw) != nil {
		return Rejection{Code: RejectMalformed}
	}
	var layout map[string]json.RawMessage
	if json.Unmarshal(layoutRaw, &layout) != nil || len(layout) != 1 {
		return Rejection{Code: RejectMalformed}
	}
	var version string
	if json.Unmarshal(layout["imageLayoutVersion"], &version) != nil || version != "1.0.0" {
		return Rejection{Code: RejectMalformed}
	}
	indexRaw, err := os.ReadFile(indexPath)
	if err != nil || rejectDuplicateJSON(indexRaw) != nil {
		return Rejection{Code: RejectMalformed}
	}
	var index ociIndex
	if json.Unmarshal(indexRaw, &index) != nil || index.SchemaVersion != 2 || len(index.Manifests) == 0 {
		return Rejection{Code: RejectMalformed}
	}
	matched := make([]bool, len(selected))
	queue := append([]descriptor(nil), index.Manifests...)
	seen := map[string]descriptor{}
	for len(queue) > 0 {
		if err := n.check(); err != nil {
			return err
		}
		d := queue[0]
		queue = queue[1:]
		match := ociDigest.FindStringSubmatch(d.Digest)
		if len(match) != 2 || d.Size < 0 || d.Size > n.limits.MaxFileBytes || d.MediaType == "" || len(d.URLs) != 0 || d.Data != "" {
			return Rejection{Code: RejectMalformed}
		}
		name := "blobs/sha256/" + match[1]
		if previous, ok := seen[name]; ok {
			if previous.Size != d.Size || previous.MediaType != d.MediaType || len(d.URLs) != 0 || d.Data != "" {
				return Rejection{Code: RejectMalformed}
			}
			continue
		}
		seen[name] = d
		member, present := byName[name]
		if !present {
			// Docker 29 preserves multi-platform index descriptors whose other
			// platform blobs are intentionally absent from a selected save.
			continue
		}
		if member.size != d.Size {
			return Rejection{Code: RejectMalformed}
		}
		switch d.MediaType {
		case "application/vnd.oci.image.index.v1+json":
			if d.Size > 16<<20 {
				return Rejection{Code: RejectResource}
			}
			raw, err := os.ReadFile(member.path)
			if err != nil || rejectDuplicateJSON(raw) != nil {
				return Rejection{Code: RejectMalformed}
			}
			var nested ociIndex
			if json.Unmarshal(raw, &nested) != nil || nested.SchemaVersion != 2 || len(nested.Manifests) == 0 {
				return Rejection{Code: RejectMalformed}
			}
			queue = append(queue, nested.Manifests...)
		case "application/vnd.oci.image.manifest.v1+json", "application/vnd.docker.distribution.manifest.v2+json":
			if d.Size > 16<<20 {
				return Rejection{Code: RejectResource}
			}
			raw, err := os.ReadFile(member.path)
			if err != nil || rejectDuplicateJSON(raw) != nil {
				return Rejection{Code: RejectMalformed}
			}
			var image ociManifest
			if json.Unmarshal(raw, &image) != nil || image.SchemaVersion != 2 || len(image.Layers) == 0 {
				return Rejection{Code: RejectMalformed}
			}
			for index, candidate := range selected {
				if dockerManifestMatchesOCI(candidate, image, byName) {
					matched[index] = true
				}
			}
		default:
			return Rejection{Code: RejectUnsupported}
		}
	}
	for _, ok := range matched {
		if !ok {
			return Rejection{Code: RejectMalformed}
		}
	}
	return nil
}

func dockerManifestMatchesOCI(selected dockerManifest, image ociManifest, byName map[string]extracted) bool {
	config := ociDigest.FindStringSubmatch(image.Config.Digest)
	if len(config) != 2 || selected.Config != "blobs/sha256/"+config[1] || len(selected.Layers) != len(image.Layers) {
		return false
	}
	configMember, ok := byName[selected.Config]
	if !ok || configMember.size != image.Config.Size {
		return false
	}
	if image.Config.MediaType != "application/vnd.oci.image.config.v1+json" && image.Config.MediaType != "application/vnd.docker.container.image.v1+json" {
		return false
	}
	for index, layer := range image.Layers {
		match := ociDigest.FindStringSubmatch(layer.Digest)
		if len(match) != 2 || selected.Layers[index] != "blobs/sha256/"+match[1] ||
			layer.MediaType != "application/vnd.oci.image.layer.v1.tar+gzip" && layer.MediaType != "application/vnd.docker.image.rootfs.diff.tar.gzip" {
			return false
		}
		member, ok := byName[selected.Layers[index]]
		if !ok || member.size != layer.Size {
			return false
		}
	}
	return true
}

func jsonObject(raw []byte) bool {
	if rejectDuplicateJSON(raw) != nil {
		return false
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	var value map[string]json.RawMessage
	if decoder.Decode(&value) != nil || value == nil {
		return false
	}
	return decoder.Decode(&struct{}{}) == io.EOF
}
