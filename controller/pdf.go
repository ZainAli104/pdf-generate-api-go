package controller

import (
	"github.com/gin-gonic/gin"
	"pdf-generate/services"
)

func NewPdf(group *gin.RouterGroup) {
	c := NewPdfController{PdfService: services.NewPdfService()}
	group.GET("", c.GetPdf)
}

type NewPdfController struct {
	PdfService *services.PdfService
}

func (u *NewPdfController) GetPdf(c *gin.Context) {
	pdfBytes, err := u.PdfService.GeneratePdf()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.Data(200, "application/pdf", pdfBytes)
}
