package cmd

import (
	"context"
	"fmt"

	"github.com/laurentiuNiculae/go-rpc-plugins/plugins"
	cli "github.com/laurentiuNiculae/go-rpc-plugins/plugins/clicommand"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// interface
type CLICommand interface {
	Command() *cobra.Command
}

var cliCommandImplementations map[string]plugins.Plugin = map[string]plugins.Plugin{}

func AllPlugins() map[string]plugins.Plugin {
	return cliCommandImplementations
}

// struct that implements the interface and calls on the client for operations
type CLICommandImpl struct {
	client cli.CLICommandClient
}

type CLIBuilder struct{}

func (clib CLIBuilder) Build(addr, port string) plugins.Plugin {
	address := fmt.Sprintf("%s:%s", addr, port)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("Can't connect")
	}

	c := cli.NewCLICommandClient(conn)

	return CLICommandImpl{client: c}
}

func (cci CLICommandImpl) Command() *cobra.Command {
	return &cobra.Command{
		Use:     "test",
		Aliases: []string{"test"},
		Short:   "This command is for test",
		Long:    "This command is for test",
		Run: func(cmd *cobra.Command, args []string) {
			response, err := cci.client.Command(
				context.Background(),
				&cli.CLIArgs{
					Args: args,
				})

			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Println(response.GetMessage())
		},
	}
}

func init() {
	plugins.PluginManager.RegisterInterface("CLICommand", cliCommandImplementations, CLIBuilder{})
}
