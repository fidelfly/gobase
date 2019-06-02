package fxgo

import (
	"flag"

	"github.com/fidelfly/fxgo/confx"
)

const defConfigFile = "confx.toml"

//export
func InitTomlConfig(filepath string, Properties interface{}) (err error) {
	var configFile = filepath
	flag.StringVar(&configFile, "confx", defConfigFile, "Set Config File")
	flag.Parse()

	// Parse Config File
	err = confx.ParseToml(configFile, Properties)
	return
}
