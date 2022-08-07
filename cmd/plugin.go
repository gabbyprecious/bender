package main

import (
	core "github.com/vincenzopalazzo/bender/pkg/plugin"

	"github.com/vincenzopalazzo/cln4go/plugin"
)

func main() {
	plugin := plugin.New(&core.PluginState{}, true, nil)
	plugin.RegisterOption("foo", "string", "Hello Go", "An example of option", false)
	plugin.RegisterRPCMethod("bender_run_server", "", "Run provider server to expose the endpoint", &core.StartServer[core.PluginState]{})
	plugin.RegisterRPCMethod("bender_set_password", "", "set password for tls certificate retrival", &core.SetPassword[core.PluginState]{})
	plugin.RegisterNotification("shutdown", &core.OnShutdown[core.PluginState]{})
	plugin.Start()
}
