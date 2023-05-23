package manager

import (
	"fmt"
	"os"

	"github.com/jeff-roche/ocivm/src/manifest"
)

// UseVersion will switch the currently active version to a new one
func UseVersion(ver string, m *manifest.VersionManifest) error {
	if !m.Installed(ver) {
		return fmt.Errorf("requested version %s is not installed", ver)
	}

	bins := []string{"oc", "kubectl", "openshift-install"}

	for _, bin := range bins {
		// basedir/current/bin-name
		originalPath := fmt.Sprintf("%s/versions/%s/%s", manifest.LocalFolderPath, ver, bin)
		currentLink := fmt.Sprintf("%s/current/%s", manifest.LocalFolderPath, bin)

		// Remove the old link
		if _, err := os.Lstat(currentLink); err == nil {
			os.Remove(currentLink)
		}

		// Make the new link
		if err := os.Link(originalPath, currentLink); err != nil {
			return err
		}
	}

	// Update the .version file
	os.WriteFile(fmt.Sprintf("%s/current/.version", manifest.LocalFolderPath), []byte(ver), 0755)

	return nil
}
