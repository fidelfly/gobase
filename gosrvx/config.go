package gosrvx

import (
	"flag"

	"github.com/fidelfly/gox/confx"
)

const defConfigFile = "config.toml"

//export
func InitTomlConfig(filepath string, properties interface{}) (err error) {
	var configFile = filepath
	flag.StringVar(&configFile, "config", defConfigFile, "Set Config File")
	flag.Parse()

	// Parse Config File
	err = confx.ParseToml(configFile, properties)
	return
}
