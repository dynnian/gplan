package main

import (
	"log"
	"os"

	"codeberg.org/dynnian/gplan/controllers"
	"codeberg.org/dynnian/gplan/repository"
	"codeberg.org/dynnian/gplan/views/cmd"
)

func run() int {
	// Initialize the repository
	repo, repoErr := repository.NewRepository()
	if repoErr != nil {
		log.Printf("Error initializing repository: %v", repoErr)
		return 1 // Exit code 1 indicates failure
	}
	defer repo.Close()

	// Initialize controllers
	projectController := controllers.NewProjectController(repo)
	taskController := controllers.NewTaskController(repo)

	// Initialize the root command with controllers
	rootCmd := cmd.NewRootCmd(projectController, taskController)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		return 1 // Return 1 on error
	}

	return 0 // Exit code 0 indicates success
}

func main() {
	// Call the run function and exit with the appropriate code
	os.Exit(run())
}
