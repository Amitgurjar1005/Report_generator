package service

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

func SendReport(c *gin.Context) {

	// parse multipart
	if err := c.Request.ParseMultipartForm(100 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "multipart error"})
		return
	}

	email := c.PostForm("email")

	file, err := c.FormFile("pdf")
	if err != nil {
		fmt.Println("FILE ERROR:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "pdf missing"})
		return
	}

	// create folder
	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	pdfPath := filepath.Join(uploadDir, "report.pdf")

	// save PDF
	if err := c.SaveUploadedFile(file, pdfPath); err != nil {
		fmt.Println("SAVE ERROR:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save failed"})
		return
	}

	// confirm file exists
	if stat, err := os.Stat(pdfPath); err != nil || stat.Size() == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid pdf"})
		return
	}

	// ================= ASYNC EMAIL =================
	go processAndSend(email, pdfPath)

	c.JSON(http.StatusOK, gin.H{
		"message": "Report processing started, email will be sent",
	})
}

func createZip(pdfPath string) (string, error) {

	zipPath := "./uploads/report.zip"

	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", err
	}

	zipWriter := zip.NewWriter(zipFile)

	fileToZip, err := os.Open(pdfPath)
	if err != nil {
		zipWriter.Close()
		zipFile.Close()
		return "", err
	}

	writer, err := zipWriter.Create("website-report.pdf")
	if err != nil {
		fileToZip.Close()
		zipWriter.Close()
		zipFile.Close()
		return "", err
	}

	_, err = io.Copy(writer, fileToZip)

	// IMPORTANT ORDER (FIXES NULL ZIP ISSUE)
	fileToZip.Close()
	zipWriter.Close()
	zipFile.Close()

	if err != nil {
		return "", err
	}

	// DEBUG
	info, _ := os.Stat(zipPath)
	fmt.Println("ZIP SIZE:", info.Size())

	return zipPath, nil
}
func processAndSend(email string, pdfPath string) {

	zipPath, err := createZip(pdfPath)
	if err != nil {
		fmt.Println("ZIP ERROR:", err)
		return
	}

	// validate zip
	if stat, err := os.Stat(zipPath); err != nil || stat.Size() == 0 {
		fmt.Println("ZIP INVALID OR EMPTY")
		return
	}

	err = sendEmail(email, zipPath)
	if err != nil {
		fmt.Println("EMAIL ERROR:", err)
	} else {
		fmt.Println("EMAIL SENT SUCCESSFULLY")
	}

	// cleanup
	os.Remove(pdfPath)
	os.Remove(zipPath)
}

func sendEmail(to string, filePath string) error {

	m := gomail.NewMessage()

	m.SetHeader("From", "webinsightreports@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Your Website Analysis Report (ZIP File)")

	m.SetBody("text/plain",
		"Hello,\n\nPlease find your website analysis report attached as ZIP file.")

	m.Attach(filePath)

	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		"webinsightreports@gmail.com",
		"qljr ured twgr rwru", // 👈 put real app password here
	)

	return d.DialAndSend(m)
}
