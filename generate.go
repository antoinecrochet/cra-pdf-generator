package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/spf13/cobra"
)

var currentDirectory string

var fileNameOutput string
var templateId int

var generatePdfCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate CRA in pdf format",
	Run: func(cmd *cobra.Command, args []string) {
		err := GeneratePdf(fileNameOutput, templateId)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}

		fmt.Printf("Successfully generated %s", fileNameOutput)
	},
}

func init() {
	// Get current directory to save generated pdf(s)
	ex, _ := os.Executable()
	currentDirectory = filepath.Dir(ex)

	// fpdf settings
	fpdf.SetDefaultCompression(false)
	fpdf.SetDefaultCatalogSort(true)
	fpdf.SetDefaultCreationDate(time.Now())
	fpdf.SetDefaultModificationDate(time.Now())

	// Initialize command flags
	generatePdfCmd.Flags().StringVarP(&fileNameOutput, "output", "o", "output.pdf", "Output file name")
	generatePdfCmd.Flags().IntVarP(&templateId, "template", "t", 1, "Pdf template identifier (only available value: 1)")

	// Add generate command to root cmd
	rootCmd.AddCommand(generatePdfCmd)
}
