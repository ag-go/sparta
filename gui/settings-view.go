package gui

import (
	"log"
	"sparta/crypto"
	"sparta/gui/widgets"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// SettingsView contains the gui information for the settings screen.
func (u *user) SettingsView(window fyne.Window, app fyne.App) fyne.CanvasObject {

	// TODO: Add setting for changing language.

	// Make it possible for the user to switch themes.
	themeSwitcher := widget.NewSelect([]string{"Dark", "Light"}, func(selected string) {
		switch selected {
		case "Dark":
			app.Settings().SetTheme(theme.DarkTheme())
		case "Light":
			app.Settings().SetTheme(theme.LightTheme())
		}

		// Set the theme to the selected one and save it using the preferences api in fyne.
		app.Preferences().SetString("Theme", selected)
	})

	// Default theme is light and thus we set the placeholder to that and then refresh it (without a refresh, it doesn't show until hovering on to widget).
	themeSwitcher.SetSelected(app.Preferences().StringWithFallback("Theme", "Light"))

	// Add the theme switcher next to a label.
	themeChanger := fyne.NewContainerWithLayout(layout.NewGridLayout(2), widget.NewLabel("Application Theme"), themeSwitcher)

	// syncTimeoutSelector gives the user options to set timeout for the send sync command.
	syncTimeoutSelector := widget.NewSelect([]string{"30 seconds", "1 minute", "2 minutes", "5 minutes"}, func(selected string) {
		switch selected {
		case "30 seconds":
			u.Timeout = 1
		case "1 minute":
			u.Timeout = 60
		case "2 minutes":
			u.Timeout = 120
		case "5 minutes":
			u.Timeout = 300
		}

		log.Println(selected)
	})

	// Set the timeout selected to be the default timeout.
	syncTimeoutSelector.SetSelected("30 seconds")

	// timeoutSettings holds the timout settings widget containers.
	timeoutSettings := fyne.NewContainerWithLayout(layout.NewGridLayout(2), widget.NewLabel("Syncronization Timeout"), syncTimeoutSelector)

	// An entry for typing the new username.
	usernameEntry := widgets.NewAdvancedEntry("New Username", false)

	// Create the button used for changing the username.
	usernameButton := widget.NewButtonWithIcon("Change Username", theme.ConfirmIcon(), func() {
		// Check that the username is valid.
		if usernameEntry.Text == u.Password || usernameEntry.Text == "" {
			dialog.ShowInformation("Please enter a valid username", "Usernames need to not be empty and not the same as the password.", window)
		} else {
			// Ask the user to confirm what we are about to do.
			dialog.ShowConfirm("Are you sure that you want to continue?", "The action will permanently change your username.", func(change bool) {
				if change {
					// Calculate the new PasswordKey.
					u.EncryptionKey = crypto.Hash(usernameEntry.Text, u.Password)

					// Set the username  to the updated username.
					u.Username = usernameEntry.Text

					// Clear out the text inside the entry.
					usernameEntry.SetText("")

					// Write the data encrypted using the new key and do so concurrently.
					go u.Data.Write(&u.EncryptionKey)
				}
			}, window)
		}

	})

	// Create the entry for updating the password.
	passwordEntry := widgets.NewAdvancedEntry("New Password", true)

	// Create the button used for changing the password.
	passwordButton := widget.NewButtonWithIcon("Change Password", theme.ConfirmIcon(), func() {
		// Check that the password is valid.
		if len(passwordEntry.Text) < 8 || passwordEntry.Text == usernameEntry.Text {
			dialog.ShowInformation("Please enter a valid password", "Passwords need to be at least eight characters long.", window)
		} else {
			// Ask the user to confirm what we are about to do.
			dialog.ShowConfirm("Are you sure that you want to continue?", "The action will permanently change your password.", func(change bool) {
				if change {
					// Calculate the new PasswordKey.
					u.EncryptionKey = crypto.Hash(u.Username, passwordEntry.Text)

					// Set the user password to the updated password.
					u.Password = passwordEntry.Text

					// Clear out the text inside the entry.
					passwordEntry.SetText("")

					// Write the data encrypted using the new key and do so concurrently.
					go u.Data.Write(&u.EncryptionKey)
				}
			}, window)
		}
	})

	// Extend our extended buttons with array entry switching and enter to change.
	usernameEntry.InitExtend(*usernameButton, widgets.MoveAction{Down: true, DownEntry: passwordEntry, Window: window})
	passwordEntry.InitExtend(*passwordButton, widgets.MoveAction{Up: true, UpEntry: usernameEntry, Window: window})

	// revertToDefaultSettings reverts all settings to their default values.
	revertToDefaultSettings := widget.NewButtonWithIcon("Reset settings to default values", theme.ViewRefreshIcon(), func() {
		// Update theme and saved settings for theme change.
		if app.Preferences().String("Theme") != "Light" {
			themeSwitcher.SetSelected("Light")

			// Set the visible theme to the light theme.
			app.Settings().SetTheme(theme.LightTheme())

			// Set the saved theme to Light.
			app.Preferences().SetString("Theme", "Light")
		}
	})
	// Create a button for clearing the data of a given profile.
	deleteButton := widget.NewButtonWithIcon("Delete all saved activities", theme.DeleteIcon(), func() {

		// Ask the user to confirm what we are about to do.
		dialog.ShowConfirm("Are you sure that you want to continue?", "Deleting your data will remove all of your exercises and activities.", func(remove bool) {
			if remove {
				// Run the delete function and do it concurrently to avoid stalling the thread with file io.
				go u.Data.Delete()

				// Notify the label that we have removed the data.
				u.EmptyExercises <- true
			}
		}, window)
	})

	// userInterfaceSettings is a group holding widgets related to user interface settings such as theme.
	userInterfaceSettings := widget.NewGroup("User Interface Settings", themeChanger)

	// syncSettings is the group for all settings related to sync support.
	syncSettings := widget.NewGroup("Syncronization Settings", timeoutSettings)

	// credentialSettings groups together all settings related to usernames and passwords.
	credentialSettings := widget.NewGroup("Login Credential Settings", fyne.NewContainerWithLayout(layout.NewGridLayout(2), usernameEntry, usernameButton, passwordEntry, passwordButton))

	// advancedSettings is a group holding widgets related to advanced settings.
	advancedSettings := widget.NewGroup("Advanced Settings", revertToDefaultSettings, widget.NewLabel(""), deleteButton)

	// settingsContentView holds all widget groups and content for the settings page.
	settingsContentView := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), userInterfaceSettings, layout.NewSpacer(), syncSettings, layout.NewSpacer(), credentialSettings, layout.NewSpacer(), advancedSettings)

	return widget.NewScrollContainer(settingsContentView)
}
