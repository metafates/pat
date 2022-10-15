package upgrader

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/metafates/pat/color"
	"github.com/metafates/pat/constant"
	"github.com/metafates/pat/filesystem"
	"github.com/metafates/pat/icon"
	"github.com/metafates/pat/log"
	"github.com/metafates/pat/util"
	"github.com/metafates/pat/where"
	"github.com/samber/lo"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"text/template"
)

var erase func()

func info(format string, args ...any) (erase func()) {
	return util.PrintErasable(
		fmt.Sprintf(
			"%s %s",
			lipgloss.NewStyle().Foreground(color.Blue).Render(icon.Progress),
			fmt.Sprintf(format, args...),
		),
	)
}

// Upgrade updates mangal to the latest version.
func Upgrade() (err error) {
	erase = info("Fetching latest version")
	version, err := LatestVersion()
	if err != nil {
		return
	}

	erase()

	comp, err := util.CompareSemVers(constant.Version, version)
	if err != nil {
		return err
	}

	if comp >= 0 {
		fmt.Printf(
			"%s %s %s\n",
			lipgloss.NewStyle().Foreground(color.Green).Render("Congrats!"),
			"You're already on the latest version of "+constant.App,
			lipgloss.NewStyle().Faint(true).Render(fmt.Sprintf("(which is %s)", constant.Version)),
		)
		return
	}

	fmt.Printf(
		"%s %s is out! You're on %s\n",
		lipgloss.NewStyle().Bold(true).Render(constant.App),
		lipgloss.NewStyle().Foreground(color.Cyan).Render(version),
		lipgloss.NewStyle().Foreground(color.Blue).Render(constant.Version),
	)

	err = update()

	if err != nil {
		return
	}

	postUpgradeTemplate := lo.Must(template.New("post-upgrade").Funcs(template.FuncMap{
		"cyan":  lipgloss.NewStyle().Foreground(color.Cyan).Render,
		"faint": lipgloss.NewStyle().Faint(true).Render,
		"title": lipgloss.NewStyle().Bold(true).Foreground(color.Green).Render,
	}).Parse(`
Welcome to {{ .App }} v{{ .Version }}

Report any bugs:

    https://github.com/metafates/{{ .App }}/issues

What's new:

    https://github.com/metafates/{{ .App }}/releases/tag/v{{ .Version }}

Changelog:

    https://github.com/metafates/{{ .App }}/compare/v{{ .OldVersion }}...v{{ .Version }}
`))

	_ = postUpgradeTemplate.Execute(os.Stdout, struct {
		App        string
		Version    string
		OldVersion string
	}{
		App:        constant.App,
		Version:    version,
		OldVersion: constant.Version,
	})

	return
}

// update self-update binary
func update() (err error) {
	erase()
	log.Infof("upgrading %s to the latest version", constant.App)

	var (
		version     string
		arch        string
		selfPath    string
		archiveName string
		archiveType string
	)

	if selfPath, err = os.Executable(); err != nil {
		return
	}

	if version, err = LatestVersion(); err != nil {
		return
	}

	switch runtime.GOARCH {
	case "amd64":
		arch = "x86_64"
	case "386":
		arch = "i386"
	default:
		arch = runtime.GOARCH
	}

	archiveType = "tar.gz"
	archiveName = fmt.Sprintf("%s_%s_%s_%s.%s", constant.App, version, util.Capitalize(runtime.GOOS), arch, archiveType)
	url := fmt.Sprintf(
		"https://github.com/metafates/%s/releases/download/v%s/%s",
		constant.App,
		version,
		archiveName,
	)

	erase = info("Downloading %s", lipgloss.NewStyle().Foreground(color.Yellow).Render(url))

	res, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return err
	}

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("error downloading binary: status code %d", res.StatusCode)
		log.Error(err)
		return
	}

	defer res.Body.Close()

	erase()
	erase = info("Extracting %s", lipgloss.NewStyle().Foreground(color.Yellow).Render(archiveName))
	out := filepath.Join(where.Temp(), "mangal_update")

	err = util.ExtractTarTo(res.Body, out)

	if err != nil {
		log.Error(err)
		return err
	}

	erase()
	erase = info("Replacing %s", lipgloss.NewStyle().Foreground(color.Yellow).Render(selfPath))
	// remove the old binary
	// it should not interrupt the running process
	file, err := filesystem.Api().OpenFile(selfPath, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Error(err)
		err = errors.New("error removing old binary, try running this as a root user")
		return err
	}

	newPat, err := filesystem.Api().OpenFile(filepath.Join(out, constant.App), os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error(err)
		return
	}

	stat, err := newPat.Stat()
	if err != nil {
		log.Error(err)
		return
	}

	if stat.Size() == 0 {
		log.Error(err)
		return err
	}

	_, err = io.Copy(file, newPat)
	if err != nil {
		log.Error(err)
		return err
	}

	erase()

	return
}
