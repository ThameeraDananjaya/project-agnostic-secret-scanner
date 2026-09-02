package request

type EngineBinding struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	BinaryDigest   string `json:"binaryDigest"`
	AdapterVersion string `json:"adapterVersion"`
}

type SourceBinding struct {
	BaseCommit         string `json:"baseCommit,omitempty"`
	HeadCommit         string `json:"headCommit"`
	MergeBase          string `json:"mergeBase,omitempty"`
	FirstRelease       bool   `json:"firstRelease"`
	HistoryRangeDigest string `json:"historyRangeDigest"`
	TrackedTreeDigest  string `json:"trackedTreeDigest"`
}

type FileBinding struct {
	Path   string `json:"path"`
	Digest string `json:"digest"`
}

type ArtifactEntry struct {
	Path   string `json:"path"`
	Type   string `json:"type"`
	Size   int64  `json:"size"`
	Digest string `json:"digest"`
}

type ArtifactManifest struct {
	Digest  string          `json:"digest"`
	Entries []ArtifactEntry `json:"entries"`
}

type FallbackRequirement struct {
	Mode string `json:"mode"`
}

type SignatureReference struct {
	TrustDomain string `json:"trustDomain"`
	Algorithm   string `json:"algorithm"`
	KeyID       string `json:"keyId"`
	Value       string `json:"value"`
}

type ProjectDocumentBinding struct {
	SchemaFamily   string             `json:"schemaFamily"`
	SchemaVersion  string             `json:"schemaVersion"`
	AdapterVersion string             `json:"adapterVersion"`
	Digest         string             `json:"digest"`
	Signature      SignatureReference `json:"signature"`
}

type ProofBindings struct {
	AdmissionLedgerSchemaVersion string `json:"admissionLedgerSchemaVersion"`
	AdmissionLedgerDigest        string `json:"admissionLedgerDigest"`
	RawClassifierVersion         string `json:"rawClassifierVersion"`
	RawClassifierDigest          string `json:"rawClassifierDigest"`
	PreparationVersion           string `json:"preparationVersion"`
	PreparationDigest            string `json:"preparationDigest"`
	ResourceProfileID            string `json:"resourceProfileId"`
	InspectionProofFormat        string `json:"inspectionProofFormat"`
	InspectionProofDigest        string `json:"inspectionProofDigest"`
}

type Limits struct {
	TimeoutSeconds      int   `json:"timeoutSeconds"`
	MaxArchiveDepth     int   `json:"maxArchiveDepth"`
	MaxArchiveEntries   int   `json:"maxArchiveEntries"`
	MaxExpandedBytes    int64 `json:"maxExpandedBytes"`
	MaxFileBytes        int64 `json:"maxFileBytes"`
	MaxCompressionRatio int   `json:"maxCompressionRatio"`
	MaxMemoryBytes      int64 `json:"maxMemoryBytes"`
	MaxCPUPercent       int   `json:"maxCpuPercent"`
}

type ScanRequest struct {
	RequestSchemaVersion  string                  `json:"requestSchemaVersion"`
	RequiredFeatures      []string                `json:"requiredFeatures,omitempty"`
	ScanID                string                  `json:"scanId"`
	SupersedesScanID      string                  `json:"supersedesScanId,omitempty"`
	Mode                  string                  `json:"mode"`
	ScannerReleaseDigest  string                  `json:"scannerReleaseDigest"`
	EngineBinding         EngineBinding           `json:"engineBinding"`
	RulePackDigest        string                  `json:"rulePackDigest"`
	PolicyDigest          string                  `json:"policyDigest"`
	AllowlistDigest       string                  `json:"allowlistDigest"`
	PolicyBinding         *ProjectDocumentBinding `json:"policyBinding,omitempty"`
	AllowlistBinding      *ProjectDocumentBinding `json:"allowlistBinding,omitempty"`
	ProofBindings         *ProofBindings          `json:"proofBindings,omitempty"`
	SourceBinding         SourceBinding           `json:"sourceBinding"`
	TrackedSourceManifest FileBinding             `json:"trackedSourceManifest"`
	BuildContextManifest  *FileBinding            `json:"buildContextManifest,omitempty"`
	ArtifactManifest      *ArtifactManifest       `json:"artifactManifest,omitempty"`
	FallbackRequirement   FallbackRequirement     `json:"fallbackRequirement"`
	Limits                Limits                  `json:"limits"`
	OfflineRequired       bool                    `json:"offlineRequired"`
	RedactionMode         string                  `json:"redactionMode"`
	RequestedAt           string                  `json:"requestedAt"`
}
