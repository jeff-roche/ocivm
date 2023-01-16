package manifest

import (
	"fmt"
	"os"
	"strings"
)

var (
	LocalFolderPath          string = ""
	versionsFolderPath       string = ""
	currentVersionFolderPath string = ""
	RemoteVersionUrl         string = "https://mirror.openshift.com/pub/openshift-v4/x86_64/clients/ocp/"
)

type VersionManifest struct {
	RemoteVersions []string
	LocalVersions  []string
	CurrentVersion string
}

func init() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("unable to locate users home directory:" + err.Error())
		os.Exit(1)
	}

	LocalFolderPath = fmt.Sprintf("%s/.ocivm", homedir)

	if _, err := os.Stat(LocalFolderPath); os.IsNotExist(err) {
		if err := initializeFolderStructure(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Successfully initialized the configuration folder (%s)\n", LocalFolderPath)
	}

	versionsFolderPath = fmt.Sprintf("%s/versions", LocalFolderPath)
	currentVersionFolderPath = fmt.Sprintf("%s/current", LocalFolderPath)
}

// RefreshVersionInfo scans the local directory and updates the list of locally installed versions and the currently active version
func (m *VersionManifest) RefreshVersionInfo() error {
	if err := m.getLocalVersions(); err != nil {
		return fmt.Errorf("unable to refresh local version list: %s", err)
	}

	if err := m.getCurrentVersion(); err != nil {
		return fmt.Errorf("unable to refresh the current version: %s", err)
	}

	return nil
}

// getRemoteVersions will reach out to the public mirror and pull the list of versions
func (m *VersionManifest) fetchRemoteVersions() error {
	// Get the document
	doc, err := getRemoteHtmlData(RemoteVersionUrl)
	if err != nil {
		return err
	}

	// Extract version numbers
	m.RemoteVersions = []string{}
	parseHtmlPageForVersionNumbers(doc, &m.RemoteVersions)

	return nil
}

func (m *VersionManifest) ListVersions(current, remote bool) {
	// Current version marking
	currentVer := ""
	if current {
		currentVer = m.CurrentVersion
	}

	// Do we need to fetch what versions are available on the remote?
	if remote {
		if err := m.fetchRemoteVersions(); err != nil {
			fmt.Printf("unable to update remote listing: %s\n", err)
			os.Exit(1)
		}

		fmt.Println("Available Versions:")
		printVersionList(m.RemoteVersions, currentVer)
	} else {
		fmt.Println("Installed versions:")
		printVersionList(m.LocalVersions, currentVer)
	}
}

func (m *VersionManifest) Installed(ver string) bool {
	if len(m.LocalVersions) == 0 {
		m.RefreshVersionInfo()
	}

	for _, lv := range m.LocalVersions {
		fmt.Println(lv)
		if lv == ver {
			return true
		}
	}

	return false
}

func (m *VersionManifest) ValidVersion(ver string) bool {
	if len(m.RemoteVersions) == 0 {
		m.fetchRemoteVersions()
	}

	for _, rv := range m.RemoteVersions {
		if rv == ver {
			return true
		}
	}

	return false
}

// getLocalVersions will check the local
func (m *VersionManifest) getLocalVersions() error {
	m.LocalVersions = []string{}

	// Validate the versions folder exists
	if _, err := os.Stat(versionsFolderPath); os.IsNotExist(err) {
		return nil
	}

	dirContents, err := os.ReadDir(versionsFolderPath)
	if err != nil {
		return fmt.Errorf("unable to check installed versions: %s", err)
	}

	for _, item := range dirContents {
		if item.IsDir() {
			m.LocalVersions = append(m.LocalVersions, item.Name())
		}
	}

	return nil
}

func (m *VersionManifest) getCurrentVersion() error {
	m.CurrentVersion = ""

	// Validate the current version file exists
	verFile := fmt.Sprintf("%s/.version", currentVersionFolderPath)
	if _, err := os.Stat(verFile); !os.IsNotExist(err) {
		ver, err := os.ReadFile(verFile)
		if err != nil {
			return fmt.Errorf("unable to parse current version: %s", err)
		}

		m.CurrentVersion = strings.TrimSpace(string(ver))
	}

	return nil
}
