package main

import (
	"github.com/LarsNieuwenhuizen/admission-test/webhook"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "help",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func main() {
	rootCmd.AddCommand(webhook.CmdWebhook)
	rootCmd.Execute()
}
