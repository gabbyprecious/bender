package plugin

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/vincenzopalazzo/cln4go/client"
	"github.com/vincenzopalazzo/cln4go/plugin"
)

type Login struct {
	Password string `json:"password" binding:"required"`
}

type SetPassword[T PluginState] struct{}

func (instance *SetPassword[T]) Call(plugin *plugin.Plugin[PluginState], request map[string]any) (map[string]any, error) {
	plugin.State.Password = fmt.Sprintf(request["password"].(string))
	return map[string]interface{}{"message": "Password set", "password": plugin.State.Password}, nil
}

type OnShutdown[T PluginState] struct{}

func (instance *OnShutdown[T]) Call(plugin *plugin.Plugin[PluginState], request map[string]any) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := plugin.State.Server.Shutdown(ctx); err != nil {
		panic(err)
	}
	os.Exit(0)
}

//nolint:all
func OnInit(plugin *plugin.Plugin[PluginState], config map[string]any) map[string]any {
	lightningDir, found := config["lightning-dir"].(string)
	if !found {
		panic(found)
	}
	rpcFileName, found := config["rpc-file"].(string)
	if !found {
		panic(found)
	}

	unixPath := strings.Join([]string{lightningDir, rpcFileName}, "/")
	rpc, err := client.NewUnix(unixPath)
	if err != nil {
		panic(err)
	}

	plugin.State.Client = rpc
	return map[string]any{}
}
