/*
Copyright © 2024 Antonio Takev tonitakev.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cli",
	Short: "Converts an application from VM to a container",
	Long: `vm2cont is a tool that migrate applications residing in VMs 
to container so that they can be run and deployed on different 
platforms and be used in Kubernetes clusters.`,
	Version: "1.0.0",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Ensure the output flag is valid
		output, _ := cmd.Flags().GetString("output")
		if output != "text" && output != "json" {
			return fmt.Errorf("invalid output format: %s (valid options are 'text' or 'json')", output)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("output", "o", "text", "Output format: text or json")
}
