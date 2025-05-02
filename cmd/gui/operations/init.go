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

// ShowInitialize displays the initialize operation screen
func ShowInitialize(window fyne.Window, client *client.Client, onBack func()) {
	// Get translator from the current app
	translator, _ := i18n.New(fyne.CurrentApp())

	// Create UI components
	components := ui.New(translator)

	// Create a progress bar
	progress := widget.NewProgressBarInfinite()
	progress.Hide()

	// Create a text area to display the initial password
	passwordDisplay := widget.NewEntry()
	passwordDisplay.Disable()
	passwordDisplay.Hide()

	// Create a button to copy the password
	copyButton := widget.NewButtonWithIcon(
		translator.Translate("init.copy"),
		theme.ContentCopyIcon(),
		func() {
			window.Clipboard().SetContent(passwordDisplay.Text)
			dialog.ShowInformation(
				translator.Translate("common.success"),
				translator.Translate("init.copied"),
				window,
			)
		},
	)
	copyButton.Hide()

	var initButton *widget.Button

	// Create a button to initialize
	initButton = widget.NewButtonWithIcon(
		translator.Translate("init.button"),
		theme.DocumentCreateIcon(),
		func() {
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

				dialog.ShowInformation(
					translator.Translate("common.success"),
					translator.Translate("init.success"),
					window,
				)
			}()
		},
	)

	// Create a back button
	backButton := components.CreateBackButton(onBack)

	// Create header
	header := components.CreateHeader(translator.Translate("init.title"))

	// Create info card with explanation
	infoCard := components.CreateInfoCard(
		theme.InfoIcon(),
		translator.Translate("init.instruction"),
	)

	// Create password section
	passwordSection := container.NewVBox(
		widget.NewLabel(translator.Translate("init.password")),
		passwordDisplay,
		copyButton,
	)

	// Create main content area
	mainContent := container.NewVBox(
		infoCard,
		initButton,
		progress,
		passwordSection,
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
