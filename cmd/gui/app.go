package gui

import (
	"errors"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/lcpu-club/kube-dify-gui/cmd/gui/operations"
	"github.com/lcpu-club/kube-dify-gui/internal/client"
)

// App represents the main GUI application
type App struct {
	fyneApp fyne.App
	window  fyne.Window
	client  *client.Client
}

// NewApp creates a new GUI application
func NewApp() *App {
	fyneApp := app.New()
	fyneApp.Settings().SetTheme(theme.LightTheme())

	window := fyneApp.NewWindow("Kube Dify GUI")
	window.Resize(fyne.NewSize(800, 600))

	return &App{
		fyneApp: fyneApp,
		window:  window,
	}
}

// Run starts the GUI application
func (a *App) Run() {
	a.showLoginScreen()
	a.window.ShowAndRun()
}

// showLoginScreen displays the login screen
func (a *App) showLoginScreen() {
	tokenEntry := widget.NewPasswordEntry()
	tokenEntry.SetPlaceHolder("Enter your HPCGame token")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Token", Widget: tokenEntry},
		},
		OnSubmit: func() {
			a.login(tokenEntry.Text)
		},
		SubmitText: "Login",
	}

	content := container.NewVBox(
		widget.NewLabel("Welcome to Kube Dify GUI"),
		widget.NewLabel("Please enter your HPCGame token to continue"),
		form,
		widget.NewLabel("Get your token from https://hpcgame.pku.edu.cn/kube/_/ui/#/tokens/"),
	)

	a.window.SetContent(content)
}

// login attempts to create a client with the provided token
func (a *App) login(token string) {
	if token == "" {
		dialog.ShowError(errors.New("Token cannot be empty"), a.window)
		return
	}

	token = strings.Trim(token, " \r\n\t")

	// Show loading dialog
	progress := dialog.NewProgress("Logging in", "Connecting to HPCGame cluster...", a.window)
	progress.Show()

	// Create client in a goroutine to avoid blocking the UI
	go func() {
		client, err := client.NewHPCGame(token)

		// Hide progress dialog
		progress.Hide()

		if err != nil {
			dialog.ShowError(err, a.window)
			return
		}

		a.client = client
		a.showDashboard()
	}()
}

// showDashboard displays the main dashboard
func (a *App) showDashboard() {
	a.window.SetTitle("Kube Dify GUI - Dashboard")

	// Create buttons for each operation
	initButton := widget.NewButtonWithIcon("Initialize", theme.DocumentCreateIcon(), func() {
		a.showInitializeScreen()
	})

	startButton := widget.NewButtonWithIcon("Start", theme.MediaPlayIcon(), func() {
		a.showStartScreen()
	})

	stopButton := widget.NewButtonWithIcon("Stop", theme.MediaStopIcon(), func() {
		a.showStopScreen()
	})

	deleteButton := widget.NewButtonWithIcon("Delete", theme.DeleteIcon(), func() {
		a.showDeleteScreen()
	})

	portForwardButton := widget.NewButtonWithIcon("Port Forward", theme.ComputerIcon(), func() {
		a.showPortForwardScreen()
	})

	// Create layout
	content := container.NewVBox(
		widget.NewLabel("Kube Dify Dashboard"),
		widget.NewLabel("Select an operation to perform:"),
		container.NewGridWithColumns(2,
			initButton,
			startButton,
			stopButton,
			deleteButton,
			portForwardButton,
		),
	)

	a.window.SetContent(content)
}

// showInitializeScreen displays the initialize operation screen
func (a *App) showInitializeScreen() {
	operations.ShowInitialize(a.window, a.client, func() {
		a.showDashboard()
	})
}

// showStartScreen displays the start operation screen
func (a *App) showStartScreen() {
	operations.ShowStart(a.window, a.client, func() {
		a.showDashboard()
	})
}

// showStopScreen displays the stop operation screen
func (a *App) showStopScreen() {
	operations.ShowStop(a.window, a.client, func() {
		a.showDashboard()
	})
}

// showDeleteScreen displays the delete operation screen
func (a *App) showDeleteScreen() {
	operations.ShowDelete(a.window, a.client, func() {
		a.showDashboard()
	})
}

// showPortForwardScreen displays the port forward operation screen
func (a *App) showPortForwardScreen() {
	operations.ShowPortForward(a.window, a.client, func() {
		a.showDashboard()
	})
}
