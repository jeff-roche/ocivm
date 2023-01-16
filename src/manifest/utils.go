package manifest

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func initializeFolderStructure() error {
	// Current version folder
	if err := os.MkdirAll(fmt.Sprintf("%s/current", LocalFolderPath), 0755); err != nil {
		return fmt.Errorf("unable to setup the folder for the current version: %s", err)
	}

	// Installed versions folder
	if err := os.MkdirAll(fmt.Sprintf("%s/versions", LocalFolderPath), 0755); err != nil {
		return fmt.Errorf("unable to setup the folder for installed versions: %s", err)
	}

	return nil
}

func getRemoteHtmlData(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve page data: %s", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to parse page data: %s", err)
	}

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("unable to parse page html: %s", err)
	}

	return doc, nil
}

func parseHtmlPageForVersionNumbers(root *html.Node, versions *[]string) {
	if root.FirstChild == nil {
		return
	}

	for c := root.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode && c.Parent.Data == "span" {
			// See if this is the right span
			for _, a := range c.Parent.Attr {
				if a.Key == "class" && a.Val == "name" {
					*versions = append(*versions, c.Data)
				}
			}
		}

		parseHtmlPageForVersionNumbers(c, versions)
	}
}

func printVersionList(versions []string, current string) {
	for _, ver := range versions {
		currentLabel := ""
		if current != "" && ver == current {
			currentLabel = " (current)"
		}

		fmt.Printf("%s%s\n", ver, currentLabel)
	}
}
