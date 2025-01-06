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
var selectedYear int
var selectedMonth int
var workedDays string
var notWorkedDays string

var generatePdfCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate CRA in pdf format",
	Run: func(cmd *cobra.Command, args []string) {
		err := GeneratePdf(fileNameOutput, templateId, selectedYear, selectedMonth, workedDays, notWorkedDays)
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
	year, month, _ := time.Now().Date()
	generatePdfCmd.Flags().StringVarP(&fileNameOutput, "output", "o", "output.pdf", "Output file name")
	generatePdfCmd.Flags().IntVarP(&templateId, "template", "t", 1, "Pdf template identifier (only available value: 1)")
	generatePdfCmd.Flags().IntVarP(&selectedYear, "year", "y", year, "Template 1 selected year (default is current year)")
	generatePdfCmd.Flags().IntVarP(&selectedMonth, "month", "m", int(month), "Template 1 selected month (default is current month)")
	generatePdfCmd.Flags().StringVarP(&workedDays, "days", "d", "", "List of worked days separated with comma")
	generatePdfCmd.Flags().StringVarP(&notWorkedDays, "ndays", "n", "", "List of not worked days separated with comma")

	// Add generate command to root cmd
	rootCmd.AddCommand(generatePdfCmd)
}
