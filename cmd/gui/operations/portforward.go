package operations

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/lcpu-club/kube-dify-gui/internal/client"
	"github.com/lcpu-club/kube-dify-gui/internal/i18n"
	"github.com/lcpu-club/kube-dify-gui/internal/ui"
)

// ShowPortForward displays the port forward operation screen
func ShowPortForward(window fyne.Window, client *client.Client, onBack func()) {
	// Get translator from the current app
	translator, _ := i18n.New(fyne.CurrentApp())

	// Create UI components
	components := ui.New(translator)

	// Create input fields
	addressEntry := widget.NewEntry()
	addressEntry.SetPlaceHolder(translator.Translate("portforward.address.placeholder"))
	addressEntry.SetText("127.0.0.1")

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder(translator.Translate("portforward.port.placeholder"))
	portEntry.SetText("28080")

	// Create status label
	statusText := fmt.Sprintf("Status: %s", translator.Translate("portforward.status.not_forwarding"))
	statusLabel := widget.NewLabel(statusText)

	// Create log display
	logDisplay := widget.NewMultiLineEntry()
	logDisplay.Disable()
	logDisplay.SetMinRowsVisible(10)

	// Variables to store port forwarding state
	var stopChan chan struct{}
	var outStream, errStream io.ReadCloser
	var isForwarding bool

	// Function to update logs from streams
	updateLogs := func(stream io.ReadCloser, prefix string) {
		buffer := make([]byte, 1024)
		for {
			n, err := stream.Read(buffer)
			if err != nil {
				if err != io.EOF && !strings.Contains(err.Error(), "file already closed") {
					fyne.CurrentApp().SendNotification(&fyne.Notification{
						Title:   translator.Translate("portforward.title"),
						Content: err.Error(),
					})
				}
				break
			}
			if n > 0 {
				text := fmt.Sprintf("[%s] %s", prefix, string(buffer[:n]))
				logDisplay.SetText(logDisplay.Text + text)
				// Scroll to bottom
				logDisplay.CursorRow = len(strings.Split(logDisplay.Text, "\n")) - 1
			}
		}
	}

	// Create buttons
	var startButton, stopButton, openButton *widget.Button

	// Stop button
	stopButton = widget.NewButtonWithIcon(
		translator.Translate("portforward.stop"),
		theme.MediaStopIcon(),
		func() {
			if stopChan != nil {
				close(stopChan)
				stopChan = nil
			}

			if outStream != nil {
				outStream.Close()
				outStream = nil
			}

			if errStream != nil {
				errStream.Close()
				errStream = nil
			}

			isForwarding = false
			statusText := fmt.Sprintf("Status: %s", translator.Translate("portforward.status.not_forwarding"))
			statusLabel.SetText(statusText)
			startButton.Enable()
			stopButton.Disable()
			openButton.Disable()
		},
	)
	stopButton.Disable()

	// Start button
	startButton = widget.NewButtonWithIcon(
		translator.Translate("portforward.start"),
		theme.MediaPlayIcon(),
		func() {
			address := addressEntry.Text
			port := portEntry.Text

			if address == "" {
				dialog.ShowError(fmt.Errorf(translator.Translate("error.address_empty")), window)
				return
			}

			if port == "" {
				dialog.ShowError(fmt.Errorf(translator.Translate("error.port_empty")), window)
				return
			}

			// Disable start button and enable stop button
			startButton.Disable()
			stopButton.Enable()

			// Clear log display
			logDisplay.SetText("")

			// Update status
			statusText := fmt.Sprintf("Status: %s", translator.Translate("portforward.status.starting"))
			statusLabel.SetText(statusText)

			// Run port forward in a goroutine
			go func() {
				var err error
				stopChan, outStream, errStream, err = client.DoPortForward(address, port)

				if err != nil {
					dialog.ShowError(err, window)
					statusText := fmt.Sprintf("Status: %s", translator.Translate("portforward.status.failed"))
					statusLabel.SetText(statusText)
					startButton.Enable()
					stopButton.Disable()
					openButton.Disable()
					return
				}

				isForwarding = true
				forwardingText := fmt.Sprintf("Forwarding to %s:%s", address, port)
				statusText := fmt.Sprintf("Status: %s", forwardingText)
				statusLabel.SetText(statusText)

				// Enable the open button
				openButton.Enable()

				// Start goroutines to read from streams
				go updateLogs(outStream, "OUT")
				go updateLogs(errStream, "ERR")

				// Add URL to log
				logDisplay.SetText(fmt.Sprintf("Port forwarding started. Access Dify at: http://%s:%s\n\n", address, port))
			}()
		},
	)

	// Open button
	openButton = widget.NewButtonWithIcon(
		translator.Translate("portforward.open"),
		theme.ComputerIcon(),
		func() {
			address := addressEntry.Text
			port := portEntry.Text
			urlStr := fmt.Sprintf("http://%s:%s", address, port)

			u, err := url.Parse(urlStr)
			if err != nil {
				errorMsg := fmt.Sprintf("%s: %s", translator.Translate("error.invalid_url"), err.Error())
				dialog.ShowError(fmt.Errorf(errorMsg), window)
				return
			}

			fyne.CurrentApp().OpenURL(u)
		},
	)
	openButton.Disable() // Disabled by default until port forwarding is active

	// Create a back button
	backButton := components.CreateBackButton(func() {
		// Stop port forwarding if active
		if isForwarding && stopChan != nil {
			close(stopChan)
		}
		onBack()
	})

	// Create header
	header := components.CreateHeader(translator.Translate("portforward.title"))

	// Create info card
	infoCard := components.CreateInfoCard(
		theme.InfoIcon(),
		translator.Translate("portforward.instruction"),
	)

	// Create form for inputs
	inputForm := container.NewGridWithColumns(2,
		widget.NewLabel(translator.Translate("portforward.address")),
		addressEntry,
		widget.NewLabel(translator.Translate("portforward.port")),
		portEntry,
	)

	// Create buttons container
	buttonsContainer := container.NewHBox(
		startButton,
		stopButton,
		openButton,
	)

	// Create log container with scroll
	logScroll := container.NewScroll(logDisplay)
	logScroll.SetMinSize(fyne.NewSize(0, 200))

	// Create main content area
	mainContent := container.NewVBox(
		infoCard,
		inputForm,
		buttonsContainer,
		statusLabel,
		widget.NewLabel(translator.Translate("portforward.logs")),
		logScroll,
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
