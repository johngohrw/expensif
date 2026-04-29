package assets

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
)

// ManifestEntry represents a single entry in Vite's manifest.json.
type ManifestEntry struct {
	File string `json:"file"`
	Src  string `json:"src"`
	Name string `json:"name"`
}

// Manifest maps Vite chunk keys to their built output.
type Manifest map[string]ManifestEntry

// LoadManifest reads and parses a Vite manifest.json file.
func LoadManifest(path string) (Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// AssetHelper provides template functions for injecting Vite assets.
type AssetHelper struct {
	Dev      bool
	Manifest Manifest
}

// ScriptTag returns a <script> tag for the given island entry point.
// In dev mode, it points to the Vite dev server. In production, it reads
// the hashed filename from the manifest.
func (a *AssetHelper) ScriptTag(entry string) template.HTML {
	if a.Dev {
		return template.HTML(fmt.Sprintf(
			`<script type="module" src="http://localhost:8081/src/entries/%s.tsx"></script>`,
			entry,
		))
	}
	for _, entryData := range a.Manifest {
		if entryData.Name == entry {
			return template.HTML(fmt.Sprintf(
				`<script type="module" src="/static/%s"></script>`,
				entryData.File,
			))
		}
	}
	// Fail loudly in production — missing manifest entry is a build bug.
	panic(fmt.Sprintf("manifest entry not found: %s", entry))
}

// DevClient returns the Vite HMR client script tag in development.
// Returns an empty string in production.
func (a *AssetHelper) DevClient() template.HTML {
	if !a.Dev {
		return ""
	}
	return template.HTML(`<script type="module" src="http://localhost:8081/@vite/client"></script>`)
}
