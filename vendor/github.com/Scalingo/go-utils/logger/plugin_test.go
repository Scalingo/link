package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type TestPlugin struct {
	Accept     bool
	PluginHook logrus.Hook
	PluginName string
}

func (p TestPlugin) Hook() (bool, logrus.Hook) { return p.Accept, p.PluginHook }
func (p TestPlugin) Name() string              { return p.PluginName }

func TestRegisterPlugin(t *testing.T) {
	examples := []struct {
		Name          string
		BasePlugins   []Plugin
		Plugin        Plugin
		ExpectPlugins []string
	}{
		{
			Name:          "when the plugin is not in memory it should add it",
			Plugin:        TestPlugin{PluginName: "test"},
			ExpectPlugins: []string{"test"},
		}, {
			Name:          "when the plugin is already in memory it should not add it",
			BasePlugins:   []Plugin{TestPlugin{PluginName: "test"}},
			Plugin:        TestPlugin{PluginName: "test"},
			ExpectPlugins: []string{"test"},
		}, {
			Name:          "whan another plugin in in memory, it should add it",
			BasePlugins:   []Plugin{TestPlugin{PluginName: "test1"}},
			Plugin:        TestPlugin{PluginName: "test"},
			ExpectPlugins: []string{"test1", "test"},
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			manager := PluginManager{plugins: example.BasePlugins}
			manager.RegisterPlugin(example.Plugin)

			for i, p := range manager.plugins {
				assert.Equal(t, example.ExpectPlugins[i], p.Name())
			}
		})
	}
}

func TestPluginHooks(t *testing.T) {
	examples := []struct {
		Name        string
		BasePlugins []Plugin
		HookLength  int
	}{
		{
			Name:       "when there is no plugins",
			HookLength: 0,
		}, {
			Name:        "when there is a plugin that do return false",
			BasePlugins: []Plugin{TestPlugin{PluginName: "test", Accept: false}},
			HookLength:  0,
		}, {
			Name:        "when there is a plugin that return true",
			BasePlugins: []Plugin{TestPlugin{PluginName: "test", Accept: true}},
			HookLength:  1,
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			manager := PluginManager{plugins: example.BasePlugins}
			assert.Equal(t, example.HookLength, len(manager.Hooks()))
		})
	}
}
