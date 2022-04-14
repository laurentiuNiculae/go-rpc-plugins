package main

import (
	"os"

	"github.com/laurentiuNiculae/go-rpc-plugins/cmd"
	"github.com/laurentiuNiculae/go-rpc-plugins/plugins"
)

func main() {
	plugins.PluginManager.LoadAll("plugin-configs")

	if err := cmd.NewCliRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
