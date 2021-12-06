package main

import (
	"fmt"
	"log"
	"os"

	"matrix.works/fmx-gateway/bootstrap"
	"matrix.works/fmx-gateway/conf"
	"matrix.works/fmx-gateway/web2/route"
)

const (
	AppName  = "fgw"
	AppOwner = "Matrixworks(ShenZhen) Information Technologies Co.,Ltd."
)

var (
	Version string = "unknown"
	BuildNo string = "unknown"
)

var appTokens = map[string]string{"fimatrix": "fimatrix2020"} // TODO: load from configure file

func newApp() *bootstrap.FgwBootstrapper {
	app := bootstrap.New(
		AppName,
		AppOwner,
		Version,
		BuildNo,
		appTokens,
	)

	app.ParseCommandLine()
	conf.CreateGlobalConfig(app.CmdOpts.ConfigPath, nil)

	app.Bootstrap()

	routeMap := conf.Cfg.RouteTable
	fmt.Printf("route table: %+v\n", routeMap)
	route.Setup(routeMap)

	return app
}

func main() {
	app := newApp()

	addr := fmt.Sprintf("%s:%d", conf.Cfg.Server.ListenAddress, conf.Cfg.Server.ListenPort)
	err := app.Serve(addr)
	if err != nil {
		log.Printf("Start server failed: %s\n", err.Error())
		os.Exit(1)
	}
}
