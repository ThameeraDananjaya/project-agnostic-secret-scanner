package artifact

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

var ociDigest = regexp.MustCompile(`^sha256:([0-9a-f]{64})$`)

type descriptor struct {
	MediaType string   `json:"mediaType"`
	Digest    string   `json:"digest"`
	Size      int64    `json:"size"`
	URLs      []string `json:"urls,omitempty"`
	Data      string   `json:"data,omitempty"`
}

type ociIndex struct {
	SchemaVersion int          `json:"schemaVersion"`
	Manifests     []descriptor `json:"manifests"`
}

type ociManifest struct {
	SchemaVersion int          `json:"schemaVersion"`
	Config        descriptor   `json:"config"`
	Layers        []descriptor `json:"layers"`
}

func isOCILayout(root string) bool {
	_, layoutErr := os.Lstat(filepath.Join(root, "oci-layout"))
	_, indexErr := os.Lstat(filepath.Join(root, "index.json"))
	return layoutErr == nil || indexErr == nil
}

func (n *normalizer) validateOCI(root string) error {
	layoutRaw, err := os.ReadFile(filepath.Join(root, "oci-layout"))
	if err != nil || len(layoutRaw) > 4096 {
		return Rejection{Code: RejectMalformed}
	}
	var layout struct {
		Version string `json:"imageLayoutVersion"`
	}
	if strictJSON(layoutRaw, &layout) != nil || layout.Version != "1.0.0" {
		return Rejection{Code: RejectMalformed}
	}
	indexRaw, err := os.ReadFile(filepath.Join(root, "index.json"))
	if err != nil || len(indexRaw) > 16<<20 {
		return Rejection{Code: RejectMalformed}
	}
	var index ociIndex
	if strictJSON(indexRaw, &index) != nil || index.SchemaVersion != 2 || len(index.Manifests) == 0 || len(index.Manifests) > n.limits.MaxEntries {
		return Rejection{Code: RejectMalformed}
	}
	referenced := map[string]bool{"oci-layout": true, "index.json": true}
	bindings := map[string]descriptor{}
	queue := append([]descriptor(nil), index.Manifests...)
	for len(queue) > 0 {
		if err := n.check(); err != nil {
			return err
		}
		d := queue[0]
		queue = queue[1:]
		if len(referenced)+len(queue) > n.limits.MaxEntries {
			return Rejection{Code: RejectResource}
		}
		match := ociDigest.FindStringSubmatch(d.Digest)
		if len(match) != 2 || d.Size < 0 || d.Size > n.limits.MaxFileBytes || len(d.URLs) != 0 || d.Data != "" {
			return Rejection{Code: RejectMalformed}
		}
		relative := filepath.ToSlash(filepath.Join("blobs", "sha256", match[1]))
		if referenced[relative] {
			bound := bindings[relative]
			if bound.Size != d.Size || bound.MediaType != d.MediaType {
				return Rejection{Code: RejectMalformed}
			}
			continue
		}
		path := filepath.Join(root, filepath.FromSlash(relative))
		info, err := regularNoLink(path)
		if err != nil || info.Size() != d.Size {
			return Rejection{Code: RejectMalformed}
		}
		digest, err := fileDigestContext(n.ctx, path)
		if err != nil || digest != match[1] {
			return Rejection{Code: RejectMalformed}
		}
		referenced[relative] = true
		bindings[relative] = d
		if d.MediaType == "application/vnd.oci.image.manifest.v1+json" || d.MediaType == "application/vnd.docker.distribution.manifest.v2+json" {
			if d.Size > 16<<20 {
				return Rejection{Code: RejectResource}
			}
			raw, readErr := os.ReadFile(path)
			var manifest ociManifest
			if readErr != nil || strictJSON(raw, &manifest) != nil || manifest.SchemaVersion != 2 || len(manifest.Layers) == 0 || len(manifest.Layers) > n.limits.MaxEntries-len(referenced) {
				return Rejection{Code: RejectMalformed}
			}
			if err := n.validateOCIImage(root, manifest); err != nil {
				return err
			}
			queue = append(queue, manifest.Config)
			queue = append(queue, manifest.Layers...)
		} else {
			if !supportedOCIBlobMediaType(d.MediaType) {
				return Rejection{Code: RejectUnsupported}
			}
			if err := n.validateOCIBlobSemantics(path, d.MediaType); err != nil {
				return err
			}
		}
	}
	actual := map[string]bool{}
	err = filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return Rejection{Code: RejectUnsafe}
		}
		if path == root || entry.IsDir() {
			return nil
		}
		if err := n.check(); err != nil {
			return err
		}
		if len(actual) >= n.limits.MaxEntries {
			return Rejection{Code: RejectResource}
		}
		relative, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return Rejection{Code: RejectUnsafe}
		}
		actual[filepath.ToSlash(relative)] = true
		return nil
	})
	if err != nil || len(actual) != len(referenced) {
		return Rejection{Code: RejectMalformed}
	}
	for path := range referenced {
		if !actual[path] {
			return Rejection{Code: RejectMalformed}
		}
	}
	return nil
}

func (n *normalizer) validateOCIImage(root string, manifest ociManifest) error {
	if manifest.Config.MediaType != "application/vnd.oci.image.config.v1+json" &&
		manifest.Config.MediaType != "application/vnd.docker.container.image.v1+json" {
		return Rejection{Code: RejectMalformed}
	}
	configPath, err := n.resolveOCIDescriptor(root, manifest.Config)
	if err != nil {
		return err
	}
	configRaw, err := os.ReadFile(configPath)
	if err != nil || len(configRaw) > 16<<20 {
		return Rejection{Code: RejectMalformed}
	}
	diffIDs := make([]string, 0, len(manifest.Layers))
	for _, layer := range manifest.Layers {
		layerPath, err := n.resolveOCIDescriptor(root, layer)
		if err != nil {
			return err
		}
		diffID, err := n.layerDiffID(layerPath, layer.MediaType)
		if err != nil {
			return err
		}
		diffIDs = append(diffIDs, diffID)
	}
	return validateImageConfiguration(configRaw, diffIDs)
}

func (n *normalizer) resolveOCIDescriptor(root string, d descriptor) (string, error) {
	match := ociDigest.FindStringSubmatch(d.Digest)
	if len(match) != 2 || d.Size < 0 || d.Size > n.limits.MaxFileBytes || len(d.URLs) != 0 || d.Data != "" {
		return "", Rejection{Code: RejectMalformed}
	}
	path := filepath.Join(root, "blobs", "sha256", match[1])
	info, err := regularNoLink(path)
	if err != nil || info.Size() != d.Size {
		return "", Rejection{Code: RejectMalformed}
	}
	digest, err := fileDigestContext(n.ctx, path)
	if err != nil || digest != match[1] {
		return "", Rejection{Code: RejectMalformed}
	}
	return path, nil
}

func (n *normalizer) validateOCIBlobSemantics(path, mediaType string) error {
	class, err := classifyNamedContext(n.ctx, path, ociLogicalName(mediaType))
	if err != nil {
		return err
	}
	switch mediaType {
	case "application/vnd.oci.image.config.v1+json", "application/vnd.docker.container.image.v1+json":
		raw, readErr := os.ReadFile(path)
		if readErr != nil || class != ClassFile || !jsonObject(raw) {
			return Rejection{Code: RejectMalformed}
		}
	case "application/vnd.oci.image.layer.v1.tar":
		if class != ClassTAR {
			return Rejection{Code: RejectMalformed}
		}
	case "application/vnd.oci.image.layer.v1.tar+gzip", "application/vnd.docker.image.rootfs.diff.tar.gzip":
		if class != ClassGZIP || !gzipStartsWithTAR(n.ctx, path) {
			return Rejection{Code: RejectMalformed}
		}
	default:
		return Rejection{Code: RejectUnsupported}
	}
	return nil
}

func ociLogicalName(mediaType string) string {
	if mediaType == "application/vnd.oci.image.layer.v1.tar" {
		return "layer.tar"
	}
	if mediaType == "application/vnd.oci.image.layer.v1.tar+gzip" || mediaType == "application/vnd.docker.image.rootfs.diff.tar.gzip" {
		return "layer.tar.gz"
	}
	return "config.json"
}

func gzipStartsWithTAR(ctx context.Context, path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	reader, err := gzip.NewReader(f)
	if err != nil {
		return false
	}
	defer reader.Close()
	header := make([]byte, 512)
	_, err = io.ReadFull(contextReader{ctx: ctx, r: reader}, header)
	return err == nil && looksLikeTARHeader(header)
}

func supportedOCIBlobMediaType(value string) bool {
	switch value {
	case "application/vnd.oci.image.config.v1+json",
		"application/vnd.docker.container.image.v1+json",
		"application/vnd.oci.image.layer.v1.tar",
		"application/vnd.oci.image.layer.v1.tar+gzip",
		"application/vnd.docker.image.rootfs.diff.tar.gzip":
		return true
	default:
		return false
	}
}

func strictJSON(raw []byte, target any) error {
	if err := rejectDuplicateJSON(raw); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return Rejection{Code: RejectMalformed}
	}
	return nil
}

func ociDescriptorFor(mediaType string, raw []byte) descriptor {
	sum := sha256.Sum256(raw)
	return descriptor{MediaType: mediaType, Digest: "sha256:" + hex.EncodeToString(sum[:]), Size: int64(len(raw))}
}
