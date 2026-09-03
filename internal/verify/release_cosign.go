package verify

import (
	"context"
	"io"
	"os"
	"os/exec"
)

type CosignCommandVerifier struct {
	path              string
	sha256            string
	trustedRoot       string
	trustedRootSHA256 string
}

func NewCosignCommandVerifier(path, expectedSHA256, trustedRoot, expectedTrustedRootSHA256 string) (*CosignCommandVerifier, error) {
	if !IsDigest(expectedSHA256) || !IsDigest(expectedTrustedRootSHA256) {
		return nil, ErrInvalidReference
	}
	info, err := os.Lstat(path)
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 {
		return nil, ErrInvalidReference
	}
	digest, _, err := digestRegularFile(path)
	if err != nil || digest != expectedSHA256 {
		return nil, ErrBindingMismatch
	}
	rootDigest, _, err := digestRegularFile(trustedRoot)
	if err != nil || rootDigest != expectedTrustedRootSHA256 {
		return nil, ErrBindingMismatch
	}
	return &CosignCommandVerifier{path: path, sha256: expectedSHA256, trustedRoot: trustedRoot, trustedRootSHA256: expectedTrustedRootSHA256}, nil
}

func (v *CosignCommandVerifier) VerifyManifest(ctx context.Context, manifestPath, bundlePath string, identity ReleaseIdentity) error {
	if v == nil || ctx == nil || identity.CertificateIdentity == "" || identity.OIDCIssuer == "" {
		return ErrInvalidReference
	}
	refreshed, err := NewCosignCommandVerifier(v.path, v.sha256, v.trustedRoot, v.trustedRootSHA256)
	if err != nil || refreshed.path != v.path {
		return ErrBindingMismatch
	}
	command := exec.CommandContext(ctx, v.path,
		"verify-blob",
		"--bundle", bundlePath,
		"--trusted-root", v.trustedRoot,
		"--certificate-identity", identity.CertificateIdentity,
		"--certificate-oidc-issuer", identity.OIDCIssuer,
		manifestPath,
	)
	command.Stdin = nil
	command.Stdout = io.Discard
	command.Stderr = io.Discard
	command.Env = []string{"COSIGN_YES=false"}
	if err := command.Run(); err != nil {
		return ErrInvalidSignature
	}
	return nil
}
