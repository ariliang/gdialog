package global

import (
	"github.com/BurntSushi/toml"
)

var Config tomlConfig

type (
	tomlConfig struct {
		Server serverToml `toml:"server"`
		Web    webToml    `toml:"web"`
		WX     wxToml     `toml:"wx_config"`
	}

	serverToml struct {
		Host string
		Port int32
	}

	webToml struct {
		SessionSecret string `toml:"session_secret"`
	}

	wxToml struct {
		AppId       string
		AppIdSecret string `toml:"appid_secret"`
	}
)

func (c *tomlConfig) GetConfig(path string) {
	toml.DecodeFile(path, &c)
}

func init() {
	Config.GetConfig("conf/config.toml")
}
