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

// ShowStop displays the stop operation screen
func ShowStop(window fyne.Window, client *client.Client, onBack func()) {
	// Get translator from the current app
	translator, _ := i18n.New(fyne.CurrentApp())

	// Create UI components
	components := ui.New(translator)

	// Create a progress bar
	progress := widget.NewProgressBarInfinite()
	progress.Hide()

	var stopButton *widget.Button

	// Create a button to stop
	stopButton = widget.NewButtonWithIcon(
		translator.Translate("stop.button"),
		theme.MediaStopIcon(),
		func() {
			// Show confirmation dialog
			dialog.ShowConfirm(
				translator.Translate("stop.title"),
				translator.Translate("stop.confirm"),
				func(confirmed bool) {
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

						dialog.ShowInformation(
							translator.Translate("common.success"),
							translator.Translate("stop.success"),
							window,
						)
						stopButton.Show()
					}()
				},
				window,
			)
		},
	)
	stopButton.Importance = widget.MediumImportance

	// Create a back button
	backButton := components.CreateBackButton(onBack)

	// Create header
	header := components.CreateHeader(translator.Translate("stop.title"))

	// Create info card
	infoCard := components.CreateInfoCard(
		theme.InfoIcon(),
		translator.Translate("stop.instruction"),
	)

	// Create main content area
	mainContent := container.NewVBox(
		infoCard,
		stopButton,
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
