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

// ShowStart displays the start operation screen
func ShowStart(window fyne.Window, client *client.Client, onBack func()) {
	// Get translator from the current app
	translator, _ := i18n.New(fyne.CurrentApp())

	// Create UI components
	components := ui.New(translator)

	// Create a progress bar
	progress := widget.NewProgressBarInfinite()
	progress.Hide()

	var startButton *widget.Button

	// Create a button to start
	startButton = widget.NewButtonWithIcon(
		translator.Translate("start.button"),
		theme.MediaPlayIcon(),
		func() {
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

				dialog.ShowInformation(
					translator.Translate("common.success"),
					translator.Translate("start.success"),
					window,
				)
				startButton.Show()
			}()
		},
	)
	startButton.Importance = widget.MediumImportance

	// Create a back button
	backButton := components.CreateBackButton(onBack)

	// Create header
	header := components.CreateHeader(translator.Translate("start.title"))

	// Create info card
	infoCard := components.CreateInfoCard(
		theme.InfoIcon(),
		translator.Translate("start.instruction"),
	)

	// Create main content area
	mainContent := container.NewVBox(
		infoCard,
		startButton,
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
