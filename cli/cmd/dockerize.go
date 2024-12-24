/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"vm2cont/cli/pkg/utils"

	"github.com/spf13/cobra"
)

// dockerizeCmd represents the dockerize command
var dockerizeCmd = &cobra.Command{
	Use:   "dockerize",
	Short: "Perform the dockerization of the application, converting it to a container",
	RunE: func(cmd *cobra.Command, args []string) error {
		dockerImageName, _ := cmd.Flags().GetString("dockerImageName")
		dockerContainerName, _ := cmd.Flags().GetString("dockerContainerName")

		// Check if the required flags are provided
		if dockerImageName == "" || dockerContainerName == "" {
			return fmt.Errorf("dockerImageName and dockerContainerName are required flags")
		}

		payload := map[string]interface{}{
			"dockerImageName":     dockerImageName,
			"dockerContainerName": dockerContainerName,
		}

		// Get the output type from the --output flag
		outputType, _ := cmd.Flags().GetString("output")

		var response []byte
		var err error

		// Perform the dockerization process
		fmt.Println("Performing dockerization...")

		// Make a call to the dockerization API
		response, err = utils.MakeRequest("POST", "http://localhost:8001/dockerize/complete", payload)
		if err != nil {
			return err
		}

		// Use HandleResponse for output
		if err := utils.HandleResponse(response, outputType); err != nil {
			return fmt.Errorf("failed to handle response: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dockerizeCmd)

	// Add flags for the dockerize command
	dockerizeCmd.Flags().StringP("dockerImageName", "", "", "Name of the Docker image to be built")
	dockerizeCmd.Flags().StringP("dockerContainerName", "", "", "Name of the Docker container to be created")
}
