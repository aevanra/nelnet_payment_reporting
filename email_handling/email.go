package email

import (
    "bytes"
    "encoding/base64"
    "fmt"
    "log"
    "mime/multipart"
    "net/smtp"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/joho/godotenv"
)

type message struct {
    To []string
    Subject string
    Body string
    Attachments map[string][]byte
}

func (m *message) attachFile(filePath string)  error {
    // add file to the attachments map
    file, err := os.ReadFile(filePath)
    if err != nil {
        return err
    }
    
    var filename string

    if strings.Contains(filePath, "/") { 
        _, filename = filepath.Split(filePath)
    } else {
        filename = filePath
    }
    
    m.Attachments[filename] = file

    return nil
}

func (m *message) build() []byte {
	hasAttachments := len(m.Attachments) > 0

    // Create an empty buffer for the email
    buf := bytes.NewBuffer(nil)

    // Write the email headers
	buf.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ",")))
	buf.WriteString("MIME-Version: 1.0\n")

    // Create a new multipart writer and set the boundary
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
    
    // If we have attachments, set the content type to multipart/mixed
	if hasAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	} else {
        // Otherwise, set the content type to text/plain
		buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
	}
    
    // Write the body of the email
	buf.WriteString(m.Body)

    // If we have attachments, write them to the email
	if hasAttachments {
		for k, v := range m.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}
        
        // Close the writer
        writer.Close()
	}

	return buf.Bytes()

}

func SendEmail(recipients []string, attachmentPath string, fullName string, ) {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    today := time.Now()

    auth := smtp.PlainAuth("", os.Getenv("gmail_address"), os.Getenv("gmail_password"), "smtp.gmail.com")
    m := new(message)
    m.To = recipients
    m.Subject = fullName + " Student Loan Repayment " + today.Format("2006-01-02")
    m.Body = "Hello,\n\nAttached is my receipt for my student loan payment this month."

    m.Attachments = make(map[string][]byte)
    m.attachFile(attachmentPath)
    
    msg := m.build()

    err = smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("gmail_address"), m.To, msg)
    if err != nil {
        log.Fatal(err)
    }

}
