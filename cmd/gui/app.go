package gui

import (
	"errors"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/lcpu-club/kube-dify-gui/cmd/gui/operations"
	"github.com/lcpu-club/kube-dify-gui/internal/client"
	"github.com/lcpu-club/kube-dify-gui/internal/i18n"
	"github.com/lcpu-club/kube-dify-gui/internal/settings"
	"github.com/lcpu-club/kube-dify-gui/internal/tray"
	"github.com/lcpu-club/kube-dify-gui/internal/ui"
)

// parseURL parses a URL string and returns a *url.URL
// If the URL is invalid, it returns nil
func parseURL(urlStr string) *url.URL {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil
	}
	return u
}

// App represents the main GUI application
type App struct {
	fyneApp    fyne.App
	window     fyne.Window
	client     *client.Client
	settings   *settings.Settings
	translator *i18n.Translator
	components *ui.Components
	trayMgr    *tray.TrayManager
}

// NewApp creates a new GUI application
func NewApp() *App {
	fyneApp := app.New()

	// Create settings manager
	settingsMgr := settings.New(fyneApp)

	// Set theme based on settings
	isDark := fyneApp.Preferences().BoolWithFallback("dark_theme", false)
	if isDark {
		fyneApp.Settings().SetTheme(theme.DarkTheme())
	} else {
		fyneApp.Settings().SetTheme(theme.LightTheme())
	}

	// Create translator
	translator, _ := i18n.New(fyneApp)

	// Set window title using translation
	window := fyneApp.NewWindow(translator.Translate("app.title"))
	window.Resize(fyne.NewSize(800, 600))

	// Create UI components
	components := ui.New(translator)

	app := &App{
		fyneApp:    fyneApp,
		window:     window,
		settings:   settingsMgr,
		translator: translator,
		components: components,
	}

	// Create tray manager
	app.trayMgr = tray.New(fyneApp, window, translator, app.quit)

	// Set window close handler to minimize to tray instead of closing
	window.SetCloseIntercept(func() {
		window.Hide()
	})

	return app
}

// Run starts the GUI application
func (a *App) Run() {
	// Setup tray icon
	a.trayMgr.Setup()

	// Check for saved token
	if a.settings.GetRememberLogin() {
		savedToken := a.settings.LoadToken()
		if savedToken != "" {
			a.login(savedToken)
		} else {
			a.showLoginScreen()
		}
	} else {
		a.showLoginScreen()
	}

	a.window.ShowAndRun()
}

// quit properly exits the application
func (a *App) quit() {
	a.fyneApp.Quit()
}

// showLoginScreen displays the login screen
func (a *App) showLoginScreen() {
	// Set window title
	a.window.SetTitle(a.translator.Translate("app.title"))

	// Create token entry
	tokenEntry := widget.NewPasswordEntry()
	tokenEntry.SetPlaceHolder(a.translator.Translate("login.token.placeholder"))

	// Create remember me checkbox
	rememberCheck := widget.NewCheck(a.translator.Translate("login.remember"), func(checked bool) {
		a.settings.SetRememberLogin(checked)
	})
	rememberCheck.SetChecked(a.settings.GetRememberLogin())

	// Create form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: a.translator.Translate("login.token"), Widget: tokenEntry},
			{Text: "", Widget: rememberCheck},
		},
		OnSubmit: func() {
			a.login(tokenEntry.Text)
		},
		SubmitText: a.translator.Translate("login.submit"),
	}

	// Create header with settings button
	settingsButton := a.components.CreateSettingsButton(func() {
		a.showSettingsScreen()
	})
	header := a.components.CreateHeader(a.translator.Translate("login.title"), settingsButton)

	// Create layout
	tokenURL := "https://hpcgame.pku.edu.cn/kube/_/ui/#/tokens/"
	tokenURLLabel := widget.NewHyperlink(a.translator.Translate("login.get_token"), parseURL(tokenURL))

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		widget.NewLabel(a.translator.Translate("login.instruction")),
		form,
		tokenURLLabel,
	)

	a.window.SetContent(container.NewPadded(content))
}

// login attempts to create a client with the provided token
func (a *App) login(token string) {
	if token == "" {
		dialog.ShowError(errors.New(a.translator.Translate("error.token_empty")), a.window)
		return
	}

	token = strings.Trim(token, " \r\n\t")

	// Save token if remember login is enabled
	if a.settings.GetRememberLogin() {
		a.settings.SaveToken(token)
	}

	// Show loading dialog
	progress := dialog.NewProgress(
		a.translator.Translate("login.submit"),
		a.translator.Translate("login.connecting"),
		a.window,
	)
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
	// Set window title
	a.window.SetTitle(a.translator.Translate("app.title") + " - " + a.translator.Translate("dashboard.title"))

	// Create buttons for each operation with tooltips
	initButton := a.components.CreateOperationButton(
		theme.DocumentCreateIcon(),
		"dashboard.init",
		"dashboard.init.tooltip",
		func() { a.showInitializeScreen() },
	)

	startButton := a.components.CreateOperationButton(
		theme.MediaPlayIcon(),
		"dashboard.start",
		"dashboard.start.tooltip",
		func() { a.showStartScreen() },
	)

	stopButton := a.components.CreateOperationButton(
		theme.MediaStopIcon(),
		"dashboard.stop",
		"dashboard.stop.tooltip",
		func() { a.showStopScreen() },
	)

	deleteButton := a.components.CreateOperationButton(
		theme.DeleteIcon(),
		"dashboard.delete",
		"dashboard.delete.tooltip",
		func() { a.showDeleteScreen() },
	)

	portForwardButton := a.components.CreateOperationButton(
		theme.ComputerIcon(),
		"dashboard.portforward",
		"dashboard.portforward.tooltip",
		func() { a.showPortForwardScreen() },
	)

	// Create header with settings and logout buttons
	settingsButton := a.components.CreateSettingsButton(func() {
		a.showSettingsScreen()
	})

	logoutButton := a.components.CreateLogoutButton(func() {
		a.logout()
	})

	header := a.components.CreateHeader(
		a.translator.Translate("dashboard.title"),
		settingsButton,
		logoutButton,
	)

	// Create layout with cards for better visual organization
	operationsCard := a.components.CreateCard(
		a.translator.Translate("dashboard.instruction"),
		container.NewGridWithColumns(2,
			initButton,
			startButton,
			stopButton,
			deleteButton,
			portForwardButton,
		),
	)

	content := container.NewVBox(
		header,
		widget.NewSeparator(),
		operationsCard,
		widget.NewLabel(a.translator.Translate("dashboard.tips")),
	)

	a.window.SetContent(container.NewPadded(content))
}

// showSettingsScreen displays the settings screen
func (a *App) showSettingsScreen() {
	operations.ShowSettings(
		a.window,
		a.translator,
		func() {
			// If we're logged in, go back to dashboard, otherwise go to login
			if a.client != nil {
				a.showDashboard()
			} else {
				a.showLoginScreen()
			}
		},
		func(lang string) {
			// Change language
			a.translator.SetLanguage(lang)

			// Refresh current screen
			if a.client != nil {
				a.showDashboard()
			} else {
				a.showLoginScreen()
			}

			// Update tray
			a.trayMgr.UpdateTranslation()
		},
		func(isDark bool) {
			// Change theme
			if isDark {
				a.fyneApp.Settings().SetTheme(theme.DarkTheme())
			} else {
				a.fyneApp.Settings().SetTheme(theme.LightTheme())
			}

			// Save preference
			a.fyneApp.Preferences().SetBool("dark_theme", isDark)
		},
	)
}

// logout logs out the current user
func (a *App) logout() {
	// Show confirmation dialog
	dialog.ShowConfirm(
		a.translator.Translate("dashboard.logout"),
		a.translator.Translate("dashboard.logout")+"?",
		func(confirmed bool) {
			if confirmed {
				// Clear client
				a.client = nil

				// Clear saved token if remember login is enabled
				if a.settings.GetRememberLogin() {
					a.settings.ClearToken()
				}

				// Show login screen
				a.showLoginScreen()
			}
		},
		a.window,
	)
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
