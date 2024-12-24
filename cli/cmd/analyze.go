/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	utils "vm2cont/cli/pkg/utils"

	"github.com/spf13/cobra"
)

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Perform analysis on the target VM where the application is running",
	RunE: func(cmd *cobra.Command, args []string) error {
		analysisType, _ := cmd.Flags().GetString("type")
		user, _ := cmd.Flags().GetString("user")
		host, _ := cmd.Flags().GetString("host")
		privateKeyPath, _ := cmd.Flags().GetString("privateKeyPath")

		// Check if the required flags are provided
		if user == "" || host == "" || privateKeyPath == "" {
			return fmt.Errorf("user, host, and privateKeyPath are required flags")
		}

		payload := map[string]interface{}{
			"analyzerApproach": analysisType,
			"user":             user,
			"host":             host,
			"privateKeyPath":   privateKeyPath,
		}

		// Get the output type from the --output flag
		outputType, _ := cmd.Flags().GetString("output")

		var response []byte
		var err error

		// If analysisType is "mixed", request more input from the user
		if analysisType == "mixed" {
			fmt.Println("Mixed analysis selected. Please provide the following details:")
			fmt.Print("Enter the strategy for collecting application files: ")
			var appFilesStrategy string
			fmt.Scanln(&appFilesStrategy)

			fmt.Print("Enter the strategy for collecting exposed ports: ")
			var exposedPortsStrategy string
			fmt.Scanln(&exposedPortsStrategy)

			fmt.Print("Enter the strategy for collecting services: ")
			var servicesStrategy string
			fmt.Scanln(&servicesStrategy)

			payload["appFilesStrategy"] = appFilesStrategy
			payload["exposedPortsStrategy"] = exposedPortsStrategy
			payload["servicesStrategy"] = servicesStrategy

			response, err = utils.MakeRequest("POST", "http://localhost:8001/analyze/complete/mixed-approach", payload)
			if err != nil {
				return err
			}
		} else {
			response, err = utils.MakeRequest("POST", "http://localhost:8001/analyze/complete/single-approach", payload)
			if err != nil {
				return err
			}
		}

		// Use HandleResponse for output
		if err := utils.HandleResponse(response, outputType); err != nil {
			return fmt.Errorf("failed to handle response: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)

	// Add a flag to add type of analysis - file system, process, or mixed
	analyzeCmd.Flags().StringP("type", "t", "fs", "Type of analysis to perform. Options are: fs, process, mixed")
	analyzeCmd.Flags().StringP("user", "", "", "Username for SSH connection")
	analyzeCmd.Flags().StringP("host", "", "", "Host for SSH connection")
	analyzeCmd.Flags().StringP("privateKeyPath", "", "", "Path to private key for SSH connection")
}
