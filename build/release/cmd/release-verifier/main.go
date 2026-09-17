package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ThameeraDananjaya/project-agnostic-secret-scanner/internal/verify"
)

var (
	releaseToolingCommit = "UNSET"
	releaseToolingTree   = "UNSET"
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
	trustedRoot := flags.String("trusted-root", "", "absolute path to the independently acquired Sigstore trusted root")
	trustedRootSHA256 := flags.String("trusted-root-sha256", "", "expected lowercase SHA-256 recorded at trusted-root acquisition")
	signerPolicy := flags.String("signer-policy", "", "absolute path to independently acquired owner-approved signer policy outside the candidate directory")
	signerPolicyDigest := flags.String("signer-policy-sha256", "", "SHA-256 obtained independently from the owner trust channel, never from candidate assets")
	if err := flags.Parse(args); err != nil || flags.NArg() != 0 || *directory == "" || *cosign == "" || *cosignSHA256 == "" || *trustedRoot == "" || *trustedRootSHA256 == "" ||
		releaseToolingCommit == "UNSET" || releaseToolingTree == "UNSET" {
		fmt.Fprintln(os.Stderr, "verification failed: complete bounded arguments are required")
		return 40
	}
	policy, err := verify.LoadSeparateSignerPolicy(*signerPolicy, *signerPolicyDigest, *directory, releaseToolingCommit, releaseToolingTree)
	if err != nil {
		fmt.Fprintln(os.Stderr, "verification failed: independently pinned signer policy rejected")
		return 40
	}
	signature, err := verify.NewCosignCommandVerifier(*cosign, *cosignSHA256, *trustedRoot, *trustedRootSHA256)
	if err != nil {
		fmt.Fprintln(os.Stderr, "verification failed: Cosign identity or digest rejected")
		return 40
	}
	result, err := verify.VerifyRelease(context.Background(), verify.ReleaseVerificationRequest{
		Directory: *directory, ManifestPath: *manifest, BundlePath: *bundle,
		Policy:    policy,
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
