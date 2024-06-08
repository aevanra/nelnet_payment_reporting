package main

import (
    "log"
    "os"
    "strings"
    "time"

    "github.com/aevanra/expenses_automation/nelnet_handling"
    "github.com/aevanra/expenses_automation/email_handling"
    "github.com/joho/godotenv"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
    // config
    now := time.Now()
    emailRecipients := strings.Split(os.Getenv("email_recipients"), ",")
    outputFilepath := os.Getenv("output_dir") + "/nelnet_payment_" + now.Format("2006-01-02") + ".png"

    // Steps
    err = nelnet.GetMostRecentPaymentScreenshot(outputFilepath)
    if err != nil {
        log.Fatalf(err.Error())
    }

    if os.Getenv("send_email") == "true" {
        email.SendEmail(emailRecipients, outputFilepath, os.Getenv("full_name"))
    }
    log.Fatalf("Finished running script")

}
