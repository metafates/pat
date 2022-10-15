package upgrader

import (
	"encoding/json"
	"fmt"
	"github.com/metafates/pat/constant"
	"github.com/samber/mo"
	"net/http"
)

var cachedLatestVersion = mo.None[string]()

// LatestVersion returns the latest version of mangal.
// It will fetch the latest version from the GitHub API.
func LatestVersion() (version string, err error) {
	if cachedLatestVersion.IsPresent() {
		return cachedLatestVersion.MustGet(), nil
	}

	resp, err := http.Get(
		fmt.Sprintf("https://api.github.com/repos/metafates/%s/releases/latest", constant.App),
	)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
	}

	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		return
	}

	// remove the v from the tag name
	if release.TagName == "" {
		return constant.Version, nil
	}

	// remove the v from the tag name
	version = release.TagName[1:]
	cachedLatestVersion = mo.Some(version)
	return
}
