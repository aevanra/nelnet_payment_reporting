Utility to take a screenshot of the latest loan payment on nelnet and send it to a specified email address for reporting.

You should be able to just run the binary in the root of this project and pull the configuration the .env file in the root 
of this project. It should handle your installs for you.

Application can be called by running the binary `main` directly.

## Setup
1. Put this dir somewhere on your maching (you did it let's go)
2. Run the application with the following env vars: install=true
    - This will run playwright.Install() and exit the program -- only needs done once unless your environment changes
3. Create a firefox profile for the app to use 
    - run the application with run_playwright_headless=false and hold_browser_open=true
    - go to about:profiles in firefox and create a new profile
4. Manually login to nelnet in the firefox instance that pops up and set nelnet to trust the browser 
5. Set the firefox_profile env var to the path of the profile you just created
6. Set the rest of the env vars in the .env file to your desired values
    - for your work email password, you will need an application password from gmail -- https://knowledge.workspace.google.com/kb/how-to-create-app-passwords-000009237 
7. Test this with your own email address before sending it to HR automatically.
8. Set up a cron job to run this every month.

### MAJOR ANNOYANCE ###
Nelnet's MFA requires you to authenticate every 90 days. This means you will need to log in in the firefox instance this spins
up and 2FA every 90 days. Currently, this does not report that, so you will need to set a reminder for yourself. You can hold the
browser open with the .env file by setting run_playwright_headless env var as false and the hold_browser_open env var as true.
At present, I don't have this automated, but feel free to extend this for yourself, I don't imagine it will end up being too
crazy to get an email code and auth with your specific email provider.


If you want to make any changes, feel free of course. I'd recommend rebuilding the binary after you do so and calling it instead 
of doing `go run main.go` as it will be faster and give you a(n ideally) working snapshot in case you want to tinker more.
