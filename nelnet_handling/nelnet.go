package nelnet

import (
    "errors"
    "log"
    "os"
    "time"

    "github.com/playwright-community/playwright-go"
)

func GetMostRecentPaymentScreenshot(outputPath string) error {
    headless := os.Getenv("headless") == "true"
    if os.Getenv("install") == "true" {
        // Install playwright
        playwright.Install()

        if !(os.Getenv("hold_browser_open") == "true") {
            return errors.New("Installed playwright")
        }
    }

    pw, err := playwright.Run()
    if err != nil {
        log.Fatalf("Could not start playwright: %v", err)
    } 

    // Browser Config
    browserConfig := playwright.BrowserTypeLaunchPersistentContextOptions{Headless: &headless, 
        RecordHarMode: playwright.HarModeFull, 
        Screen: &playwright.Size{Width: 2560, Height: 1440},
    }

    // Start Browser
    browser, err := pw.Firefox.LaunchPersistentContext(os.Getenv("firefox_profile"), browserConfig)
    defer browser.Close()
    if err != nil {
        log.Fatalf("Could not start browser: %v", err)
    }

    // Get the open page
    page := browser.Pages()[0]

    // Navigate to Nelnet Login Page
    _, err = page.Goto("https://nelnet.studentaid.gov/dashboard", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
    if err != nil {
        log.Fatalf("Could not navigate to page: %v", err)
    }
    // There are annoying pages to get through
    page.Locator("[id='login-button']").Click()
    page.Locator("[id='user-name-button']").Click()

    // Login
    time.Sleep(1*time.Second) // Not sleeping caused issues
    page.Locator("[id='username-textfield']").Fill(os.Getenv("nelnet_username"))
    page.Locator("[id='password-textfield']").Fill(os.Getenv("nelnet_password"))
    time.Sleep(1*time.Second) // Not sleeping caused issues
    page.GetByText("Continue").Click()


    // Hold browser open for 5 minutes if needed
    if !headless && os.Getenv("hold_browser_open") == "true" {
        log.Println("Holding browser open for 5 minutes")
        time.Sleep(10*time.Minute)
        return errors.New("Held browser open for maintenance, exiting now.")
    }

    // Go to Payment Activity Screen
    page.Goto("https://nelnet.studentaid.gov/payments/payment-activity")

    // Most Recent Payment
    page.Locator(".wide-table > tbody:nth-child(2) > tr:nth-child(1) > td:nth-child(1) > a:nth-child(1)").Click()

    fullPage := true
    _, err = page.Screenshot(playwright.PageScreenshotOptions{Path: &outputPath, FullPage: &fullPage})
    if err != nil {
        log.Fatalf("Could not take screenshot: %v", err)
    }

    return nil
}
