package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cra-pdf-generator",
	Short: "CRA in pdf format generator",
}

func main() {
	// Load environment file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Warning: error loading .env file, some templates might not be avalaible.")
		fmt.Println(err)
	}
	// Execute command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
