package plugin

import (
	"os"

	"github.com/vincenzopalazzo/cln4go/plugin"
)

type PluginState struct{}

type Hello[T PluginState] struct{}

func (instance *Hello[T]) Call(plugin *plugin.Plugin[T], request map[string]any) (map[string]any, error) {
	return map[string]any{"message": "hello from cln4go.template"}, nil
}

type OnShutdown[T PluginState] struct{}

func (instance *OnShutdown[T]) Call(plugin *plugin.Plugin[T], request map[string]any) {
	os.Exit(0)
}
