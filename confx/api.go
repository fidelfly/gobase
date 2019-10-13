package confx

import (
	"os"

	"github.com/fidelfly/gox/logx"

	"github.com/BurntSushi/toml"
)

func ParseToml(file string, target interface{}) (err error) {
	if _, err = os.Stat(file); err == nil {
		logx.Infof("Config file : %s Found!", file)
		_, err = toml.DecodeFile(file, target)
		if err != nil {
			logx.Error(err)
		}
	} else {
		logx.Errorf("Config file : %s is not found!", file)
	}
	return
}
