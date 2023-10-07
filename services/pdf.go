package services

import (
	"bytes"
	"encoding/json"
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
	pdf.Ln(25)

	pdf.SetFont("Arial", "", 12)

	headers := []string{"Attribute", "Value"}
	colWidths := []float64{70, 120}

	// Set light purple fill for header
	pdf.SetFillColor(204, 204, 255) // Light purple fill

	// Display Member Data
	for i, h := range headers {
		pdf.CellFormat(colWidths[i], 10, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Set light grey fill for data rows
	pdf.SetFillColor(240, 240, 240) // Light grey fill

	data := [][]string{
		{"First Name", member.User.FirstName},
		{"Last Name", member.User.LastName},
		{"Identifier", member.Identifier},
	}

	if member.JoiningDate != nil {
		jd := member.JoiningDate.Format("2006-01-02")
		data = append(data, []string{"Joining Date", jd})
	}

	for rowIndex, row := range data {
		fill := rowIndex%2 != 0
		for colIndex, d := range row {
			pdf.CellFormat(colWidths[colIndex], 10, d, "1", 0, "L", fill, 0, "")
		}
		pdf.Ln(-1)
	}
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 20, "Attendance Data")
	pdf.Ln(20)
	pdf.SetFont("Arial", "", 12)

	headers = []string{"Punch Time", "Punch", "Status"}
	colWidths = []float64{70, 60, 60}
	pdf.SetFillColor(204, 204, 255) // Light purple fill

	for i, str := range headers {
		pdf.CellFormat(colWidths[i], 10, str, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFillColor(240, 240, 240) // Light grey fill

	var prevDate string
	for rowIndex, a := range attendance {
		currentDate := a.PunchTime.Format("2006-01-02")
		if currentDate != prevDate {
			// Add a new row for the date
			pdf.SetFont("Arial", "B", 12)   // Bold font for date rows
			pdf.SetFillColor(98, 98, 160)   // Darker purple fill for date rows
			pdf.SetTextColor(255, 255, 255) // Set text color to white for date rows
			pdf.CellFormat(190, 10, currentDate, "1", 0, "C", true, 0, "")
			pdf.Ln(-1)
			prevDate = currentDate
			pdf.SetFont("Arial", "", 12)    // Reset to regular font for other rows
			pdf.SetTextColor(0, 0, 0)       // Reset text color to black for other rows
			pdf.SetFillColor(240, 240, 240) // Reset to light grey fill for other rows
		}
		pdf.CellFormat(colWidths[0], 10, a.PunchTime.Format("15:04:05"), "1", 0, "C", rowIndex%2 != 0, 0, "")
		pdf.CellFormat(colWidths[1], 10, helpers.PunchToString(a.Punch), "1", 0, "C", rowIndex%2 != 0, 0, "")
		pdf.CellFormat(colWidths[2], 10, helpers.AttendanceStatusToString(a.Status), "1", 0, "C", rowIndex%2 != 0, 0, "")
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
