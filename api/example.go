package main

import "fmt"

const API_NAME string = "vm2cont"

func main() {
	// print the title of my API
	// Info: Use `:=` to infer the type of the variable from the value, the compiler will decide the type
	var version string = "1.0"
	// Array of strings with the name of modules
	var modules = [2]string{"analyze", "dockerize"}
	fmt.Println("## Title of the API")
	// Make sure to add `\n` to avoid appearance of weird symbols at the end of the print statement
	fmt.Printf("%v API, version: %v\n", API_NAME, version)
	// Iterate over the modules array and print the name of each module
	fmt.Println("## Modules of the API")
	for i, module := range modules {
		fmt.Printf("Module %v: %v\n", i, module)
	}
	// Create a slice of strings with the name of the endpoints for module analyze
	var analyzeEndpoints = []string{"analyze", "analyze/summary", "analyze/summary/wordcount"}
	fmt.Println("## Endpoints of the API module analyze")
	// Iterate over the analyzeEndpoints slice and print the name of each endpoint
	for i, endpoint := range analyzeEndpoints {
		fmt.Printf("Endpoint %v: %v\n", i, endpoint)
	}

	// Create a struct to represent a user
	type User struct {
		ID       int
		Username string
		Email    string
	}
	// Create a new user
	user := User{ID: 1, Username: "user1", Email: "user@gmail.com"}
	fmt.Println("## User of the API")
	// Print the user information
	fmt.Printf("User: %v\n", user)

	// Create a map with the names of the modules as keys and the endpoints as values
	var modulesEndpoints = map[string][]string{
		"analyze":   analyzeEndpoints,
		"dockerize": {"dockerize", "dockerize/summary"},
	}
	fmt.Println("## Modules and endpoints of the API")
	// Iterate over the modulesEndpoints map and print the name of each module and its endpoints
	for module, endpoints := range modulesEndpoints {
		fmt.Printf("Module: %v\n", module)
		for i, endpoint := range endpoints {
			fmt.Printf("Endpoint %v: %v\n", i, endpoint)
		}
	}
}
