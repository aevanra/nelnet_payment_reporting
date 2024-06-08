package email

import (
    "bytes"
    "encoding/base64"
    "fmt"
    "log"
    "mime"
    "mime/multipart"
    "net/smtp"
    "net/http"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "time"

    "github.com/joho/godotenv"
    "github.com/emersion/go-message/charset"
    "github.com/emersion/go-imap/v2/imapclient"
    "github.com/emersion/go-imap/v2"
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
    defer writer.Close()
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
	}

	return buf.Bytes()

}

func SendEmail(recipients []string, attachmentPath string, fullName string, ) {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    today := time.Now()

    auth := smtp.PlainAuth("", os.Getenv("sender_address"), os.Getenv("sender_password"), "smtp.gmail.com")
    m := new(message)
    m.To = recipients
    m.Subject = fullName + " Student Loan Repayment " + today.Format("2006-01-02")
    m.Body = "Hello,\n\nAttached is my receipt for my student loan payment this month."

    m.Attachments = make(map[string][]byte)
    m.attachFile(attachmentPath)
    
    msg := m.build()

    err = smtp.SendMail("smtp.gmail.com:587", auth, os.Getenv("sender_address"), m.To, msg)
    if err != nil {
        log.Fatal(err)
    }
}

func GetMFACode() string {
    options := &imapclient.Options{
        WordDecoder: &mime.WordDecoder{CharsetReader: charset.Reader},
    }
    client, err := imapclient.DialTLS("imap.gmail.com:993", options)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    log.Print("Connected")
    
    // Login
    client.Login(os.Getenv("mfa_email_address"), os.Getenv("mfa_email_password")).Wait()
    defer client.Logout()
    log.Print("Logged in")

    // Select inbox
    opt := &imap.SelectOptions{ReadOnly: true}
    client.Select("INBOX", opt).Wait()
    log.Print("Selected inbox")

    // Search for email
    headerSearch := []imap.SearchCriteriaHeaderField{}
    headerSearch = append(headerSearch, imap.SearchCriteriaHeaderField{Key: "SUBJECT", Value: "Your authentication code"})
    // the Since only evaluates the date, not the time. So time.Now() means just today
    searchCriteria := imap.SearchCriteria{Since: time.Now(), Header: headerSearch}
    dat, err := client.Search(&searchCriteria, &imap.SearchOptions{ReturnAll: true}).Wait()
    log.Print(dat.All)

    bodySectionList :=[]*imap.FetchItemBodySection{
            {Specifier: imap.PartSpecifierText},
        }

    // fetch the auth email
    fetchOptions := imap.FetchOptions{
        Flags: true, 
        Envelope: true, 
        BodySection: bodySectionList,
    }
    messages, err := client.Fetch(dat.All, &fetchOptions).Collect()
    if err != nil {
        log.Fatal("Could not fetch given seq")
    }

    var msg string
    for _, v := range messages[0].BodySection{
        msg = string(v[:])
        break
    }

    re, _ := regexp.Compile("[0-9]{6}</p>")
    match := re.Find([]byte(msg))
    matchStr := string(match[:])

    return strings.Replace(matchStr, "</p>", "", -1)
}
