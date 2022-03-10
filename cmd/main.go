package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tommartensen/tfui/pkg/client/plan"
	"github.com/tommartensen/tfui/pkg/client/system"
	"github.com/tommartensen/tfui/pkg/config"
	"github.com/tommartensen/tfui/pkg/server"
)

var rootCmd = &cobra.Command{
	Use:     "tfui",
	Short:   "The Terraform UI",
	Long:    "tfui is a tool to manage the Terraform UI server, e.g. upload plans, or reset the server.",
	Version: "0.0.1",
}

func currentConfig(cmd *cobra.Command, args []string) {
	configuration := config.New()
	fmt.Printf("Server Address: %s\n", configuration.Addr)
	fmt.Printf("Client Token:   %s\n", configuration.ClientToken)
}

func addConfigCommands(rootCmd *cobra.Command) {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manages local config",
	}
	configShowCmd := &cobra.Command{
		Use:   "show",
		Short: "Shows local config",
		Run:   currentConfig,
	}
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
}

func addSystemCommands(rootCmd *cobra.Command) {
	systemCmd := &cobra.Command{
		Use:   "system",
		Short: "Commands to manage the system",
	}
	resetSystemCmd := &cobra.Command{
		Use:   "reset",
		Short: "Resets system",
		Run:   system.Reset,
	}
	statusSystemCmd := &cobra.Command{
		Use:   "status",
		Short: "Show system status",
		Run:   system.Status,
	}
	systemCmd.AddCommand(resetSystemCmd)
	systemCmd.AddCommand(statusSystemCmd)
	rootCmd.AddCommand(systemCmd)
}

func addPlanCommands(rootCmd *cobra.Command) {
	planCmd := &cobra.Command{
		Use:   "plan",
		Short: "Commands to manage the Terraform plans",
	}
	uploadPlanCmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload a plan to the server",
		Run:   plan.UploadPlan,
	}
	uploadPlanCmd.Flags().StringP("file", "f", "", "path to the plan in JSON format to upload")
	uploadPlanCmd.MarkFlagRequired("file")
	planCmd.AddCommand(uploadPlanCmd)
	rootCmd.AddCommand(planCmd)
}

func addServerCommands(rootCmd *cobra.Command) {
	serverCommand := &cobra.Command{
		Use:   "server",
		Short: "Start the server",
		Run:   server.RunServer,
	}
	rootCmd.AddCommand(serverCommand)
}

func init() {
	addConfigCommands(rootCmd)
	addSystemCommands(rootCmd)
	addPlanCommands(rootCmd)
	addServerCommands(rootCmd)
}

func main() {
	rootCmd.Execute()
}
