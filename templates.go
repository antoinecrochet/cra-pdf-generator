package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
)

// Struct representing a table entry
type tableEntry struct {
	day     string
	present bool
	fill    bool
}

// Translated value of time pkg
var translatedMonths = []string{"Janvier", "Février", "Mars", "Avril", "Mai", "Juin", "Juillet", "Août", "Septembre", "Octobre", "Novembre", "Décembre"}
var translatedDays = []string{"Dimanche", "Lundi", "Mardi", "Mercredi", "Jeudi", "Vendredi", "Samedi"}

// Public holidays in France (days per month number)
var publicHolidays = map[int][]int{
	1:  {1},
	4:  {21},
	5:  {1, 8, 29},
	6:  {9},
	7:  {14},
	8:  {15},
	11: {1, 11},
	12: {25},
}

// List of available templates identifiers
var availableTemplateIds = []int{1}

func GeneratePdf(fileName string, templateId int, selectedYear int, selectedMonth int, workedDays string) error {
	// Validate parameters
	if !slices.Contains(availableTemplateIds, templateId) {
		return fmt.Errorf("template %d is not a valid template (available templates are: %v)", templateId, availableTemplateIds)
	}
	// Convert workedDays into int array
	var workedDaysArray []int
	if workedDays != "" {
		workedDaysStringArray := strings.Split(workedDays, ",")
		for _, value := range workedDaysStringArray {
			intWorkedDay, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			workedDaysArray = append(workedDaysArray, intWorkedDay)
		}
	}

	pdf := fpdf.New("P", "mm", "A4", "./assets/fonts")
	pdf.AddUTF8Font("calibri", "", "calibri-font-family/calibri-regular.ttf")
	pdf.AddUTF8Font("calibri", "I", "calibri-font-family/calibri-italic.ttf")
	pdf.AddUTF8Font("calibri", "B", "calibri-font-family/calibri-bold.ttf")
	pdf.AddUTF8Font("calibri", "BI", "calibri-font-family/calibri-bold-italic.ttf")

	if templateId == 1 {
		buildTemplate1(pdf, selectedYear, selectedMonth, workedDaysArray)
	}

	fileStr := filepath.Join(currentDirectory, fileName)

	err := pdf.OutputFileAndClose(fileStr)
	if err != nil {
		return err
	}
	return nil
}

// Build template 1
func buildTemplate1(pdf *fpdf.Fpdf, selectedYear int, selectedMonth int, workedDays []int) {
	// constants
	lineHeight := 6.
	marginX := 25.
	left := marginX
	crossCharacter := "x"
	offsetSecondColumn := 16

	// Add new page
	pdf.AddPage()

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
	pdf.Cell(pdfAreaWidth/2, lineHeight, fmt.Sprintf("%s %s %d", os.Getenv("TEMPLATE1_SUBJECT_PREFIX"), translatedMonths[selectedMonth-int(time.January)], selectedYear))
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

	// Build table
	tableEntries := buildTableContent(selectedYear, selectedMonth, workedDays)
	for i := 0; i < 16; i++ {
		pdf.SetX(left)
		pdf.CellFormat(columnWidth[0], tableLineHeight, tableEntries[i].day, "1", 0, "", tableEntries[i].fill, 0, "")
		presentValue := ""
		if tableEntries[i].present {
			presentValue = crossCharacter
		}
		pdf.CellFormat(columnWidth[1], tableLineHeight, presentValue, "1", 0, "C", tableEntries[i].fill, 0, "")

		if len(tableEntries)-1 < i+offsetSecondColumn {
			pdf.CellFormat(columnWidth[2], tableLineHeight, "", "1", 0, "", false, 0, "")
			pdf.CellFormat(columnWidth[3], tableLineHeight, "", "1", 0, "C", false, 0, "")
			pdf.Ln(-1)
			continue
		}

		pdf.CellFormat(columnWidth[2], tableLineHeight, tableEntries[i+offsetSecondColumn].day, "1", 0, "", tableEntries[i+offsetSecondColumn].fill, 0, "")
		presentValueSecondary := ""
		if tableEntries[i+offsetSecondColumn].present {
			presentValueSecondary = crossCharacter
		}
		pdf.CellFormat(columnWidth[3], tableLineHeight, presentValueSecondary, "1", 0, "C", tableEntries[i+offsetSecondColumn].fill, 0, "")
		pdf.Ln(-1)
	}
	pdf.SetX(left)
	pdf.SetFontStyle("B")
	pdf.CellFormat(columnWidth[0]+columnWidth[1], tableLineHeight, os.Getenv("TEMPLATE1_TOTAL_TITLE"), "1", 0, "", false, 0, "")
	pdf.CellFormat(columnWidth[2]+columnWidth[3], tableLineHeight, strconv.Itoa(len(workedDays)), "1", 0, "C", false, 0, "")
	pdf.Ln(-1)
	pdf.Ln(-1)

	// Page footer
	pdf.SetX(left)
	pdf.SetFontStyle("")
	pdf.CellFormat(pdfAreaWidth/2, lineHeight, os.Getenv("TEMPLATE1_SENDER_SIGNATURE_TITLE"), "", 0, "L", false, 0, "")
	pdf.SetX(pageWidth - marginX - pdfAreaWidth/2)
	pdf.CellFormat(pdfAreaWidth/2, lineHeight, os.Getenv("TEMPLATE1_RECEIVER_SIGNATURE_TITLE"), "", 0, "R", false, 0, "")
}

// Build table entries representing weekdays of the selected month
func buildTableContent(selectedYear int, selectedMonth int, workedDays []int) []tableEntry {
	daysInMonth := time.Date(selectedYear, time.Month(selectedMonth)+1, 0, 0, 0, 0, 0, time.UTC).Day()
	firstWeekdayInMonth := time.Date(selectedYear, time.Month(selectedMonth), 1, 0, 0, 0, 0, time.UTC).Weekday()

	tableEntries := make([]tableEntry, daysInMonth)
	currentWeekDay := firstWeekdayInMonth
	for i := 1; i <= daysInMonth; i++ {
		tableEntries[i-1] = tableEntry{
			day:     fmt.Sprintf("%s %d", translatedDays[currentWeekDay], i),
			present: slices.Contains(workedDays, i),
			fill:    currentWeekDay == time.Saturday || currentWeekDay == time.Sunday || slices.Contains(publicHolidays[selectedMonth], i),
		}
		currentWeekDay = (currentWeekDay + 1) % 7
	}

	return tableEntries
}
