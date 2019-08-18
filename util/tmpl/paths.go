package tmpl

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// GetLocalPath returns the path of the local templates
func GetLocalPath() string {
	var dir string

	if viper.IsSet("templates.paths.local") {
		dir, _ = homedir.Expand(viper.GetString("templates.paths.local"))

		return dir
	}

	dir, _ = homedir.Expand("~/.lavra/templates")

	return dir
}

// GetCachePath returns the path of the local cache of fetched
// remote templates.
func GetCachePath() string {
	var dir string

	if viper.IsSet("templates.paths.cache") {
		dir, _ = homedir.Expand(viper.GetString("templates.paths.cache"))

		return dir
	}

	dir, _ = homedir.Expand("~/.lavra/.cache/templates")

	return dir
}
