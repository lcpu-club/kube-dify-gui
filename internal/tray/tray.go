package tray

import (
	"fyne.io/fyne/v2"
	"fyne.io/systray"
	"github.com/lcpu-club/kube-dify-gui/internal/i18n"
)

// TrayManager handles the system tray functionality
type TrayManager struct {
	app        fyne.App
	window     fyne.Window
	translator *i18n.Translator
	onQuit     func()
}

// New creates a new tray manager
func New(a fyne.App, w fyne.Window, t *i18n.Translator, onQuit func()) *TrayManager {
	return &TrayManager{
		app:        a,
		window:     w,
		translator: t,
		onQuit:     onQuit,
	}
}

// Setup initializes the system tray
func (t *TrayManager) Setup() {
	// Start the systray in a goroutine
	go systray.Run(t.onReady, t.onExit)
}

// onReady is called when the systray is ready
func (t *TrayManager) onReady() {
	// Set the icon (this will use the app icon by default)
	systray.SetIcon(t.app.Icon().Content())
	systray.SetTitle("Kube Dify GUI")
	systray.SetTooltip("Kube Dify GUI")

	// Create menu items
	mShow := systray.AddMenuItem(t.translator.Translate("tray.show"), "Show the application window")
	mHide := systray.AddMenuItem(t.translator.Translate("tray.hide"), "Hide the application window")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem(t.translator.Translate("tray.quit"), "Quit the application")

	// Handle menu item clicks
	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				t.window.Show()
				t.window.RequestFocus()
			case <-mHide.ClickedCh:
				t.window.Hide()
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

// onExit is called when the systray is exiting
func (t *TrayManager) onExit() {
	if t.onQuit != nil {
		t.onQuit()
	}
}

// UpdateTranslation updates the tray menu with the current language
func (t *TrayManager) UpdateTranslation() {
	// This is a placeholder for now
	// In a real implementation, we would need to recreate the menu items
	// with the new translations, but systray doesn't support updating
	// existing menu items directly
}
