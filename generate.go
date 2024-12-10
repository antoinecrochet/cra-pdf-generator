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

var generatePdfCmd = &cobra.Command{
	Use:   "generate [CSV file]",
	Short: "Generate CRA in pdf format",
	Run: func(cmd *cobra.Command, args []string) {
		generatePdf("cra.pdf")
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

	// Add generate command to root cmd
	rootCmd.AddCommand(generatePdfCmd)
}

func generatePdf(fileName string) {
	pdf := fpdf.New("P", "mm", "A4", "./assets/fonts")
	pdf.AddUTF8Font("calibri", "", "calibri-font-family/calibri-regular.ttf")
	pdf.AddUTF8Font("calibri", "I", "calibri-font-family/calibri-italic.ttf")
	pdf.AddUTF8Font("calibri", "B", "calibri-font-family/calibri-bold.ttf")
	pdf.AddUTF8Font("calibri", "BI", "calibri-font-family/calibri-bold-italic.ttf")

	pdf.AddPage()

	lineHeight := 6.
	marginX := 25.
	left := marginX
	pageWidth, _ := pdf.GetPageSize()
	pdfAreaWidth := pageWidth - 2*marginX
	pdf.SetFont("calibri", "", 12)

	// Sender
	pdf.SetXY(left, 25)
	pdf.MultiCell(pdfAreaWidth/2, lineHeight, "Lorem ipsum dolor\nsit amet, consectetur\nadipiscing elit.", "", "L", false)

	// Receiver
	receiverCellWidth := pdfAreaWidth / 2
	pdf.SetX(pageWidth - marginX - receiverCellWidth)
	pdf.MultiCell(receiverCellWidth, lineHeight, "Lorem ipsum dolor\nsit amet, consectetur\nadipiscing elit.", "", "R", false)

	pdf.Ln(-1)
	// Subject
	pdf.SetX(left)
	pdf.SetFontStyle("B")
	objectText := "Objet : "
	pdf.Cell(pdf.GetStringWidth(objectText), lineHeight, objectText)
	pdf.SetFontStyle("")
	pdf.Cell(pdfAreaWidth/2, lineHeight, "Etiam sit amet arcu sodales, iaculis velit sed, imperdiet magna")
	pdf.Ln(-1)
	pdf.Ln(-1)

	// Table
	pdf.SetFontStyle("B")
	pdf.SetX(left)
	tableLineHeight := 9.
	columnWidth := []float64{40, 35, 40, 35}
	header := []string{"Jour", "Présent", "Jour", "Présent"}
	for i, str := range header {
		pdf.CellFormat(columnWidth[i], tableLineHeight, str, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)
	// Color and font restoration
	pdf.SetFillColor(239, 239, 239)
	pdf.SetFontStyle("")
	for i := 0; i < 16; i++ {
		pdf.SetX(left)
		pdf.CellFormat(columnWidth[0], tableLineHeight, "", "1", 0, "", i%2 == 0, 0, "")
		pdf.CellFormat(columnWidth[1], tableLineHeight, "", "1", 0, "C", i%2 == 0, 0, "")
		pdf.CellFormat(columnWidth[2], tableLineHeight, "", "1", 0, "", i%2 == 0, 0, "")
		pdf.CellFormat(columnWidth[3], tableLineHeight, "", "1", 0, "C", i%2 == 0, 0, "")
		pdf.Ln(-1)
	}
	fileStr := filepath.Join(currentDirectory, fileName)

	err := pdf.OutputFileAndClose(fileStr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("Successfully generated %s", fileName)
}