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
)

// ShowPortForward displays the port forward operation screen
func ShowPortForward(window fyne.Window, client *client.Client, onBack func()) {
	// Create input fields
	addressEntry := widget.NewEntry()
	addressEntry.SetPlaceHolder("Address (e.g., 0.0.0.0)")
	addressEntry.SetText("127.0.0.1")

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Port (e.g., 28080)")
	portEntry.SetText("28080")

	// Create status label
	statusLabel := widget.NewLabel("Status: Not forwarding")

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
						Title:   "Port Forward Error",
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
	stopButton = widget.NewButtonWithIcon("Stop Port Forward", theme.MediaStopIcon(), func() {
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
		statusLabel.SetText("Status: Not forwarding")
		startButton.Enable()
		stopButton.Disable()
		openButton.Disable()
	})
	stopButton.Disable()

	// Start button
	startButton = widget.NewButtonWithIcon("Start Port Forward", theme.MediaPlayIcon(), func() {
		address := addressEntry.Text
		port := portEntry.Text

		if address == "" {
			dialog.ShowError(fmt.Errorf("address cannot be empty"), window)
			return
		}

		if port == "" {
			dialog.ShowError(fmt.Errorf("port cannot be empty"), window)
			return
		}

		// Disable start button and enable stop button
		startButton.Disable()
		stopButton.Enable()

		// Clear log display
		logDisplay.SetText("")

		// Update status
		statusLabel.SetText("Status: Starting port forward...")

		// Run port forward in a goroutine
		go func() {
			var err error
			stopChan, outStream, errStream, err = client.DoPortForward(address, port)

			if err != nil {
				dialog.ShowError(err, window)
				statusLabel.SetText("Status: Failed to start port forward")
				startButton.Enable()
				stopButton.Disable()
				openButton.Disable()
				return
			}

			isForwarding = true
			statusLabel.SetText(fmt.Sprintf("Status: Forwarding to %s:%s", address, port))

			// Enable the open button
			openButton.Enable()

			// Start goroutines to read from streams
			go updateLogs(outStream, "OUT")
			go updateLogs(errStream, "ERR")

			// Add URL to log
			logDisplay.SetText(fmt.Sprintf("Port forwarding started. Access Dify at: http://%s:%s\n\n", address, port))
		}()
	})

	// Open button
	openButton = widget.NewButtonWithIcon("Open in Browser", theme.ComputerIcon(), func() {
		address := addressEntry.Text
		port := portEntry.Text
		urlStr := fmt.Sprintf("http://%s:%s", address, port)

		u, err := url.Parse(urlStr)
		if err != nil {
			dialog.ShowError(fmt.Errorf("invalid URL: %v", err), window)
			return
		}

		fyne.CurrentApp().OpenURL(u)
	})
	openButton.Disable() // Disabled by default until port forwarding is active

	// Create a back button
	backButton := widget.NewButtonWithIcon("Back", theme.NavigateBackIcon(), func() {
		// Stop port forwarding if active
		if isForwarding && stopChan != nil {
			close(stopChan)
		}
		onBack()
	})

	logScroll := container.NewScroll(logDisplay)
	logScroll.SetMinSize(logDisplay.MinSize())

	// Create layout
	content := container.NewVBox(
		widget.NewLabel("Port Forward"),
		widget.NewLabel("Forward the Dify application to a local port"),
		container.NewGridWithColumns(2,
			widget.NewLabel("Address:"),
			addressEntry,
			widget.NewLabel("Port:"),
			portEntry,
		),
		container.NewHBox(
			startButton,
			stopButton,
			openButton,
		),
		statusLabel,
		widget.NewLabel("Logs:"),
		logScroll,
		backButton,
	)

	window.SetContent(content)
}
