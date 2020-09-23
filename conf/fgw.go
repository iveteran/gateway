package conf

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
)

var (
	Cfg     *FimatrixConfig
	once    sync.Once
	cfgLock sync.RWMutex
)

var (
	SysTimeLocation  *time.Location
	SysTimeForm      string
	SysTimeFormShort string
)

type ServerConfig struct {
	Host               string
	ListenAddress      string
	ListenPort         int
	LogFile            string
	SysTimeForm        string
	SysTimeFormShort   string
	SysTimeLocation    string
	CookieSecret       string
	SignSecret         string
	AppId              string
	AppToken           string
	UrlWhiteList       []string
	UrlPrefixWhiteList []string
	BlackList          []string
}

type RedisConfig struct {
	Host      string
	Port      int
	Database  int
	IsRunning bool
}

type FimatrixConfig struct {
	FilePath   string
	Server     ServerConfig
	Redises    map[string]RedisConfig
	RouteTable map[string]string
}

func (this *FimatrixConfig) Load(filename string, logger *log.Logger) {
	this.FilePath = filename

	cfgLock.RLock()
	defer cfgLock.RUnlock()

	this.ReloadConfig(logger)
}

func (this *FimatrixConfig) ReloadConfig(logger *log.Logger) {
	filePath, err := filepath.Abs(this.FilePath)
	if err != nil {
		panic(err)
	}
	if _, err := toml.DecodeFile(filePath, this); err != nil {
		panic(err)
	}
	if logger != nil {
		logger.Printf("Parse toml file once. filePath: %s\n", filePath)
		logger.Printf("Config: %+v\n", this)
	} else {
		fmt.Printf("Parse toml file once. filePath: %s\n", filePath)
		fmt.Printf("Config: %+v\n", this)
	}

	SysTimeLocation, _ = time.LoadLocation(this.Server.SysTimeLocation)
	SysTimeForm = this.Server.SysTimeForm
	SysTimeFormShort = this.Server.SysTimeFormShort

	Cfg = this
}

func CreateGlobalConfig(filename string, logger *log.Logger) *FimatrixConfig {
	if filename == "" {
		fmt.Printf("Please give a configure file\n")
		return nil
	}
	Cfg = new(FimatrixConfig)
	Cfg.Load(filename, logger)
	return Cfg
}
