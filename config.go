package fxgo

import (
	"flag"

	"github.com/fidelfly/fxgo/confx"
)

const defConfigFile = "config.toml"

//export
func InitTomlConfig(filepath string, Properties interface{}) (err error) {
	var configFile = filepath
	flag.StringVar(&configFile, "config", defConfigFile, "Set Config File")
	flag.Parse()

	// Parse Config File
	err = confx.ParseToml(configFile, Properties)
	return
}
