package installer

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/jeff-roche/ocivm/src/manifest"
)

var installerBaseUrlx86 = "https://mirror.openshift.com/pub/openshift-v4/x86_64/clients/ocp"

func GetNewInstaller(requestedVer string, activeManifest *manifest.VersionManifest) error {
	// Run the local check first to avoid network calls if possible
	if activeManifest.Installed(requestedVer) {
		return fmt.Errorf("requested version \"%s\" is already installed", requestedVer)
	}

	if !activeManifest.ValidVersion(requestedVer) {
		return fmt.Errorf("requested version \"%s\" is not a valid version", requestedVer)
	}

	archive, err := downloadCurrentPlatformInstaller(requestedVer)
	if err != nil {
		return fmt.Errorf("unable to download the installer archive: %s", archive)
	}

	// Make the destination folder
	dest, _ := getVersionFolderPath(requestedVer, activeManifest)
	os.Mkdir(dest, 0755)

	if err := extractInstallerToFolder(dest, archive); err != nil {
		return fmt.Errorf("unable to extract the installer: %s", err)
	}

	return nil
}

func downloadCurrentPlatformInstaller(ver string) ([]byte, error) {
	var installerName string
	switch runtime.GOOS {
	case "darwin":
		if runtime.GOARCH == "arm64" {
			installerName = fmt.Sprintf("openshift-install-mac-arm64-%s.tar.gz", ver)
		} else {
			installerName = fmt.Sprintf("openshift-install-mac-%s.tar.gz", ver)
		}
	case "linux":
		installerName = fmt.Sprintf("openshift-install-linux-%s.tar.gz", ver)
	case "windows":
		return nil, fmt.Errorf("windows functionality not currently implemented")
	default:
		return nil, fmt.Errorf("no valid binary found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	fmt.Printf("Downloading version %s for %s/%s...\n", ver, runtime.GOOS, runtime.GOARCH)

	installerUrl := fmt.Sprintf("%s/%s/%s", installerBaseUrlx86, ver, installerName)

	return fetchInstaller(installerUrl)
}

func fetchInstaller(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve installer archive (%s): %s", url, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to decode installer archive: %s", err)
	}

	return body, nil
}

func extractInstallerToFolder(dest string, archive []byte) error {
	fmt.Println("Extracting the installer...")

	// Unzip the archive
	unzipped, err := gzip.NewReader(bytes.NewReader(archive))
	if err != nil {
		return fmt.Errorf("unable to unzip the archive: %s", err)
	}
	defer unzipped.Close()

	// Parse the tar
	tr := tar.NewReader(unzipped)

	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF: // End of the archive
			return fmt.Errorf("could not find the installer in the archive")
		case err != nil:
			return err
		case header == nil:
			continue
		}

		// The file we are looking for
		if header.Name == "openshift-install" {
			f, err := os.OpenFile(fmt.Sprintf("%s/openshift-install", dest), os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("unable to extract the openshift-installer from the archive: %s", err)
			}

			if _, err := io.Copy(f, tr); err != nil {
				return fmt.Errorf("unable to copy the file contents from the archive: %s", err)
			}

			f.Close()

			break
		}
	}

	return nil
}
