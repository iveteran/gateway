package conf

import (
	"log"

	"matrix.works/fmx-common/conf"
)

var Cfg *MyConfig

type UrlPermission struct {
	UrlWhiteList       []string
	UrlPrefixWhiteList []string
	UrlUserAccessList  []string
	BlackList          []string
}

type Misc struct {
	GuestUserId uint32
}

type MyConfig struct {
	conf.FmxConfig
	UrlPermission UrlPermission
	RouteTable    map[string]string
	Misc          Misc
}

func CreateGlobalConfig(filename string, logger *log.Logger) *MyConfig {
	Cfg = &MyConfig{
		FmxConfig: conf.FmxConfig{
			FilePath: filename,
		},
	}
	conf.Cfg = &Cfg.FmxConfig
	conf.LoadConfig(filename, Cfg, logger)
	conf.CheckRequiredOptions()
	return Cfg
}
