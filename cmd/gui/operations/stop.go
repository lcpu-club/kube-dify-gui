package operations

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/lcpu-club/kube-dify-gui/internal/client"
)

// ShowStop displays the stop operation screen
func ShowStop(window fyne.Window, client *client.Client, onBack func()) {
	// Create a progress bar
	progress := widget.NewProgressBarInfinite()
	progress.Hide()

	var stopButton *widget.Button

	// Create a button to stop
	stopButton = widget.NewButtonWithIcon("Stop", theme.MediaStopIcon(), func() {
		// Show confirmation dialog
		dialog.ShowConfirm("Confirm Stop", "Are you sure you want to stop Dify?", func(confirmed bool) {
			if !confirmed {
				return
			}

			// Hide the stop button and show progress
			stopButton.Hide()
			progress.Show()

			// Run the stop operation in a goroutine
			go func() {
				// Call the DoStop function
				err := client.DoStop()

				// Hide progress
				progress.Hide()

				if err != nil {
					dialog.ShowError(err, window)
					stopButton.Show()
					return
				}

				dialog.ShowInformation("Success", "Dify stopped successfully", window)
				stopButton.Show()
			}()
		}, window)
	})

	// Create a back button
	backButton := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), func() {
		onBack()
	})

	// Create layout
	content := container.NewVBox(
		widget.NewLabel("Stop Dify"),
		widget.NewLabel("This will stop the Dify application"),
		stopButton,
		progress,
		backButton,
	)

	window.SetContent(content)
}
