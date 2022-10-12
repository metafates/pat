package shell

import (
	"fmt"
	"github.com/metafates/pat/filesystem"
	"github.com/metafates/pat/log"
	"github.com/metafates/pat/shell/bash"
	"github.com/metafates/pat/shell/fish"
	"github.com/metafates/pat/shell/sh"
	"github.com/metafates/pat/shell/zsh"
	"github.com/metafates/pat/where"
	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"time"
)

type Shell interface {
	Name() string
	AddPath(path string) error
	RemovePath(path string) error
	Overwrite(paths []string) error
	Paths() ([]string, error)
	Available() bool
}

func AvailableShells() []Shell {
	shells := lo.Filter([]Shell{
		fish.New(),
		zsh.New(),
		bash.New(),
		sh.New(),
	}, func(shell Shell, _ int) bool {
		return shell.Available()
	})

	log.Infof("found %d available shells", len(shells))
	return shells
}

func Backup(shell Shell) {
	paths, err := shell.Paths()
	if err != nil {
		log.Errorf("failed to backup paths: %v", err)
		return
	}

	log.Infof("backing up %d paths", len(paths))
	fileName := filepath.Join(where.Backup(), fmt.Sprintf("backup.%s.toml", shell.Name()))

	type Backup struct {
		Date  string   `toml:"date"`
		Paths []string `toml:"paths" multiline:"true"`
	}

	type Backups struct {
		Backups []Backup `toml:"backups"`
	}

	b := Backup{
		Date:  time.Now().Format("15:04:05 02 Jan 2006"),
		Paths: paths,
	}

	var (
		backups  Backups
		contents []byte
	)

	exists, err := filesystem.Api().Exists(fileName)
	if err != nil {
		log.Errorf("failed to check if backup file exists: %v", err)
		return
	}

	if exists {
		contents, err = filesystem.Api().ReadFile(fileName)
		if err != nil {
			log.Errorf("failed to read backup file: %v", err)
			return
		}
	}

	if exists {
		err = toml.Unmarshal(contents, &backups)
		if err != nil {
			log.Errorf("failed to decode backup file: %v", err)
			return
		}

		if len(backups.Backups) >= 3 {
			// remove older backups if there are more than 3
			log.Info("truncate backups")
			backups.Backups = []Backup{
				backups.Backups[len(backups.Backups)-2],
				backups.Backups[len(backups.Backups)-1],
			}
		}

		backups.Backups = append(backups.Backups, b)
	} else {
		backups = Backups{
			Backups: []Backup{b},
		}
	}

	marshalled, err := toml.Marshal(backups)
	if err != nil {
		log.Errorf("failed to encode backup file: %v", err)
		return
	}

	err = filesystem.Api().WriteFile(fileName, marshalled, os.ModePerm)
	if err != nil {
		log.Errorf("failed to write backup file: %v", err)
	}
}
