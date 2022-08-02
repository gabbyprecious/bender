package main

import (
	core "github.com/cln-reckless/cln4go.plugin/pkg/plugin"

	"github.com/vincenzopalazzo/cln4go/plugin"
)

func main() {
	plugin := plugin.New(&core.PluginState{}, true, nil)
	plugin.RegisterOption("foo", "string", "Hello Go", "An example of option", false)
	plugin.RegisterRPCMethod("hello", "", "an example of rpc method", &core.Hello[core.PluginState]{})
	plugin.RegisterNotification("shutdown", &core.OnShutdown[core.PluginState]{})
	plugin.Start()
}
