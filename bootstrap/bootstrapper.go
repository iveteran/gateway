package bootstrap

import (
	"log"
	"net/http"

	"matrix.works/fmx-common/web/bootstrap"
)

type Bootstrapper = bootstrap.Bootstrapper
type CommandOptions = bootstrap.CommandOptions

type FgwCommandOptions struct {
	*CommandOptions
	// add more options here
}

type FgwBootstrapper struct {
	*Bootstrapper
	CmdOpts *FgwCommandOptions
}

func New(
	appName, appOwner, appVersion, appBuildNo string,
	tokenTable map[string]string, cfgList ...bootstrap.Configurator,
) *FgwBootstrapper {

	b := &FgwBootstrapper{
		Bootstrapper: bootstrap.New(
			appName,
			appOwner,
			appVersion,
			appBuildNo,
			tokenTable,
			cfgList...,
		),
		CmdOpts: &FgwCommandOptions{
			&CommandOptions{},
		},
	}

	return b
}

func (this *FgwBootstrapper) ParseCommandLine() {
	this.Bootstrapper.ParseCommandLine(this.CmdOpts.CommandOptions)
}

func (this *FgwBootstrapper) Serve(addr string) error {
	server := &http.Server{
		Addr: addr,
	}

	log.Printf("Listening on %s...", addr)
	return server.ListenAndServe()
}
