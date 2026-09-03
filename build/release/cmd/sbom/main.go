package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"sort"
	"strings"
)

type module struct {
	Path    string  `json:"Path"`
	Version string  `json:"Version"`
	Sum     string  `json:"Sum"`
	Main    bool    `json:"Main"`
	Replace *module `json:"Replace"`
}

type packageEntry struct {
	SPDXID           string              `json:"SPDXID"`
	Name             string              `json:"name"`
	VersionInfo      string              `json:"versionInfo"`
	DownloadLocation string              `json:"downloadLocation"`
	FilesAnalyzed    bool                `json:"filesAnalyzed"`
	ExternalRefs     []map[string]string `json:"externalRefs,omitempty"`
}

type relationship struct {
	SPDXElementID      string `json:"spdxElementId"`
	RelationshipType   string `json:"relationshipType"`
	RelatedSPDXElement string `json:"relatedSpdxElement"`
}

type document struct {
	SPDXVersion       string         `json:"spdxVersion"`
	DataLicense       string         `json:"dataLicense"`
	SPDXID            string         `json:"SPDXID"`
	Name              string         `json:"name"`
	DocumentNamespace string         `json:"documentNamespace"`
	CreationInfo      map[string]any `json:"creationInfo"`
	Packages          []packageEntry `json:"packages"`
	Relationships     []relationship `json:"relationships"`
}

func main() {
	input := flag.String("input", "", "go list -m -json all stream")
	output := flag.String("output", "", "SPDX 2.3 JSON output")
	revision := flag.String("revision", "", "exact 40-character source revision")
	created := flag.String("created", "", "canonical UTC timestamp")
	flag.Parse()
	if flag.NArg() != 0 || *input == "" || *output == "" || len(*revision) != 40 || *created == "" {
		fatal("complete deterministic SBOM arguments are required")
	}
	modules, err := readModules(*input)
	if err != nil {
		fatal(err.Error())
	}
	packages := make([]packageEntry, 0, len(modules))
	var relations []relationship
	rootID := "SPDXRef-Package-project-agnostic-secret-scanner"
	for index, value := range modules {
		version := value.Version
		if value.Main {
			version = "v1.0.0+" + *revision
		}
		id := rootID
		if !value.Main {
			id = fmt.Sprintf("SPDXRef-Package-%03d", index+1)
		}
		purl := "pkg:golang/" + url.PathEscape(value.Path)
		if version != "" {
			purl += "@" + url.PathEscape(version)
		}
		packages = append(packages, packageEntry{SPDXID: id, Name: value.Path, VersionInfo: version, DownloadLocation: "NOASSERTION", FilesAnalyzed: false,
			ExternalRefs: []map[string]string{{"referenceCategory": "PACKAGE-MANAGER", "referenceType": "purl", "referenceLocator": purl}}})
		if !value.Main {
			relations = append(relations, relationship{SPDXElementID: rootID, RelationshipType: "DEPENDS_ON", RelatedSPDXElement: id})
		}
	}
	relations = append([]relationship{{SPDXElementID: "SPDXRef-DOCUMENT", RelationshipType: "DESCRIBES", RelatedSPDXElement: rootID}}, relations...)
	value := document{
		SPDXVersion: "SPDX-2.3", DataLicense: "CC0-1.0", SPDXID: "SPDXRef-DOCUMENT",
		Name: "project-agnostic-secret-scanner-v1.0.0", DocumentNamespace: "https://github.com/ThameeraDananjaya/project-agnostic-secret-scanner/releases/v1.0.0/spdx/" + *revision,
		CreationInfo: map[string]any{"created": *created, "creators": []string{"Tool: scanner-release-builder-1.0.0"}},
		Packages:     packages, Relationships: relations,
	}
	raw, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fatal(err.Error())
	}
	raw = append(raw, '\n')
	if err := os.WriteFile(*output, raw, 0o644); err != nil {
		fatal(err.Error())
	}
}

func readModules(path string) ([]module, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(io.LimitReader(file, 16<<20))
	var values []module
	for {
		var value module
		if err := decoder.Decode(&value); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		if value.Replace != nil {
			value = *value.Replace
		}
		if value.Path == "" || strings.ContainsAny(value.Path, "\r\n\x00") {
			return nil, fmt.Errorf("invalid module identity")
		}
		values = append(values, value)
	}
	if len(values) == 0 {
		return nil, fmt.Errorf("empty module graph")
	}
	sort.Slice(values, func(i, j int) bool {
		if values[i].Main != values[j].Main {
			return values[i].Main
		}
		return values[i].Path < values[j].Path
	})
	return values, nil
}

func fatal(message string) {
	fmt.Fprintln(os.Stderr, "sbom:", message)
	os.Exit(1)
}
