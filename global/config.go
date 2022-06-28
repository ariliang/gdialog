package global

import (
	"github.com/BurntSushi/toml"
)

var Config tomlConfig

type (
	tomlConfig struct {
		Server     serverToml     `toml:"server"`
		Web        webToml        `toml:"web"`
		DialogCore dialogCoreToml `toml:"dialog_core"`
		WX         wxToml         `toml:"wx_config"`
	}

	serverToml struct {
		Host string
		Port int32
	}

	webToml struct {
		SessionSecret string `toml:"session_secret"`
	}

	dialogCoreToml struct {
		Host        string
		DialogModel []string `toml:"dialog_model"`
	}

	wxToml struct {
		AppId       string
		AppIdSecret string `toml:"appid_secret"`
		AuthSimu    bool   `toml:"auth_simu"`
	}
)

func (c *tomlConfig) GetConfig(path string) {
	toml.DecodeFile(path, &c)
}

func init() {
	Config.GetConfig("conf/config.toml")
}
