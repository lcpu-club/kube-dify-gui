package operations

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/lcpu-club/kube-dify-gui/internal/client"
)

// ShowDelete displays the delete operation screen
func ShowDelete(window fyne.Window, client *client.Client, onBack func()) {
	// Create a progress bar
	progress := widget.NewProgressBarInfinite()
	progress.Hide()

	var deleteButton *widget.Button

	// Create a button to delete
	deleteButton = widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
		// Show confirmation dialog with warning
		dialog.ShowConfirm(
			"Confirm Delete",
			"WARNING: This will delete all Dify resources including data. This action cannot be undone. Are you sure?",
			func(confirmed bool) {
				if !confirmed {
					return
				}

				// Hide the delete button and show progress
				deleteButton.Hide()
				progress.Show()

				// Run the delete operation in a goroutine
				go func() {
					// Call the DoDelete function
					err := client.DoDelete()

					// Hide progress
					progress.Hide()

					if err != nil {
						dialog.ShowError(err, window)
						deleteButton.Show()
						return
					}

					dialog.ShowInformation("Success", "Dify resources deleted successfully", window)
					deleteButton.Show()
				}()
			},
			window,
		)
	})

	// Create a back button
	backButton := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), func() {
		onBack()
	})

	// Create layout with warning
	content := container.NewVBox(
		widget.NewLabel("Delete Dify"),
		widget.NewLabelWithStyle(
			"WARNING: This will delete all Dify resources including data!",
			fyne.TextAlignCenter,
			fyne.TextStyle{Bold: true},
		),
		deleteButton,
		progress,
		backButton,
	)

	window.SetContent(content)
}
