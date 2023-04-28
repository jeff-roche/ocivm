package installer

import (
	"fmt"
	"os"

	"github.com/jeff-roche/ocivm/src/manifest"
)

// getVersionFolderPath will generate the path to the folder for the specified version
//
//	A path and a boolean specifying if the folder exists will be returned
func getVersionFolderPath(version string, activeManifest *manifest.VersionManifest) (string, bool) {
	dest := fmt.Sprintf("%s/versions/%s", manifest.LocalFolderPath, version)

	exists := false
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		exists = true
	}

	return dest, exists
}
