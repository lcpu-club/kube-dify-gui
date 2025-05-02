package operations

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/lcpu-club/kube-dify-gui/internal/client"
)

// ShowStart displays the start operation screen
func ShowStart(window fyne.Window, client *client.Client, onBack func()) {
	// Create a progress bar
	progress := widget.NewProgressBarInfinite()
	progress.Hide()

	var startButton *widget.Button

	// Create a button to start
	startButton = widget.NewButtonWithIcon("Start", theme.MediaPlayIcon(), func() {
		// Hide the start button and show progress
		startButton.Hide()
		progress.Show()

		// Run the start operation in a goroutine
		go func() {
			// Call the DoStart function
			err := client.DoStart()

			// Hide progress
			progress.Hide()

			if err != nil {
				dialog.ShowError(err, window)
				startButton.Show()
				return
			}

			dialog.ShowInformation("Success", "Dify started successfully", window)
			startButton.Show()
		}()
	})

	// Create a back button
	backButton := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), func() {
		onBack()
	})

	// Create layout
	content := container.NewVBox(
		widget.NewLabel("Start Dify"),
		widget.NewLabel("This will start the Dify application"),
		startButton,
		progress,
		backButton,
	)

	window.SetContent(content)
}
