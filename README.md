Utility to take a screenshot of the latest loan payment on nelnet and send it to a specified email address for reporting.

You should be able to just run the binary in the root of this project and pull the configuration the .env file in the root 
of this project. It should handle your installs for you.

Application can be called by running the binary `main` directly.

## Setup
1. Put this dir somewhere on your machine
2. Run the application with the following env vars: install=true
    - This will run playwright.Install() and exit the program -- only needs done once unless your environment changes
3. Set the rest of the env vars in the .env file to your desired values
    - for your work email password, you will need an application password from gmail -- https://knowledge.workspace.google.com/kb/how-to-create-app-passwords-000009237 
4. Test this with your own email address before sending it to anyone automatically.
    -- IF IT FAILS BE CAREFUL, IT CAN LOCK YOU OUT OF NELNET IF IT FAILS TOO MANY TIMES
5. Set up a cron job to run this every month.


If you want to make any changes, feel free of course. I'd recommend rebuilding the binary after you do so and calling it instead 
of doing `go run main.go` as it will be faster and give you a(n ideally) working snapshot in case you want to tinker more.
