package utils

import (
	"encoding/json"
	"fmt"
)

// HandleResponse formats and prints the API response based on the output type.
func HandleResponse(response []byte, outputType string) error {
	switch outputType {
	case "json":
		// Pretty-print JSON response, machine-readable format
		var formattedResponse map[string]interface{}
		if err := json.Unmarshal(response, &formattedResponse); err != nil {
			return fmt.Errorf("failed to parse JSON response: %w", err)
		}
		prettyJSON, err := json.MarshalIndent(formattedResponse, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format JSON response: %w", err)
		}
		fmt.Println(string(prettyJSON))
	case "text":
		// Human-readable format, assuming the response is a json object
		var formattedResponse map[string]interface{}
		if err := json.Unmarshal(response, &formattedResponse); err != nil {
			return fmt.Errorf("failed to parse JSON response: %w", err)
		}
		for key, value := range formattedResponse {
			fmt.Printf("%s: %v\n", key, value)
		}
	default:
		return fmt.Errorf("invalid output type: %s", outputType)
	}
	return nil
}
