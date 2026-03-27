package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Release struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	HTMLURL     string `json:"html_url"`
	PublishedAt string `json:"published_at"`
}

func CheckForUpdates(currentVersion string) (*Release, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/fadinflame/simple-ac/releases/latest", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "SimpleAC-Updater")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error checking for updates:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch release: %s", resp.Status)
	}

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var release Release
	if err := json.Unmarshal(raw, &release); err != nil {
		return nil, err
	}

	// Normalize versions for comparison
	vLatest := strings.TrimPrefix(release.TagName, "v")
	vCurrent := strings.TrimPrefix(currentVersion, "v")

	if vLatest != "" && vLatest != vCurrent {
		return &release, nil
	}

	return nil, nil
}
