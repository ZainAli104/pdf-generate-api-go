package services

import (
	"bytes"
	"encoding/json"
	"github.com/jung-kurt/gofpdf"
	"net/http"
	"os"
	"pdf-generate/helpers"
	"pdf-generate/models"
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

	// Display Attendance Data
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

	for rowIndex, a := range attendance {
		fill := rowIndex%2 != 0
		pdf.CellFormat(colWidths[0], 10, a.PunchTime.Format("2006-01-02 15:04:05"), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colWidths[1], 10, helpers.PunchToString(a.Punch), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colWidths[2], 10, helpers.AttendanceStatusToString(a.Status), "1", 0, "C", fill, 0, "")
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
