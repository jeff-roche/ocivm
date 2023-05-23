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

var downloadBaseUrlx86 = "https://mirror.openshift.com/pub/openshift-v4/x86_64/clients/ocp"

func GetNewBinaries(requestedVer string, activeManifest *manifest.VersionManifest) error {
	// Run the local check first to avoid network calls if possible
	if activeManifest.Installed(requestedVer) {
		return fmt.Errorf("requested version \"%s\" is already installed", requestedVer)
	}

	if !activeManifest.ValidVersion(requestedVer) {
		return fmt.Errorf("requested version \"%s\" is not a valid version", requestedVer)
	}

	binaries := []string{"openshift-client", "openshift-install"}

	for _, bin := range binaries {
		archive, err := downloadCurrentPlatformBinary(bin, requestedVer)
		if err != nil {
			return fmt.Errorf("unable to download the %s binary archive: %s", bin, archive)
		}

		// Make the destination folder
		dest, _ := getVersionFolderPath(requestedVer, activeManifest)
		os.Mkdir(dest, 0755)

		if err := extractBinaryToFolder(dest, archive); err != nil {
			return fmt.Errorf("unable to extract the %s binary: %s", bin, err)
		}
	}

	return nil
}

func downloadCurrentPlatformBinary(bin, ver string) ([]byte, error) {
	var binName string
	switch runtime.GOOS {
	case "darwin":
		if runtime.GOARCH == "arm64" {
			binName = fmt.Sprintf("%s-mac-arm64-%s.tar.gz", bin, ver)
		} else {
			binName = fmt.Sprintf("%s-mac-%s.tar.gz", bin, ver)
		}
	case "linux":
		binName = fmt.Sprintf("%s-linux-%s.tar.gz", bin, ver)
	case "windows":
		return nil, fmt.Errorf("windows functionality not currently implemented")
	default:
		return nil, fmt.Errorf("no valid binary found for %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	fmt.Printf("Downloading %s version %s for %s/%s...\n", bin, ver, runtime.GOOS, runtime.GOARCH)

	binUrl := fmt.Sprintf("%s/%s/%s", downloadBaseUrlx86, ver, binName)

	return fetchBinZip(binUrl)
}

func fetchBinZip(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve bin archive (%s): %s", url, err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to decode bin archive: %s", err)
	}

	return body, nil
}

func extractBinaryToFolder(dest string, archive []byte) error {
	fmt.Println("Extracting the binary...")

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
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		// The file we are looking for
		switch header.Name {
		case "openshift-install":
			fallthrough
		case "oc":
			fallthrough
		case "kubectl":
			f, err := os.OpenFile(fmt.Sprintf("%s/%s", dest, header.Name), os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("unable to extract the %s binary from the archive: %s", header.Name, err)
			}

			if _, err := io.Copy(f, tr); err != nil {
				return fmt.Errorf("unable to copy the file contents from the archive: %s", err)
			}

			f.Close()
		}
	}
}
