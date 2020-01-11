package gui

import (
	"sparta/src/file"
	"sparta/src/file/encrypt"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

// PasswordKey contains the key taken from the username and password.
var PasswordKey [32]byte

// CheckValidInput checks if the inputed username and passwords are valid adn creates a message if they are not.
func CheckValidInput(username, password string, window fyne.Window) (valid bool) {
	if username == "" || password == "" {
		dialog.ShowInformation("Missing username/password", "Please provide both username and password.", window)
		valid = false
	} else if username == password {
		dialog.ShowInformation("Identical username and password", "Please do not use identical username and password.", window)
		valid = false
	} else if len(password) < 8 {
		dialog.ShowInformation("Too short password", "The password should be eight characters or longer.", window)
		valid = false
	}

	return true
}

// ShowLoginPage shows the login page that handles the inertaface for logging in.
func ShowLoginPage(window fyne.Window) {
	// Initialize the login form that we are to be using.
	username := NewExtendedEntry("Username", false)

	// Initialize the password input box that we are to be using.
	password := NewExtendedEntry("Password", true)

	// Create the login button that will calculate the 32bit long sha256 hash.
	loginButton := widget.NewButton("Login", func() {
		// Check the inputed data to handle invalid inputs.
		valid := CheckValidInput(username.Text, password.Text, window)
		if !valid {
			return
		}

		// Adapt the window to a good size and make it resizable again.
		window.SetFixedSize(false)
		window.Resize(fyne.NewSize(800, 500))

		// Calculate the sha256 hash of the username and password.
		PasswordKey = encrypt.EncryptionKey(username.Text, password.Text)

		// Create a channel for sending activity data through. Let's us avoid reading the file every time we add a new activity.
		newAddedExercise := make(chan string)

		// Check for the file where we store the data.
		XMLData := file.Check(&PasswordKey)

		ShowMainDataView(window, &XMLData, newAddedExercise)
	})

	// Add the Action component to make actions work inside the struct. This is used to press the loginButton on pressing enter/return ton the keyboard.
	username.Action, password.Action = &Action{*loginButton}, &Action{*loginButton}

	// Set the content to be displayed. It is the userName, userPassword fields and the login button inside a layout.
	window.SetContent(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), username, password, loginButton))

	// Set a sane default for the window size on login and make sure that it isn't resizable.
	window.Resize(fyne.NewSize(400, 100))
	window.SetFixedSize(true)

	// Show all of our set content and run the gui.
	window.ShowAndRun()
}