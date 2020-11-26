package main

import (
	"fmt"

	appAuth "matrix.works/fmx-common/web/middleware/authenticate"
	"matrix.works/fmx-gateway/bootstrap"
	"matrix.works/fmx-gateway/conf"
	//userAuth "matrix.works/fmx-gateway/web/middleware/authenticate"
	"matrix.works/fmx-gateway/web/middleware/reverseProxy"
	"matrix.works/fmx-gateway/web/routes"
)

const (
	APP_NAME  = "fgw"
	APP_OWNER = "Matrixworks(ShenZhen) Information Technologies Co.,Ltd."
)

var (
	Version string = "unknown"
	BuildNo string = "unknown"
)

var appTokens = map[string]string{"fimatrix": "fimatrix2020"} // TODO: load from configure file

func newApp() *bootstrap.FgwBootstrapper {
	/// 初始化应用

	app := bootstrap.New(
		APP_NAME,
		APP_OWNER,
		Version,
		BuildNo,
		appTokens,
	)

	app.ParseCommandLine()
	conf.CreateGlobalConfig(app.CmdOpts.ConfigPath, nil)

	app.Bootstrap()

	app.Configure(
		reverseProxy.Configure,
		appAuth.Configure,
		//userAuth.Configure,
		routes.Configure,
	)

	return app
}

func main() {
	app := newApp()

	address := fmt.Sprintf("%s:%d", conf.Cfg.Server.ListenAddress, conf.Cfg.Server.ListenPort)
	app.Listen(address)
}
