package where

import (
	"fmt"
	"github.com/metafates/pat/constant"
	"github.com/metafates/pat/util"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

// Config path
// Will create the directory if it doesn't exist
func Config() string {
	var path string

	if customDir, present := os.LookupEnv(EnvConfigPath); present {
		path = customDir
	} else {
		path = filepath.Join(lo.Must(os.UserConfigDir()), constant.App)
	}

	return mkdir(path)
}

// Logs path
// Will create the directory if it doesn't exist
func Logs() string {
	return mkdir(filepath.Join(Config(), "logs"))
}

// Backup path
// Will create the directory if it doesn't exist
func Backup() string {
	return mkdir(filepath.Join(Config(), "backup"))
}

// Cache path
// Will create the directory if it doesn't exist
func Cache() string {
	genericCacheDir, err := os.UserCacheDir()
	if err != nil {
		genericCacheDir = "."
	}

	cacheDir := filepath.Join(genericCacheDir, constant.PrefixCache)
	return mkdir(cacheDir)
}

// Temp path
// Will create the directory if it doesn't exist
func Temp() string {
	tempDir := filepath.Join(os.TempDir(), constant.PrefixTemp)
	return mkdir(tempDir)
}

func FishScript() string {
	dir := mkdir(util.ResolveTilde(viper.GetString(constant.FishScriptPath)))
	return filepath.Join(dir, fmt.Sprintf(".%s.%s", constant.App, constant.Fish))
}

func ZshScript() string {
	dir := mkdir(util.ResolveTilde(viper.GetString(constant.ZshScriptPath)))
	return filepath.Join(dir, fmt.Sprintf(".%s.%s", constant.App, constant.Zsh))
}

func BashScript() string {
	dir := mkdir(util.ResolveTilde(viper.GetString(constant.BashScriptPath)))
	return filepath.Join(dir, fmt.Sprintf(".%s.%s", constant.App, constant.Bash))
}
