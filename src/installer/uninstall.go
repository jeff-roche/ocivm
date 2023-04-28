package installer

import (
	"fmt"
	"os"

	"github.com/jeff-roche/ocivm/src/manifest"
)

func UninstallIfInstalled(version string, activeManifest *manifest.VersionManifest) error {
	// Is this a valid version?
	if !activeManifest.ValidVersion(version) {
		return fmt.Errorf("requested version is not a valid version")
	}

	// Make sure the version is installed
	if !activeManifest.Installed(version) {
		return nil
	}

	// Get the path and make sure it exists
	vFolder, exists := getVersionFolderPath(version, activeManifest)
	if !exists {
		return nil
	}

	// Check if we need to deactivate this version first
	if activeManifest.CurrentVersion == version {
		return fmt.Errorf("requested version is currently active, please switch versions before uninstalling")
	}

	if err := os.RemoveAll(vFolder); err != nil {
		return fmt.Errorf("unable to delete the specified version folder: %s", err)
	}

	return nil
}
