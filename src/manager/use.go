package manager

import (
	"fmt"
	"ocivm/src/manifest"
	"os"
)

// UseVersion will switch the currently active version to a new one
func UseVersion(ver string, m *manifest.VersionManifest) error {
	if !m.Installed(ver) {
		return fmt.Errorf("requested version %s is not installed", ver)
	}

	// basedir/current/openshift-install
	originalPath := fmt.Sprintf("%s/versions/%s/openshift-install", manifest.LocalFolderPath, ver)
	currentLink := fmt.Sprintf("%s/current/openshift-install", manifest.LocalFolderPath)

	// Remove the old link
	if _, err := os.Lstat(currentLink); err == nil {
		os.Remove(currentLink)
	}

	// Make the new link
	if err := os.Link(originalPath, currentLink); err != nil {
		return err
	}

	// Update the .version file
	os.WriteFile(fmt.Sprintf("%s/current/.version", manifest.LocalFolderPath), []byte(ver), 0755)

	return nil
}
