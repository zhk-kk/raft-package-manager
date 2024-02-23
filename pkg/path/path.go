package path

import "path"

var (
	MetadataDir           = "metadata"
	IgnoreDir             = ".ignore"
	CopyDataDir           = "cpdata"
	IntegrationScriptsDir = "iscripts"
	ManifestFile          = path.Join(MetadataDir, "manifest.json")
)
