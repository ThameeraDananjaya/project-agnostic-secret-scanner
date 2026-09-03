package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/verify"
)

const (
	repository          = "ThameeraDananjaya/project-agnostic-secret-scanner"
	repositoryOwnerID   = int64(50274860)
	workflow            = ".github/workflows/release.yml"
	releaseRef          = "refs/tags/v1.0.0"
	oidcIssuer          = "https://token.actions.githubusercontent.com"
	certificateIdentity = "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/.github/workflows/release.yml@refs/tags/v1.0.0"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	flags := flag.NewFlagSet("scanner-release-verifier", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	directory := flags.String("directory", "", "directory containing the complete acquired release set")
	manifest := flags.String("manifest", "release-manifest.json", "manifest path relative to directory")
	bundle := flags.String("bundle", "release-manifest.sigstore.json", "Cosign bundle path relative to directory")
	cosign := flags.String("cosign", "", "absolute path to the separately acquired pinned Cosign executable")
	cosignSHA256 := flags.String("cosign-sha256", "", "expected lowercase SHA-256 of the Cosign executable")
	if err := flags.Parse(args); err != nil || flags.NArg() != 0 || *directory == "" || *cosign == "" || *cosignSHA256 == "" {
		fmt.Fprintln(os.Stderr, "verification failed: complete bounded arguments are required")
		return 40
	}
	signature, err := verify.NewCosignCommandVerifier(*cosign, *cosignSHA256)
	if err != nil {
		fmt.Fprintln(os.Stderr, "verification failed: Cosign identity or digest rejected")
		return 40
	}
	result, err := verify.VerifyRelease(context.Background(), verify.ReleaseVerificationRequest{
		Directory: *directory, ManifestPath: *manifest, BundlePath: *bundle,
		Policy: verify.ReleaseTrustPolicy{
			Repository: repository, RepositoryOwnerID: repositoryOwnerID, Workflow: workflow,
			Ref: releaseRef, OIDCIssuer: oidcIssuer, CertificateIdentity: certificateIdentity,
			ReleaseVersion: "v1.0.0",
		},
		Signature: signature,
		Now:       time.Now().UTC(),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "verification failed: release evidence rejected")
		return 40
	}
	fmt.Printf("verified release=%s manifest_sha256=%s assets=%d\n", result.ReleaseVersion, result.ManifestDigest, result.AssetCount)
	return 0
}
