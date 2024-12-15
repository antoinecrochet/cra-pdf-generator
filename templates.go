package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/go-pdf/fpdf"
)

var availableTemplateIds = []int{1}

func GeneratePdf(fileName string, templateId int) error {
	if !slices.Contains(availableTemplateIds, templateId) {
		return fmt.Errorf("template %d is not a valid template (available templates are: %v)", templateId, availableTemplateIds)
	}
	pdf := fpdf.New("P", "mm", "A4", "./assets/fonts")
	pdf.AddUTF8Font("calibri", "", "calibri-font-family/calibri-regular.ttf")
	pdf.AddUTF8Font("calibri", "I", "calibri-font-family/calibri-italic.ttf")
	pdf.AddUTF8Font("calibri", "B", "calibri-font-family/calibri-bold.ttf")
	pdf.AddUTF8Font("calibri", "BI", "calibri-font-family/calibri-bold-italic.ttf")

	if templateId == 1 {
		buildTemplate1(pdf)
	}

	fileStr := filepath.Join(currentDirectory, fileName)

	err := pdf.OutputFileAndClose(fileStr)
	if err != nil {
		return err
	}
	return nil
}

func buildTemplate1(pdf *fpdf.Fpdf) {
	pdf.AddPage()

	lineHeight := 6.
	marginX := 25.
	left := marginX
	pageWidth, _ := pdf.GetPageSize()
	pdfAreaWidth := pageWidth - 2*marginX
	pdf.SetFont("calibri", "", 12)

	// Sender
	pdf.SetXY(left, 25)
	pdf.MultiCell(pdfAreaWidth/2, lineHeight, os.Getenv("TEMPLATE1_SENDER"), "", "L", false)

	// Receiver
	receiverCellWidth := pdfAreaWidth / 2
	pdf.SetX(pageWidth - marginX - receiverCellWidth)
	pdf.MultiCell(receiverCellWidth, lineHeight, os.Getenv("TEMPLATE1_RECEIVER"), "", "R", false)

	pdf.Ln(-1)
	// Subject
	pdf.SetX(left)
	pdf.SetFontStyle("B")
	objectText := "Objet : "
	pdf.Cell(pdf.GetStringWidth(objectText), lineHeight, objectText)
	pdf.SetFontStyle("")
	pdf.Cell(pdfAreaWidth/2, lineHeight, os.Getenv("TEMPLATE1_SUBJECT_PREFIX")+" Novembre 2024")
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
	pdf.SetX(left)
	pdf.SetFontStyle("B")
	pdf.CellFormat(columnWidth[0]+columnWidth[1], tableLineHeight, os.Getenv("TEMPLATE1_TOTAL_TITLE"), "1", 0, "", false, 0, "")
	pdf.CellFormat(columnWidth[2]+columnWidth[3], tableLineHeight, "10", "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.Ln(-1)

	pdf.SetX(left)
	pdf.SetFontStyle("")
	pdf.CellFormat(pdfAreaWidth/2, lineHeight, os.Getenv("TEMPLATE1_SENDER_SIGNATURE_TITLE"), "", 0, "L", false, 0, "")
	pdf.SetX(pageWidth - marginX - pdfAreaWidth/2)
	pdf.CellFormat(pdfAreaWidth/2, lineHeight, os.Getenv("TEMPLATE1_RECEIVER_SIGNATURE_TITLE"), "", 0, "R", false, 0, "")
}
