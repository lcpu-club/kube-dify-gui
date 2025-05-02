package operations

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/lcpu-club/kube-dify-gui/internal/client"
	"github.com/lcpu-club/kube-dify-gui/internal/i18n"
	"github.com/lcpu-club/kube-dify-gui/internal/ui"
)

// ShowDelete displays the delete operation screen
func ShowDelete(window fyne.Window, client *client.Client, onBack func()) {
	// Get translator from the current app
	translator, _ := i18n.New(fyne.CurrentApp())

	// Create UI components
	components := ui.New(translator)

	// Create a progress bar
	progress := widget.NewProgressBarInfinite()
	progress.Hide()

	var deleteButton *widget.Button

	// Create a button to delete
	deleteButton = widget.NewButtonWithIcon(
		translator.Translate("delete.button"),
		theme.DeleteIcon(),
		func() {
			// Show confirmation dialog with warning
			dialog.ShowConfirm(
				translator.Translate("delete.title"),
				translator.Translate("delete.confirm"),
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

						dialog.ShowInformation(
							translator.Translate("common.success"),
							translator.Translate("delete.success"),
							window,
						)
						deleteButton.Show()
					}()
				},
				window,
			)
		},
	)
	deleteButton.Importance = widget.HighImportance

	// Create a back button
	backButton := components.CreateBackButton(onBack)

	// Create header
	header := components.CreateHeader(translator.Translate("delete.title"))

	// Create warning card
	warningCard := components.CreateInfoCard(
		theme.WarningIcon(),
		translator.Translate("delete.instruction"),
	)

	// Create main content area
	mainContent := container.NewVBox(
		warningCard,
		deleteButton,
		progress,
	)

	// Create card for main content
	contentCard := components.CreateCard("", mainContent)

	// Create layout
	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		contentCard,
		backButton,
	)

	window.SetContent(container.NewPadded(content))
}
