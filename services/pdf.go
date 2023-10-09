package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"net/http"
	"os"
	"pdf-generate/helpers"
	"pdf-generate/models"
	"sort"
)

func NewPdfService() *PdfService {
	return &PdfService{}
}

type PdfService struct{}

func (p *PdfService) GeneratePdf() ([]byte, error) {
	member, err := p.fetchMemberData()
	if err != nil {
		return nil, err
	}

	attendance, err := p.fetchMemberAttendanceData()
	if err != nil {
		return nil, err
	}

	sort.Slice(attendance, func(i, j int) bool {
		return attendance[i].PunchTime.Before(attendance[j].PunchTime)
	})

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 24)
	pdf.Cell(0, 20, "Member Report")
	pdf.Ln(20)

	pdf.SetFont("Arial", "", 12)

	headers := []string{"Attribute", "Value"}
	colWidths := []float64{70, 120}

	// Display => Member Data
	headers = []string{"First Name", "Last Name", "Identifier", "Joining Date"}
	colWidths = []float64{50, 50, 60, 30}

	// Set light purple fill for header
	pdf.SetFillColor(204, 204, 255) // Light purple fill

	for i, h := range headers {
		pdf.CellFormat(colWidths[i], 10, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Set light grey fill for data rows
	pdf.SetFillColor(240, 240, 240) // Light grey fill

	data := []string{
		member.User.FirstName,
		member.User.LastName,
		member.Identifier,
	}

	if member.JoiningDate != nil {
		data = append(data, member.JoiningDate.Format("2006-01-02"))
	} else {
		data = append(data, "") // If JoiningDate is nil, we add an empty string to align the table cells
	}

	for i, d := range data {
		pdf.CellFormat(colWidths[i], 10, d, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Display => Attendance Data
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 20, "Attendance Data")
	pdf.Ln(20)
	pdf.SetFont("Arial", "", 12)

	headers = []string{"Punch Time", "Type", "Status", "Fine"}
	colWidths = []float64{55, 55, 55, 25}
	pdf.SetFillColor(204, 204, 255) // Light purple fill

	for i, str := range headers {
		pdf.CellFormat(colWidths[i], 10, str, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFillColor(240, 240, 240) // Light grey fill

	var prevDate string
	totalWidth := float64(0)
	for _, w := range colWidths {
		totalWidth += w
	}

	for rowIndex, a := range attendance {
		currentDate := a.PunchTime.Format("2006-01-02")
		if currentDate != prevDate {
			// Add a new row for the date
			pdf.SetFont("Arial", "B", 12)   // Bold font for date rows
			pdf.SetFillColor(98, 98, 160)   // Darker purple fill for date rows
			pdf.SetTextColor(255, 255, 255) // Set text color to white for date rows
			pdf.CellFormat(totalWidth, 10, currentDate, "1", 0, "C", true, 0, "")
			pdf.Ln(-1)
			prevDate = currentDate
			pdf.SetFont("Arial", "", 12)    // Reset to regular font for other rows
			pdf.SetTextColor(0, 0, 0)       // Reset text color to black for other rows
			pdf.SetFillColor(240, 240, 240) // Reset to light grey fill for other rows
		}

		// Calculate the fine percentage
		var finePercentage float64 = 0
		if a.Status == models.AttendanceLate {
			if a.PunchTime.Hour() == 9 {
				if a.PunchTime.Minute() > 0 && a.PunchTime.Minute() <= 10 {
					finePercentage = 1
				} else if a.PunchTime.Minute() > 10 && a.PunchTime.Minute() <= 20 {
					finePercentage = 2
				} else if a.PunchTime.Minute() > 20 {
					finePercentage = 3
				}
			} else if a.PunchTime.Hour() > 9 {
				finePercentage = 3
			}
		}
		fine := fmt.Sprintf("%.0f%%", finePercentage)

		pdf.CellFormat(colWidths[0], 10, a.PunchTime.Format("15:04:05"), "1", 0, "C", rowIndex%2 != 0, 0, "")
		pdf.CellFormat(colWidths[1], 10, helpers.PunchToString(a.Punch), "1", 0, "C", rowIndex%2 != 0, 0, "")
		pdf.CellFormat(colWidths[2], 10, helpers.AttendanceStatusToString(a.Status), "1", 0, "C", rowIndex%2 != 0, 0, "")
		pdf.CellFormat(colWidths[3], 10, fine, "1", 0, "C", rowIndex%2 != 0, 0, "")
		pdf.Ln(-1)
	}

	var pdfBytes bytes.Buffer
	err = pdf.Output(&pdfBytes)
	if err != nil {
		return nil, err
	}

	return pdfBytes.Bytes(), nil
}

func (p *PdfService) fetchMemberData() (*models.MemberPopulated, error) {
	response, err := http.Get(os.Getenv("MEMBER_DATA_API"))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	member := &models.MemberPopulated{}
	err = json.NewDecoder(response.Body).Decode(member)
	if err != nil {
		return nil, err
	}

	return member, nil
}

type AttendanceResponse struct {
	Data []models.Attendance `json:"data"`
}

func (p *PdfService) fetchMemberAttendanceData() ([]models.Attendance, error) {
	response, err := http.Get(os.Getenv("MEMBER_ATTENDANCE_API"))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var attendanceResp AttendanceResponse
	err = json.NewDecoder(response.Body).Decode(&attendanceResp)
	if err != nil {
		return nil, err
	}

	return attendanceResp.Data, nil
}
