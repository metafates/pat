package config

import (
	"github.com/metafates/pat/constant"
	"github.com/metafates/pat/util"
	"github.com/spf13/viper"
)

var fields = []Field{
	// LOGS
	{
		constant.LogsWrite,
		false,
		"Write logs to file",
	},
	{
		constant.LogsLevel,
		"info",
		`Logs level.
Available options are: (from less to most verbose)
panic, fatal, error, warn, info, debug, trace`,
	},
	// END LOGS

	// BACKUP
	{
		constant.BackupEnabled,
		true,
		"Enable backups",
	},
	// END BACKUP

	// FISH
	{
		constant.FishScriptPath,
		util.ResolveTilde("~/.pat"),
		`Fish script path
You can use ~ to represent your home directory`,
	},
	// END FISH

	// ZSH
	{
		constant.ZshScriptPath,
		util.ResolveTilde("~/.pat"),
		`Zsh script path
You can use ~ to represent your home directory`,
	},
	// END ZSH

	// BASH
	{
		constant.BashScriptPath,
		util.ResolveTilde("~/.pat"),
		`Bash script path
You can use ~ to represent your home directory`,
	},
	// END BASH
}

func setDefaults() {
	Default = make(map[string]Field)
	for _, f := range fields {
		Default[f.Key] = f
		viper.SetDefault(f.Key, f.Value)
		viper.MustBindEnv(f.Key)
	}
}

var Default map[string]Field
