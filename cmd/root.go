/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// "zli" - client-side cli.
func NewCliRootCmd() *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "go-rpc-plugins",
		Short: "",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// additional cmds
	enableCliPlugins(rootCmd)

	return rootCmd
}

func enableCliPlugins(cmd *cobra.Command) {
	// init clients for each config
	for k, p := range AllPlugins() {
		fmt.Println("Adding plugin: ", k)
		cmd.AddCommand(p.(CLICommand).Command())
	}
}
