package main

import (
	"fmt"
	"github.com/long250038728/web/cmd/gen/client"
	"github.com/spf13/cobra"
)

func main() {
	devopsCron := client.NewDevopsCorn()
	serverCron := client.NewServerCornCorn()

	rootCmd := &cobra.Command{
		Use:   "serverGen",
		Short: "快速生成工具",
	}

	rootCmd.AddCommand(serverCron.Server())
	rootCmd.AddCommand(devopsCron.Devops())
	rootCmd.AddCommand()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
	}
}
