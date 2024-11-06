/*
Copyright © 2024 Antonio Takev tonitakev@gmail.com
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// dockerizeCmd represents the dockerize command
var dockerizeCmd = &cobra.Command{
	Use:   "dockerize",
	Short: "Create a Dockerfile for the application",
	Long: `Based on the provided application profile, 
create a Dockerfile that will be used to build the container image.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dockerize called")
	},
}

func init() {
	rootCmd.AddCommand(dockerizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dockerizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dockerizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
