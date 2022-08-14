package plugin

import (
	"fmt"

	"github.com/vincenzopalazzo/cln4go/plugin"
)

type LNDashboardRune[T PluginState] struct{}

func (instance *LNDashboardRune[T]) Call(plugin *plugin.Plugin[PluginState], request map[string]any) (map[string]any, error) {
	runeClnapp, err := plugin.State.Client.Call("commando-rune", map[string]any{
		"restrictions": "[\"method^list|method^get|method=decode|method=fetchinvoice|method=ping\",\"method/listdatastore\"]",
	})
	if err != nil {
		plugin.Log("broken", fmt.Sprintf("%s", err))
		return nil, err
	}
	runeVal, _ := runeClnapp["rune"].(string)
	return map[string]interface{}{"lndashboard-rune": runeVal}, nil
}

type ClnAppRune[T PluginState] struct{}

func (instance *ClnAppRune[T]) Call(plugin *plugin.Plugin[PluginState], request map[string]any) (map[string]any, error) {
	runeClnapp, err := plugin.State.Client.Call("commando-rune", map[string]any{
		"restrictions": "[\"method^list|method^get|method=decode|method=fetchinvoice|method=ping\",\"method/listdatastore\"]",
	})
	if err != nil {
		plugin.Log("broken", fmt.Sprintf("%s", err))
		return nil, err
	}
	runeVal := runeClnapp["rune"].(string)
	return map[string]interface{}{"clnapp-rune": runeVal}, nil
}
