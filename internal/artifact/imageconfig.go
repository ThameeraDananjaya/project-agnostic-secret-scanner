package artifact

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var platformToken = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,63}$`)

type imageConfiguration struct {
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
	RootFS       struct {
		Type    string   `json:"type"`
		DiffIDs []string `json:"diff_ids"`
	} `json:"rootfs"`
}

func validateImageConfiguration(raw []byte, layerDiffIDs []string) error {
	if rejectDuplicateJSON(raw) != nil {
		return Rejection{Code: RejectMalformed}
	}
	var config imageConfiguration
	if json.Unmarshal(raw, &config) != nil || !platformToken.MatchString(config.Architecture) ||
		!platformToken.MatchString(config.OS) || config.RootFS.Type != "layers" ||
		len(config.RootFS.DiffIDs) != len(layerDiffIDs) || len(layerDiffIDs) == 0 {
		return Rejection{Code: RejectMalformed}
	}
	for index, digest := range config.RootFS.DiffIDs {
		if !ociDigest.MatchString(digest) || digest != layerDiffIDs[index] {
			return Rejection{Code: RejectMalformed}
		}
	}
	return nil
}

func (n *normalizer) layerDiffID(path, mediaType string) (string, error) {
	if digest, ok := n.semanticLayers[path]; ok {
		return digest, nil
	}
	switch mediaType {
	case "application/vnd.oci.image.layer.v1.tar":
		stats, err := preflightTAR(n.ctx, path, n.limits.MaxEntries-n.count-n.semanticEntries,
			n.limits.MaxExpandedBytes-n.total-n.semanticBytes, n.limits.MaxFileBytes)
		if err != nil {
			return "", err
		}
		if err := n.reserveSemantic(stats.entries, stats.expanded); err != nil {
			return "", err
		}
		digest, err := fileDigestContext(n.ctx, path)
		if err != nil {
			return "", Rejection{Code: RejectMalformed}
		}
		value := "sha256:" + digest
		n.semanticLayers[path] = value
		return value, nil
	case "application/vnd.oci.image.layer.v1.tar+gzip", "application/vnd.docker.image.rootfs.diff.tar.gzip":
		value, err := n.gzipLayerDiffID(path)
		if err == nil {
			n.semanticLayers[path] = value
		}
		return value, err
	default:
		return "", Rejection{Code: RejectUnsupported}
	}
}

func (n *normalizer) gzipLayerDiffID(path string) (string, error) {
	if err := n.reserveSemantic(1, 0); err != nil {
		return "", err
	}
	f, err := os.Open(path)
	if err != nil {
		return "", Rejection{Code: RejectMalformed}
	}
	defer f.Close()
	reader, err := gzip.NewReader(f)
	if err != nil {
		return "", Rejection{Code: RejectMalformed}
	}
	reader.Multistream(false)
	target := filepath.Join(n.staging, "oci-diff-"+strings.Repeat("0", 8)+"-"+hex.EncodeToString([]byte{byte(len(n.result.Entries) % 255)}))
	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		_ = reader.Close()
		return "", Rejection{Code: RejectInvariant}
	}
	h := sha256.New()
	budget := &semanticWriter{normalizer: n, writer: io.MultiWriter(out, h)}
	written, copyErr := io.Copy(budget, io.LimitReader(contextReader{ctx: n.ctx, r: reader}, n.limits.MaxFileBytes+1))
	if closeErr := out.Close(); copyErr == nil {
		copyErr = closeErr
	}
	if closeErr := reader.Close(); copyErr == nil {
		copyErr = closeErr
	}
	if copyErr != nil || written > n.limits.MaxFileBytes {
		_ = os.Remove(target)
		return "", Rejection{Code: RejectResource}
	}
	compressed, err := regularNoLink(path)
	if err != nil || exceedsRatio(written, compressed.Size(), n.limits.MaxCompressionRatio) {
		_ = os.Remove(target)
		return "", Rejection{Code: RejectResource}
	}
	stats, err := preflightTAR(n.ctx, target, n.limits.MaxEntries-n.count-n.semanticEntries,
		n.limits.MaxExpandedBytes-n.total-n.semanticBytes, n.limits.MaxFileBytes)
	if err != nil {
		_ = os.Remove(target)
		return "", err
	}
	if err := n.reserveSemantic(stats.entries, stats.expanded); err != nil {
		_ = os.Remove(target)
		return "", err
	}
	if err := os.Remove(target); err != nil {
		return "", Rejection{Code: RejectInvariant}
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil)), nil
}

type semanticWriter struct {
	normalizer *normalizer
	writer     io.Writer
}

func (w *semanticWriter) Write(value []byte) (int, error) {
	if err := w.normalizer.check(); err != nil {
		return 0, err
	}
	if err := w.normalizer.reserveSemantic(0, int64(len(value))); err != nil {
		return 0, err
	}
	n, err := w.writer.Write(value)
	if n < len(value) {
		w.normalizer.semanticBytes -= int64(len(value) - n)
	}
	return n, err
}

func strictJSONObject(raw []byte) bool {
	if rejectDuplicateJSON(raw) != nil {
		return false
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	var value map[string]json.RawMessage
	return decoder.Decode(&value) == nil && value != nil && decoder.Decode(&struct{}{}) == io.EOF
}
