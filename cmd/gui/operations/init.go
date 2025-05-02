package operations

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/lcpu-club/kube-dify-gui/internal/client"
)

// ShowInitialize displays the initialize operation screen
func ShowInitialize(window fyne.Window, client *client.Client, onBack func()) {
	// Create a progress bar
	progress := widget.NewProgressBarInfinite()
	progress.Hide()

	// Create a text area to display the initial password
	passwordDisplay := widget.NewEntry()
	passwordDisplay.Disable()
	passwordDisplay.Hide()

	// Create a button to copy the password
	copyButton := widget.NewButtonWithIcon("Copy Password", theme.ContentCopyIcon(), func() {
		window.Clipboard().SetContent(passwordDisplay.Text)
		dialog.ShowInformation("Copied", "Password copied to clipboard", window)
	})
	copyButton.Hide()

	var initButton *widget.Button

	// Create a button to initialize
	initButton = widget.NewButtonWithIcon("Initialize", theme.DocumentCreateIcon(), func() {
		// Hide the init button and show progress
		initButton.Hide()
		progress.Show()

		// Run the initialization in a goroutine
		go func() {
			// Call the DoInit function
			initPassword, err := client.DoInit()

			// Hide progress
			progress.Hide()

			if err != nil {
				dialog.ShowError(err, window)
				initButton.Show()
				return
			}

			// Show the password
			passwordDisplay.SetText(initPassword)
			passwordDisplay.Show()
			copyButton.Show()

			dialog.ShowInformation("Success", "Initialization completed successfully", window)
		}()
	})

	// Create a back button
	backButton := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), func() {
		onBack()
	})

	// Create layout
	content := container.NewVBox(
		widget.NewLabel("Initialize Dify"),
		widget.NewLabel("This will create the necessary Kubernetes resources for Dify"),
		initButton,
		progress,
		container.NewVBox(
			widget.NewLabel("Initial Password:"),
			passwordDisplay,
			copyButton,
		),
		backButton,
	)

	window.SetContent(content)
}
