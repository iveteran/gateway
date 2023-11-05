package conf

import (
	"log"
)

var Cfg *FgwConfig

type UrlPermission struct {
	UrlWhiteList       []string
	UrlPrefixWhiteList []string
	UrlUserAccessList  []string
	BlackList          []string
}

type Misc struct {
	GuestUserId uint32
}

type FgwConfig struct {
	Config
	UrlPermission UrlPermission
	RouteTable    map[string]string
	Misc          Misc
}

func CreateGlobalConfig(filename string, logger *log.Logger) *FgwConfig {
	Cfg = &FgwConfig{
		Config: Config{
			FilePath: filename,
		},
	}
	BaseCfg = &Cfg.Config
	LoadConfig(filename, Cfg, logger)
	CheckRequiredOptions()
	return Cfg
}
